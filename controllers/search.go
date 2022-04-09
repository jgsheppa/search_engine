package controllers

import (
	"encoding/json"
	"fmt"
	"github.com/RediSearch/redisearch-go/redisearch"
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"strconv"
)

type SearchResponse struct {
	Total      int                      `json:"total"`
	Response   []map[string]interface{} `json:"response"`
	Suggestion []redisearch.Suggestion  `json:"suggestion"`
}

type swaggerResponse struct {
	// Author of article
	Author string `json:"author" example:"Anne Applebaum"`
	ID     int    `json:"id" example:"1"`
	// URL of article
	URL string `json:"url" example:"www.bestpractice.com/awesome-article"`
	// Title of article
	Title string `json:"title" example:"How to be awesome"`
	// Topics of article
	Topic string `json:"topic" example:"Awesome Stuff"`
}

type swaggerSuggestion struct {
	Term    string  `json:"term" example:"Pair"`
	Score   float64 `json:"score" example:"70.70"`
	Payload string  `json:"payload" example:"Pair"`
	Incr    bool    `json:"incr" example:"false"`
}

type SwaggerSearchResponse struct {
	Suggestion []swaggerSuggestion `json:"suggestion"`
	Response   []swaggerResponse   `json:"response"`
	Total      int                 `json:"total" example:"1"`
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
// @Param term path string true "Search by keyword"
// @Param sort query string false "Sort by field"
// @Param ascending query string false "Accepted terms: true, false"
// @Param limit query int false "Limit number of results"
// @Produce json
// @Success 200 {object} SwaggerSearchResponse "Ok"
// @Failure 404 {string} string "Not Found"
// @Failure 500 {string} string "Server Error"
// @Router /api/search/{term} [get]
func (rdb *RedisDB) Search(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	term := mux.Vars(r)["term"]
	sort := r.FormValue("sort")
	ascending := r.FormValue("ascending")
	limit := r.FormValue("limit")

	queryLimit := 5
	if len(limit) > 0 {
		limitAsInt, err := strconv.Atoi(limit)
		if err != nil {
			fmt.Fprintln(w, "Error with limit query")
		}
		queryLimit = limitAsInt
	}

	sortBy := "date"
	if len(sort) > 0 {
		sortBy = sort
	}

	isAscending := true
	if ascending == "false" {
		isAscending = false
	}

	SearchAndSuggest(w, rdb, isAscending, queryLimit, term, sortBy)
}

func SearchAndSuggest(w http.ResponseWriter, rdb *RedisDB, order bool, limit int, term, sortBy string) {
	var HighlightedFields = []string{"title", "author", "topic"}

	// Searching with limit and sorting
	docs, total, err := rdb.redisSearch.Search(redisearch.NewQuery(term).
		Limit(0, limit).
		Highlight(HighlightedFields, "<b>", "</b>").
		SetSortBy(sortBy, order))

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
