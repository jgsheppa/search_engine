package models

import (
	"github.com/RediSearch/redisearch-go/redisearch"
)

type SearchResponse struct {
	Total      int                      `json:"total"`
	Response   []map[string]interface{} `json:"response"`
	Suggestion []redisearch.Suggestion  `json:"suggestion"`
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
	// Searching with limit and sorting
	docs, total, err := s.Redisearch.Search(redisearch.NewQuery(term).
		Limit(0, limit).
		Highlight(highlightedFields, "<b>", "</b>").
		SetSortBy(sortBy, order))
	if err != nil {
		return SearchResponse{}, err
	}

	var response []map[string]interface{}
	for _, doc := range docs {
		response = append(response, doc.Properties)
	}

	var auto []redisearch.Suggestion

	if len(response) == 0 {
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
			docs, total, err := s.Redisearch.Search(redisearch.NewQuery(auto[0].Payload).
				Limit(0, limit).
				Highlight(highlightedFields, "<b>", "</b>").
				SetSortBy(sortBy, order))

			if err != nil {
				return SearchResponse{}, err
			}

			for _, doc := range docs {
				response = append(response, doc.Properties)
			}

			result := SearchResponse{
				total,
				response,
				auto,
			}

			return result, nil
		} else {
			return SearchResponse{}, err
		}

	} else {
		result := SearchResponse{
			total,
			response,
			auto,
		}

		return result, nil
	}

}

func (s *Services) GeoSearch() []redisearch.Document {
	// Todo: Get document by name and look up location that way
	// Searching for 100KM radius should only output Catania
	docs, _, _ := s.Redisearch.Search(redisearch.NewQuery("*").AddFilter(
		redisearch.Filter{
			Field: "location",
			Options: redisearch.GeoFilterOptions{
				Lon:    48.20849,
				Lat:    16.37208,
				Radius: 20,
				Unit:   redisearch.KILOMETERS,
			},
		},
	).
		Limit(0, 30).
		SetSortBy("city", true))

	return docs
}
