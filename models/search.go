package models

import (
	"github.com/RediSearch/redisearch-go/redisearch"
	"github.com/gomodule/redigo/redis"
	"log"
	"strings"
	"time"
)

type RedisDB struct {
	redisClient redis.Conn
	redisSearch redisearch.Client
}

type Todo struct {
	Document string `json:"document"`
	Author   string `json:"author"`
	ID       int    `json:"id"`
	URL      string `json:"url"`
	Title    string `json:"title"`
}

type Field struct {
	Name string `json:"name"`
	Type string `json:"type"`
}

type Todos []Todo

func CreateSuggestions(todo Todo) []redisearch.Suggestion {
	var suggestion []redisearch.Suggestion

	splitTitle := strings.Split(todo.Title, " ")

	for _, word := range splitTitle {
		suggestion = append(suggestion, redisearch.Suggestion{
			Term:    word,
			Score:   100,
			Payload: todo.Title,
			Incr:    false,
		})
	}

	splitAuthor := strings.Split(todo.Author, " ")

	for _, word := range splitAuthor {
		suggestion = append(suggestion, redisearch.Suggestion{
			Term:    word,
			Score:   100,
			Payload: todo.Author,
			Incr:    false,
		})
	}

	return suggestion
}

func CreateDocument(rs redisearch.Client, autoCompleter redisearch.Autocompleter, todos Todos) {
	var documents []redisearch.Document

	for _, todo := range todos {
		suggestion := CreateSuggestions(todo)

		autoCompleter.AddTerms(suggestion...)

		doc := redisearch.NewDocument(todo.Document, 1.0)
		doc.Set("author", todo.Author).
			Set("id", todo.ID).
			Set("title", todo.Title).
			Set("url", todo.URL).
			Set("date", time.Now().Unix())
		documents = append(documents, doc)
	}

	// Index the document. The API accepts multiple documents at a time
	if err := rs.Index(documents...); err != nil {
		log.Fatal(err)
	}
}

func DeleteDocument(rs redisearch.Client, document string) error {
	err := rs.DeleteDocument(document)
	if err != nil {
		return err
	}
	return nil
}
