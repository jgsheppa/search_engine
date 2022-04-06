package models

import (
	"github.com/RediSearch/redisearch-go/redisearch"
	"strings"
)

func CreateSuggestions(article Article) []redisearch.Suggestion {
	var suggestion []redisearch.Suggestion

	splitTitle := strings.Split(article.Title, " ")

	for _, word := range splitTitle {
		suggestion = append(suggestion, redisearch.Suggestion{
			Term:    word,
			Score:   100,
			Payload: word,
			Incr:    false,
		})
	}

	splitAuthor := strings.Split(article.Author, " ")

	for _, word := range splitAuthor {
		suggestion = append(suggestion, redisearch.Suggestion{
			Term:    word,
			Score:   100,
			Payload: word,
			Incr:    false,
		})
	}

	splitTopic := strings.Split(article.Topic, " ")

	for _, word := range splitTopic {
		suggestion = append(suggestion, redisearch.Suggestion{
			Term:    word,
			Score:   100,
			Payload: word,
			Incr:    false,
		})
	}

	return suggestion
}
