package models

import (
	"github.com/RediSearch/redisearch-go/redisearch"
	"github.com/gomodule/redigo/redis"
)

// CreateSchema is used to create the schema for your Redisearch documents,
// which will allow you to add your data in the form of these documents
func CreateSchema() *redisearch.Schema {
	sc := redisearch.NewSchema(redisearch.DefaultOptions).
		AddField(redisearch.NewNumericField("date")).
		AddField(redisearch.NewTextFieldOptions("text", redisearch.TextFieldOptions{Weight: 5.0, Sortable: true})).
		AddField(redisearch.NewTextFieldOptions("url", redisearch.TextFieldOptions{Weight: 5.0, Sortable: true})).
		AddField(redisearch.NewTextFieldOptions("topic", redisearch.TextFieldOptions{Weight: 5.0, Sortable: true}))
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
