package models

import (
	"github.com/RediSearch/redisearch-go/redisearch"
	"log"
	"time"
)

type Guide struct {
	// Document name, if possible a UUID
	Document string `json:"document"`
	// URL of article
	URL string `json:"url"`
	// Title of article
	Text string `json:"title"`
	// Topics of article
	Topic string `json:"topic"`
}

type Guides []Guide

func CreateGuideDocument(rs redisearch.Client, autoCompleter redisearch.Autocompleter, guides Guides) {
	var documents []redisearch.Document

	for _, guide := range guides {
		suggestion := CreateGuideSuggestions(guide)

		err := autoCompleter.AddTerms(suggestion...)
		if err != nil {
			log.Println("Error adding term for autocomplete")
		}

		doc := redisearch.NewDocument(guide.Document, 1.0)
		doc.Set("text", guide.Text).
			Set("url", guide.URL).
			Set("topic", guide.Topic).
			Set("date", time.Now().Unix())
		documents = append(documents, doc)
	}

	// Index the document. The API accepts multiple documents at a time
	if err := rs.Index(documents...); err != nil {
		log.Fatal(err)
	}
}
