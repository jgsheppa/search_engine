package models

import (
	"github.com/RediSearch/redisearch-go/redisearch"
	"strings"
)

func CreateSuggestions(document Document) []redisearch.Suggestion {
	var suggestion []redisearch.Suggestion

	splitCountry := strings.Split(document.Country, " ")
	text := SuggestionFactory(splitCountry)
	suggestion = append(suggestion, text...)

	splitCity := strings.Split(document.Name, " ")
	topic := SuggestionFactory(splitCity)
	suggestion = append(suggestion, topic...)

	return suggestion
}

func SuggestionFactory(wordArray []string) []redisearch.Suggestion {
	var suggestion []redisearch.Suggestion

	var score float64
	score = 100
	incr := false

	for index, word := range wordArray {
		wordArrayLength := len(wordArray) - 1
		increment := 2
		limit := increment + index

		if index == wordArrayLength {
			// Add single word as suggestion
			suggestion = append(suggestion, redisearch.Suggestion{
				Term:    word,
				Score:   score,
				Payload: word,
				Incr:    incr,
			})

			return suggestion
		}

		// Add beginning of phrase to limit as suggestion
		stringPortion := strings.Join(wordArray[:limit], " ")
		suggestion = append(suggestion, redisearch.Suggestion{
			Term:    stringPortion,
			Score:   score,
			Payload: stringPortion,
			Incr:    incr,
		})

		if index > 0 {
			// Add section of phrase to suggestions
			stringPortion = strings.Join(wordArray[index:limit], " ")
			suggestion = append(suggestion, redisearch.Suggestion{
				Term:    stringPortion,
				Score:   score,
				Payload: stringPortion,
				Incr:    incr,
			})
		}

		// Add single word as suggestion
		suggestion = append(suggestion, redisearch.Suggestion{
			Term:    word,
			Score:   score,
			Payload: word,
			Incr:    incr,
		})
	}

	return suggestion
}
