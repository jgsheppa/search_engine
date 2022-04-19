package models

import (
	"github.com/RediSearch/redisearch-go/redisearch"
	"strconv"
)

type RedisDB struct {
	redisSearch redisearch.Client
}

type Document struct {
	// URL of article
	Country string `json:"country"`
	// Title of article
	Name string `json:"name"`
	// Latitude of city
	Lat string `json:"lat"`
	// Longitude of city
	Long string `json:"lng"`
}

type Field struct {
	Name string `json:"name"`
	Type string `json:"type"`
}

type Documents []Document

func (s *Services) CreateDocument(documents Documents) error {
	var redisDocuments []redisearch.Document

	for i, document := range documents {
		suggestion := CreateSuggestions(document)

		err := s.Autocomplete.AddTerms(suggestion...)
		if err != nil {
			return &serverError
		}

		doc := redisearch.NewDocument("document:"+document.Name+strconv.Itoa(i), 1.0)
		doc.Set("city", document.Name).
			Set("country", document.Country).
			Set("location", document.Lat+","+document.Long).
			Set("id", i)
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
