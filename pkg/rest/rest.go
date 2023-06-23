package rest

import (
	"embed"
	"github.com/Dimss/raspanme/pkg/store"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"io/fs"
	"net/http"
)

//go:embed ui/*
var HtmlAssets embed.FS

type Server struct {
}

func NewServer() *Server {
	return &Server{}
}

func mustFS() http.FileSystem {
	sub, err := fs.Sub(HtmlAssets, "ui")

	if err != nil {
		panic(err)
	}

	return http.FS(sub)
}

func (s *Server) Start() {
	r := gin.Default()
	r.StaticFS("/ui", mustFS())
	v1 := r.Group("/v1")
	{
		v1.GET("/category", s.GetCategories)
		v1.GET("/question/:catID", s.GetQuestions)
	}
	r.Run()
}

func (s *Server) GetCategories(c *gin.Context) {
	var categories []store.Category
	if err := store.Db().Model(&store.Category{}).Order("description").Find(&categories).Error; err != nil {
		c.JSON(http.StatusBadGateway, gin.H{"message": err})
		return
	}
	c.JSON(http.StatusOK, categories)
}

func (s *Server) GetQuestions(c *gin.Context) {
	catId := c.Param("catID")
	var questions []store.Question
	if err := store.Db().Model(&store.Question{}).
		Preload("Answer").
		Where("category_id = ? and lang_id = ?", catId, 1).
		Find(&questions).Error; err != nil {
		zap.S().Error(err)
		c.JSON(http.StatusBadGateway, gin.H{"message": err})
		return
	}
	var res [][]store.Question
	for _, q := range questions {
		ruQ := store.Question{}
		if err := store.Db().Model(&store.Question{}).
			Preload("Answer").
			Where("q_id = ? and lang_id = ?", q.QID, 2).First(&ruQ).Error; err != nil {
			zap.S().Error(err)
		}
		res = append(res, []store.Question{q, ruQ})
	}
	c.JSON(http.StatusOK, res)
}
