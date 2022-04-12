package models

import (
	"github.com/RediSearch/redisearch-go/redisearch"
	"log"
	"time"
)

type RedisDB struct {
	redisSearch redisearch.Client
}

type Article struct {
	// Document name, if possible a UUID
	Document string `json:"document"`
	// Author of article
	Author string `json:"author"`
	ID     int    `json:"id"`
	// URL of article
	URL string `json:"url"`
	// Title of article
	Title string `json:"title"`
	// Topics of article
	Topic string `json:"topic"`
}

type Field struct {
	Name string `json:"name"`
	Type string `json:"type"`
}

type Articles []Article

func CreateDocument(rs redisearch.Client, autoCompleter redisearch.Autocompleter, articles Articles) error {
	var documents []redisearch.Document

	for _, article := range articles {
		suggestion := CreateSuggestions(article)

		err := autoCompleter.AddTerms(suggestion...)
		if err != nil {
			log.Println("Error adding term for autocomplete")
		}

		doc := redisearch.NewDocument(article.Document, 1.0)
		doc.Set("author", article.Author).
			Set("id", article.ID).
			Set("title", article.Title).
			Set("url", article.URL).
			Set("topic", article.Topic).
			Set("date", time.Now().Unix())
		documents = append(documents, doc)
	}

	// Index the document. The API accepts multiple documents at a time
	if err := rs.Index(documents...); err != nil {
		return err
	}
	return nil
}

func DeleteDocument(rs redisearch.Client, document string) error {
	err := rs.DeleteDocument(document)
	if err != nil {
		return err
	}
	return nil
}
