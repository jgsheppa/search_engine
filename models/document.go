package models

import (
	"github.com/RediSearch/redisearch-go/redisearch"
	"log"
	"time"
)

type RedisDB struct {
	redisSearch redisearch.Client
}

type Document struct {
	// Document name, if possible a UUID
	Document string `json:"document"`
	// URL of article
	URL string `json:"url"`
	// Title of article
	Text string `json:"title"`
	// Topics of article
	Topic string `json:"topic"`
}

type Field struct {
	Name string `json:"name"`
	Type string `json:"type"`
}

type Documents []Document

func (s *Services) CreateDocument(documents Documents) error {
	var redisDocuments []redisearch.Document

	for _, document := range documents {
		suggestion := CreateSuggestions(document)

		err := s.Autocomplete.AddTerms(suggestion...)
		if err != nil {
			log.Println("Error adding term for autocomplete")
		}

		doc := redisearch.NewDocument(document.Document, 1.0)
		doc.Set("text", document.Text).
			Set("url", document.URL).
			Set("topic", document.Topic).
			Set("date", time.Now().Unix())
		redisDocuments = append(redisDocuments, doc)
	}

	// Index the document. The API accepts multiple documents at a time
	if err := s.Redisearch.Index(redisDocuments...); err != nil {
		return err
	}
	return nil
}

func (s *Services) DeleteDocument(document string) error {
	err := s.Redisearch.DeleteDocument(document)
	if err != nil {
		return err
	}
	return nil
}
