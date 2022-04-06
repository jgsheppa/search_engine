package controllers

import (
	"encoding/json"
	"fmt"
	"github.com/RediSearch/redisearch-go/redisearch"
	"github.com/gomodule/redigo/redis"
	"github.com/gorilla/mux"
	"github.com/jgsheppa/search_engine/models"
	"io"
	"io/ioutil"
	"log"
	"net/http"
)

var (
	ErrorNotFound = "No matches found"
)

func NewArticle(redisClient redis.Conn, redisSearch redisearch.Client, autoCompleter redisearch.Autocompleter) *RedisDB {
	return &RedisDB{
		redisClient:   redisClient,
		redisSearch:   redisSearch,
		autoCompleter: autoCompleter,
	}
}

type RedisDB struct {
	// Redis client
	redisClient redis.Conn
	// Redisearch client
	redisSearch redisearch.Client
	// Redis autocomplete
	autoCompleter redisearch.Autocompleter
} // @name RedisDB

type Error struct {
	// Error Message
	Error string `json:"error"`
	// Error Code
	Code int `json:"code"`
} // @name Error

type Field struct {
	// Field name
	Name string `json:"name"`
	// Field type
	Type string `json:"type"`
} // @name Field

// PostDocuments godoc
// @Summary Post documents to Redisearch
// @Tags Document
// @Param Body body models.Articles true "The body to create a Redis document for an article"
// @Success 201 {object} models.Articles
// @Failure 422
// @Router /api/documents [post]
func (rdb *RedisDB) PostDocuments(w http.ResponseWriter, r *http.Request) {
	var articles models.Articles
	body, err := ioutil.ReadAll(io.LimitReader(r.Body, 1048576))
	if err != nil {
		panic(err)
	}
	if err := r.Body.Close(); err != nil {
		panic(err)
	}
	if err := json.Unmarshal(body, &articles); err != nil {
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		w.WriteHeader(422) // unprocessable entity
		if err := json.NewEncoder(w).Encode(err); err != nil {
			panic(err)
		}
	}

	models.CreateDocument(rdb.redisSearch, rdb.autoCompleter, articles)

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusCreated)
	if err := json.NewEncoder(w).Encode(articles); err != nil {
		panic(err)
	}
	fmt.Fprintln(w, "Document successfully uploaded")
}

// DeleteDocument godoc
// @Summary Delete documents from Redisearch
// @Tags Document
// @Param documentName path string true "search term"
// @ID documentName
// @Success 200 {string} string "Ok"
// @Router /api/document/delete/{documentName} [delete]
func (rdb *RedisDB) DeleteDocument(w http.ResponseWriter, r *http.Request) {
	term := mux.Vars(r)["documentName"]

	err := models.DeleteDocument(rdb.redisSearch, term)
	if err != nil {
		panic(err)
	}
	w.WriteHeader(http.StatusOK)
	fmt.Fprintln(w, http.StatusOK, "Document successfully deleted")
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

func SearchAndSuggest(w http.ResponseWriter, rdb *RedisDB, term string, sortBy string) {
	var HighlightedFields []string
	HighlightedFields = append(HighlightedFields, "title")
	HighlightedFields = append(HighlightedFields, "author")

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
