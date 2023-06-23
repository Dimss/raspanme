package translate

import (
	"bytes"
	"encoding/json"
	"github.com/Dimss/raspanme/pkg/store"
	"go.uber.org/zap"
	"log"
	"net/http"
	"net/url"
)

type translateRequest struct {
	Text string
}

type translationsSubResponse struct {
	Text string
	To   string
}

type translateResponse struct {
	Translations []translationsSubResponse
}

type Translator struct {
	ApiKey   string
	Endpoint string
}

func NewTranslator(apiKey, endpoint string) *Translator {
	return &Translator{
		ApiKey:   apiKey,
		Endpoint: endpoint,
	}
}

func (t *Translator) Run() {
	var questions []store.Question
	if err := store.Db().Model(&store.Question{}).
		Preload("Answer").Find(&questions).Error; err != nil {
		zap.S().Error(err)
		return
	}

	for _, q := range questions {
		tr := []translateRequest{{Text: q.Question}}
		for _, a := range q.Answer {
			tr = append(tr, translateRequest{Text: a.Answer})
		}
		translated := t.translate(tr)
		trQuestion := &store.Question{
			Question:   translated[0].Translations[0].Text,
			QID:        q.QID,
			LangID:     2,
			CategoryID: q.CategoryID,
		}
		for i := 1; i < 5; i++ {
			trAnswer := store.Answer{
				Answer:      translated[i].Translations[0].Text,
				Right:       q.Answer[i-1].Right,
				AnswerIndex: q.Answer[i-1].AnswerIndex,
				LangID:      2,
			}
			trQuestion.Answer = append(trQuestion.Answer, trAnswer)
		}
		if res := store.Db().Create(trQuestion); res.Error != nil {
			zap.S().Error(res.Error)
		}
		zap.S().Infof("translated: %s", q.QID)

	}

}

func (t *Translator) translate(translateRequest []translateRequest) []translateResponse {
	u, _ := url.Parse(t.Endpoint)
	q := u.Query()
	q.Add("from", "he")
	q.Add("to", "ru")
	u.RawQuery = q.Encode()

	// Create an anonymous struct for your request body and encode it to JSON
	//body := []translateRequest{{Text: "אחד"}, {Text: "שתים"}, {Text: "שלוש"}}
	body := translateRequest

	b, _ := json.Marshal(body)

	// Build the HTTP POST request
	req, err := http.NewRequest("POST", u.String(), bytes.NewBuffer(b))
	if err != nil {
		log.Fatal(err)
	}
	// Add required headers to the request
	req.Header.Add("Ocp-Apim-Subscription-Key", t.ApiKey)
	req.Header.Add("Ocp-Apim-Subscription-Region", "eastus")
	req.Header.Add("Content-Type", "application/json")

	// Call the Translator API
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Fatal(err)
	}
	defer res.Body.Close()

	var result []translateResponse
	if err := json.NewDecoder(res.Body).Decode(&result); err != nil {
		log.Fatal(err)
	}
	return result

}
