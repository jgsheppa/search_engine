package models

import (
	"github.com/RediSearch/redisearch-go/redisearch"
	"strings"
)

func CreateSuggestions(document Document) []redisearch.Suggestion {
	var suggestion []redisearch.Suggestion

	splitText := strings.Split(document.Text, " ")
	text := SuggestionFactory(splitText)
	suggestion = append(suggestion, text...)

	splitTopic := strings.Split(document.Topic, " ")
	topic := SuggestionFactory(splitTopic)
	suggestion = append(suggestion, topic...)

	return suggestion
}

func SuggestionFactory(wordArray []string) []redisearch.Suggestion {
	var suggestion []redisearch.Suggestion

	for index, word := range wordArray {
		wordArrayLength := len(wordArray) - 1
		increment := 2
		limit := increment + index

		if index == wordArrayLength {
			// Add single word as suggestion
			suggestion = append(suggestion, redisearch.Suggestion{
				Term:    word,
				Score:   100,
				Payload: word,
				Incr:    false,
			})

			return suggestion
		}

		// Add beginning of phrase to limit as suggestion
		stringPortion := strings.Join(wordArray[:limit], " ")
		suggestion = append(suggestion, redisearch.Suggestion{
			Term:    stringPortion,
			Score:   100,
			Payload: stringPortion,
			Incr:    false,
		})

		if index > 0 {
			// Add section of phrase to suggestions
			stringPortion = strings.Join(wordArray[index:limit], " ")
			suggestion = append(suggestion, redisearch.Suggestion{
				Term:    stringPortion,
				Score:   100,
				Payload: stringPortion,
				Incr:    false,
			})
		}

		// Add single word as suggestion
		suggestion = append(suggestion, redisearch.Suggestion{
			Term:    word,
			Score:   100,
			Payload: word,
			Incr:    false,
		})
	}

	return suggestion
}
