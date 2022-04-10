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

func CreateGuideSuggestions(guide Guide) []redisearch.Suggestion {
	var suggestion []redisearch.Suggestion

	splitTitle := strings.Split(guide.Text, " ")
	title := SuggestionFactory(splitTitle)
	suggestion = append(suggestion, title...)

	splitTopic := strings.Split(guide.Topic, " ")
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
