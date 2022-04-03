package models

import (
	"github.com/RediSearch/redisearch-go/redisearch"
	"strings"
)

func CreateSuggestions(todo Todo) []redisearch.Suggestion {
	var suggestion []redisearch.Suggestion

	splitTitle := strings.Split(todo.Title, " ")

	for _, word := range splitTitle {
		suggestion = append(suggestion, redisearch.Suggestion{
			Term:    word,
			Score:   100,
			Payload: word,
			Incr:    false,
		})
	}

	splitAuthor := strings.Split(todo.Author, " ")

	for _, word := range splitAuthor {
		suggestion = append(suggestion, redisearch.Suggestion{
			Term:    word,
			Score:   100,
			Payload: word,
			Incr:    false,
		})
	}

	return suggestion
}
