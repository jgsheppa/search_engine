package models

import (
	"testing"
)

var expectedResult = [12]string{"New Topic", "New", "Topic"}

func TestCreateSuggestions(t *testing.T) {

	article := Document{
		Name: "New Topic",
	}

	t.Run("suggestion output", func(t *testing.T) {
		got := CreateSuggestions(article)
		want := expectedResult
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

func TestSuggestionFactory(t *testing.T) {
	words := []string{"On", "Pair", "Programming"}

	t.Run("suggestion output", func(t *testing.T) {
		got := SuggestionFactory(words)
		want := []string{"On Pair", "On", "On Pair Programming", "Pair Programming",
			"Pair", "Programming"}

		for i, suggestion := range got {
			if suggestion.Term != want[i] {
				t.Errorf("got error %q want %q", suggestion.Term, want[i])
			}
		}
	})
}
