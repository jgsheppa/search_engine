package controllers

import (
	"encoding/json"
	"fmt"
	"github.com/RediSearch/redisearch-go/redisearch"
	"github.com/gorilla/mux"
	"log"
	"net/http"
)

type SearchResponse struct {
	Total      int                      `json:"total"`
	Response   []map[string]interface{} `json:"response"`
	Suggestion []redisearch.Suggestion  `json:"suggestion"`
}

type SuggestOptions struct {
	Num          int
	Fuzzy        bool
	WithPayloads bool
	WithScores   bool
}

// Search godoc
// @Summary Search Redisearch documents
// @Tags Search
// @ID term
// @Param term path string true "search term"
// @Success 200 {string} string "Ok"
// @Failure 404 {string} string "Not Found"
// @Failure 500 {string} string "Server Error"
// @Router /api/search/{term} [get]
func (rdb *RedisDB) Search(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	term := mux.Vars(r)["term"]
	sortBy := "date"

	SearchAndSuggest(w, rdb, term, sortBy)
}

// SearchAndSort godoc
// @Summary Search and sort Redisearch documents
// @Tags Search
// @Param term path string true "search term"
// @Param sortBy path string true "sort by"
// @Success 200 {string} string "Ok"
// @Failure 404 {string} string "Not Found"
// @Failure 500 {string} string "Server Error"
// @Router /api/search/{term}/{sortBy} [get]
func (rdb *RedisDB) SearchAndSort(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	term := mux.Vars(r)["term"]
	sortBy := mux.Vars(r)["sortBy"]

	SearchAndSuggest(w, rdb, term, sortBy)
}

func SearchAndSuggest(w http.ResponseWriter, rdb *RedisDB, term string, sortBy string) {
	var HighlightedFields []string
	HighlightedFields = append(HighlightedFields, "title")
	HighlightedFields = append(HighlightedFields, "author")
	HighlightedFields = append(HighlightedFields, "topic")

	// Searching with limit and sorting
	docs, total, err := rdb.redisSearch.Search(redisearch.NewQuery(term).
		Limit(0, 5).
		Highlight(HighlightedFields, "<b>", "</b>").
		SetSortBy(sortBy, true))

	if err != nil {
		log.Println(err)
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
		auto, err := rdb.autoCompleter.SuggestOpts(term, opts)
		if err != nil {
			w.WriteHeader(http.StatusNotFound)
			fmt.Fprintln(w, http.StatusNotFound, "No suggestions found")
			return
		}

		if len(auto) > 0 {
			docs, total, err := rdb.redisSearch.Search(redisearch.NewQuery(auto[0].Payload).
				Limit(0, 5).
				Highlight(HighlightedFields, "<b>", "</b>").
				SetSortBy(sortBy, true))

			if err != nil {
				w.WriteHeader(http.StatusNotFound)
				fmt.Fprintln(w, http.StatusNotFound, "No query matches these suggestions")
				return
			}

			for _, doc := range docs {
				response = append(response, doc.Properties)
			}

			result := SearchResponse{
				total,
				response,
				auto,
			}

			err = json.NewEncoder(w).Encode(result)
		} else {
			error := Error{
				ErrorNotFound,
				404,
			}
			err = json.NewEncoder(w).Encode(error)
		}

	} else {
		result := SearchResponse{
			total,
			response,
			auto,
		}

		err = json.NewEncoder(w).Encode(result)
	}
	if err != nil {
		log.Println("failed to encode response")
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
