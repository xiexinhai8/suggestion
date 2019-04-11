package es

import (
	// "encoding/json";
	"gopkg.in/olivere/elastic.v6"
	"encoding/json"
	log "github.com/sirupsen/logrus"
	"context"
)

type Suggestion struct {
	Name  string  `json:"name"`
	City  string  `json:"city"`
	Score float64 `json:"score"`
}

type Source struct {
	Country string	`json:"country"`
	City	string	`json:"city"`
	Name	string	`json:"name"`
}

var (
	elasticClient *elastic.Client
)

type EsSuggestor struct {
	SuggestName string
	Field       string
	Pretty      bool
	Client      *elastic.Client
}

func (esSuggestor *EsSuggestor) Search(query string, esIndex string) []Suggestion {
	s := elastic.NewCompletionSuggester(esSuggestor.SuggestName).Prefix(query).Field(esSuggestor.Field)
	src, err := s.Source(true)
	if err != nil {
		log.Errorf("extract query occurred error: ", err)
	}
	data, err := json.Marshal(src)
	if err != nil {
		log.Errorf("json dumps occurred error: ", err)
	}
	got := string(data)
	log.Infof("request_body: %v", got)
	searchResult, err := esSuggestor.Client.Search().
		Index(esIndex).
		Query(elastic.NewMatchAllQuery()).
		Suggester(s).Pretty(true).
		Do(context.TODO())
	if err != nil {
		log.Panicf("search suggest occurred error: ", err)
	}
	if searchResult.Suggest == nil {
		log.Errorf("search suggest result is nil")
	}
	result, fond := searchResult.Suggest[esSuggestor.SuggestName]
	log.Infof("request_body: %v, find: %v, result_num: %v", got, fond, len(result))

	suggestions := make([]Suggestion, 0)
	for _, suggests := range result {
		for _, option := range suggests.Options {
			var source = new(Source)
			err := json.Unmarshal(*option.Source, source)
			if err != nil {
				log.Errorf("parse es source info occurred error %v", err)
			}
			suggestions = append(suggestions, Suggestion{
				Name: source.Name,
				City: source.City,
				Score: option.ScoreUnderscore})
		}
	}
	log.Infof("request_body: %v, result: %v", got, suggestions)
	return suggestions
}
