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
	redisClient   redis.Conn
	redisSearch   redisearch.Client
	autoCompleter redisearch.Autocompleter
}

type Error struct {
	Error string `json:"error"`
	Code  int    `json:"code"`
}

type Field struct {
	Name string `json:"name"`
	Type string `json:"type"`
}

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

func (rdb *RedisDB) DeleteDocument(w http.ResponseWriter, r *http.Request) {
	term := mux.Vars(r)["documentName"]

	err := models.DeleteDocument(rdb.redisSearch, term)
	if err != nil {
		panic(err)
	}
	fmt.Fprintln(w, "Document successfully deleted")
}

func (rdb *RedisDB) Search(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	term := mux.Vars(r)["term"]
	sortBy := "date"

	SearchAndSuggest(w, rdb, term, sortBy)
}

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
			fmt.Fprintln(w, "No suggestions found")
			return
		}

		if len(auto) > 0 {
			docs, total, err := rdb.redisSearch.Search(redisearch.NewQuery(auto[0].Payload).
				Limit(0, 5).
				Highlight(HighlightedFields, "<b>", "</b>").
				SetSortBy(sortBy, true))

			if err != nil {
				fmt.Fprintln(w, "No query matches these suggestions")
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
