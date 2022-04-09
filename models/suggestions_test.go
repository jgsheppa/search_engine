package models

import (
	"testing"
)

func TestCreateSuggestions(t *testing.T) {
	words := [9]string{"On Pair", "On", "On Pair Programming", "Pair Programming",
		"Pair", "The Dude", "The"}

	article := Article{
		Document: "article1",
		Author:   "The Dude",
		ID:       1,
		URL:      "www.npr.org",
		Title:    "On Pair Programming",
	}

	t.Run("suggestion output", func(t *testing.T) {
		got := CreateSuggestions(article)
		want := words

		for i, suggestion := range got {
			if suggestion.Term != want[i] {
				t.Errorf("got error %q want %q", suggestion.Term, want[i])
			}

			if suggestion.Payload != want[i] {
				t.Errorf("got error %q want %q", suggestion.Payload, want[i])
			}
		}
	})
}
