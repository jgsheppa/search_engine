package models

import (
	"github.com/RediSearch/redisearch-go/redisearch"
	"testing"
)

func TestCreateSuggestions(t *testing.T) {
	words := [10]string{"How", "to", "lose", "a", "guy", "in", "ten", "days", "The", "Dude"}

	article := Todo{
		Document: "article1",
		Author:   "The Dude",
		ID:       1,
		URL:      "www.npr.org",
		Title:    "How to lose a guy in ten days",
	}

	t.Run("known word", func(t *testing.T) {
		var wantedSuggestions []redisearch.Suggestion
		for _, word := range words {
			helper := HelperSuggestions(word)
			wantedSuggestions = append(wantedSuggestions, helper)
		}

		got := CreateSuggestions(article)
		want := wantedSuggestions

		for i, suggestion := range got {
			if suggestion.Term != want[i].Term {
				t.Errorf("got error %q want %q", suggestion.Term, want[i].Term)
			}

			if suggestion.Payload != want[i].Payload {
				t.Errorf("got error %q want %q", suggestion.Payload, want[i].Payload)
			}
		}
	})
}

func HelperSuggestions(word string) redisearch.Suggestion {
	suggestion := redisearch.Suggestion{
		Term:    word,
		Score:   100,
		Payload: word,
		Incr:    false,
	}
	return suggestion
}
