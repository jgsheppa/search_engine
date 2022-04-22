package models

import (
	"github.com/RediSearch/redisearch-go/redisearch"
)

type SearchResponse struct {
	Total      int                     `json:"total"`
	Response   []redisearch.Document   `json:"response"`
	Suggestion []redisearch.Suggestion `json:"suggestion"`
}

type SearchService interface {
	SearchAndSuggest(order bool,
		limit int,
		highlightedFields []string,
		term,
		sortBy string) (SearchResponse, error)
}

func (s *Services) SearchAndSuggest(
	order bool,
	limit int,
	highlightedFields []string,
	term,
	sortBy string,
) (SearchResponse, error) {
	totalHits := 0
	// Searching with limit and sorting
	docs, total, err := s.Redisearch.Search(redisearch.NewQuery(term).
		Limit(0, limit).
		Highlight(highlightedFields, "<b>", "</b>").
		SetSortBy(sortBy, order))
	if err != nil {
		return SearchResponse{}, err
	}

	var response []redisearch.Document

	for _, doc := range docs {
		response = append(response, doc)
	}

	var auto []redisearch.Suggestion

	if len(response) == 0 {
		totalHits := 0

		opts := redisearch.SuggestOptions{
			Num:          5,
			Fuzzy:        false,
			WithPayloads: true,
			WithScores:   true,
		}
		auto, err := s.Autocomplete.SuggestOpts(term, opts)
		if err != nil {
			return SearchResponse{}, err
		}

		if len(auto) > 0 {

			for _, suggestion := range auto {
				docs, total, err = s.Redisearch.Search(redisearch.NewQuery(suggestion.Payload).
					Limit(0, limit).
					Highlight(highlightedFields, "<b>", "</b>").
					SetSortBy(sortBy, order))

				if err != nil {
					return SearchResponse{}, err
				}

				for _, doc := range docs {
					response = append(response, doc)
				}
				totalHits += total
			}

			result := SearchResponse{
				totalHits,
				response,
				auto,
			}

			return result, nil
		} else {
			return SearchResponse{}, err
		}

	} else {
		result := SearchResponse{
			totalHits,
			response,
			auto,
		}

		return result, nil
	}

}
