package parser

import (
	"encoding/base64"
	"fmt"
	"github.com/Dimss/raspanme/pkg/store"
	"github.com/xuri/excelize/v2"
	"go.uber.org/zap"
	"strings"
)

type Parser struct {
	file          *excelize.File
	sheetName     string
	categoryIndex int
	titleIndex    int
	rows          [][]string
}

func NewParser(path string) *Parser {
	f, err := excelize.OpenFile(path)
	if err != nil {
		zap.S().Fatal(err)
	}

	p := &Parser{
		file:          f,
		sheetName:     "rptSheelotAndTshuvot",
		categoryIndex: 6,
		titleIndex:    1,
	}
	p.loadRows()
	return p
}

func (p *Parser) loadRows() {
	rows, err := p.file.GetRows(p.sheetName)
	if err != nil {
		zap.S().Error(err)
		return
	}
	p.rows = rows
}

func (p *Parser) categoryKey(category string) string {
	return base64.StdEncoding.EncodeToString([]byte(p.parseCategory(category)))
}

func (p *Parser) LoadCategory() {
	categories := map[string]*store.Category{}
	for i := p.titleIndex + 2; i < len(p.rows); i++ {
		category, err := p.file.GetCellValue(p.sheetName, fmt.Sprintf("G%d", i))
		if err != nil {
			zap.S().Error(err)
		}
		if category == "" {
			continue
		}
		categoryKey := p.categoryKey(category)
		if _, exists := categories[categoryKey]; !exists {
			categories[categoryKey] = &store.Category{Category: categoryKey, Description: category, LangID: 1}
		}
	}
	for _, c := range categories {
		if res := store.Db().Create(c); res.Error != nil {
			zap.S().Info(res.Error)
		}
	}
}

func (p *Parser) isRightAnswer(rowNumber int) bool {
	if p.rows[rowNumber][1] == "+" {
		return true
	}
	return false
}

func (p *Parser) categoryId(category string) uint {
	c := &store.Category{}
	store.Db().Where("category = ?", p.categoryKey(category)).First(c)
	return c.ID
}

func (p *Parser) LoadQuestionsAndAnswers() {
	i := p.titleIndex + 1

	for {
		if len(p.rows[i]) <= 9 && p.rows[i][8] != "" {
			q := &store.Question{LangID: 1}

			q.QID = p.rows[i][8]
			q.Question = p.rows[i][4]
			q.CategoryID = p.categoryId(p.rows[i][6])

			for {

				if len(p.rows) <= i || p.rows[i] == nil {
					break
				}

				q.Answer = append(q.Answer, store.Answer{
					LangID:      1,
					Answer:      p.rows[i][2],
					Right:       p.isRightAnswer(i),
					AnswerIndex: p.rows[i][3],
				})
				i++
			}

			if res := store.Db().Create(q); res.Error != nil {
				zap.S().Error(res.Error)
			}

		}

		i++
		if i >= len(p.rows) {
			break
		}
	}
}

func (p *Parser) parseCategory(category string) (parsedCategory string) {
	for _, c := range strings.Split(category, " ") {
		parsedCategory += strings.TrimSpace(c)
	}
	return
}
