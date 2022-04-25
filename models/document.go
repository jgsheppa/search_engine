package models

import (
	"github.com/RediSearch/redisearch-go/redisearch"
)

type RedisDB struct {
	redisSearch redisearch.Client
}

func (s *Services) CreateDocument(documents Documents) error {
	var redisDocuments []redisearch.Document

	for _, document := range documents {
		suggestion := CreateSuggestions(document)

		err := s.Autocomplete.AddTerms(suggestion...)
		if err != nil {
			return &serverError
		}

		doc := redisearch.NewDocument("document:"+document.Link, 1.0)
		doc.Set("name", document.Name).
			Set("link", document.Link).
			Set("active", document.Active)
		redisDocuments = append(redisDocuments, doc)
	}

	// Index the document. The API accepts multiple documents at a time
	if err := s.Redisearch.Index(redisDocuments...); err != nil {
		return &serverError
	}
	return nil
}

func (s *Services) DeleteDocument(document string) error {
	err := s.Redisearch.DeleteDocument(document)
	if err != nil {
		return &serverError
	}
	return nil
}
