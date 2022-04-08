package models

import (
	"github.com/RediSearch/redisearch-go/redisearch"
	"strings"
)

func CreateSuggestions(article Article) []redisearch.Suggestion {
	var suggestion []redisearch.Suggestion

	splitTitle := strings.Split(article.Title, " ")
	title := SuggestionFactory(splitTitle)
	suggestion = append(suggestion, title...)

	splitAuthor := strings.Split(article.Author, " ")
	author := SuggestionFactory(splitAuthor)
	suggestion = append(suggestion, author...)

	splitTopic := strings.Split(article.Topic, " ")
	topic := SuggestionFactory(splitTopic)
	suggestion = append(suggestion, topic...)

	return suggestion
}

func SuggestionFactory(wordArray []string) []redisearch.Suggestion {
	var suggestion []redisearch.Suggestion

	for index, word := range wordArray {
		wordArrayLength := len(wordArray) - 1
		increment := 2
		stopSlice := increment + index

		if index == wordArrayLength {
			return suggestion
		}
		stringPortion := strings.Join(wordArray[0:stopSlice], " ")

		suggestion = append(suggestion, redisearch.Suggestion{
			Term:    stringPortion,
			Score:   100,
			Payload: stringPortion,
			Incr:    false,
		})

		stringPortion = strings.Join(wordArray[index:stopSlice], " ")

		suggestion = append(suggestion, redisearch.Suggestion{
			Term:    stringPortion,
			Score:   100,
			Payload: stringPortion,
			Incr:    false,
		})

		suggestion = append(suggestion, redisearch.Suggestion{
			Term:    word,
			Score:   100,
			Payload: word,
			Incr:    false,
		})

	}

	return suggestion
}
