package models

import (
	"github.com/RediSearch/redisearch-go/redisearch"
	"github.com/gomodule/redigo/redis"
)

// These values are used to build the schema and in the controller package
// for sortable and highlighted fields.
const (
	Name   = "name"
	Link   = "link"
	Active = "active"
)

// Document should reflect the data structure you would like to index with Redis
type Document struct {
	// Name of NBA player
	Name string `json:"name"`
	// Link to player page
	Link string `json:"link"`
	// Active represents the years a player was active
	Active string `json:"active"`
}

type Documents []Document

// CreateSchema is used to create the schema for your Redisearch documents,
// which will allow you to add your data in the form of these documents
func CreateSchema() *redisearch.Schema {
	sc := redisearch.NewSchema(redisearch.DefaultOptions).
		AddField(redisearch.NewTextFieldOptions(Name,
			redisearch.TextFieldOptions{Weight: 5.0, Sortable: true})).
		AddField(redisearch.NewTextFieldOptions(Link,
			redisearch.TextFieldOptions{Weight: 5.0, Sortable: true})).
		AddField(redisearch.NewTextFieldOptions(Active,
			redisearch.TextFieldOptions{Weight: 5.0, Sortable: true}))

	return sc
}

func CreateIndex(pool *redis.Pool, indexName string) (*redisearch.Client, *redisearch.Autocompleter, error) {
	// Create a client. By default a client is schemaless
	// unless a schema is provided when creating the index
	c := redisearch.NewClientFromPool(pool, indexName)
	autocomplete := redisearch.NewAutocompleterFromPool(pool, indexName)

	// Create a schema
	sc := CreateSchema()

	// Create the index with the given schema
	err := c.CreateIndex(sc)

	return c, autocomplete, err
}
