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

type Todo struct {
	Document string `json:"document"`
	Author   string `json:"author"`
	ID       int    `json:"id"`
	URL      string `json:"url"`
	Title    string `json:"title"`
}

type Field struct {
	Name string `json:"name"`
	Type string `json:"type"`
}

type Todos []Todo

func (rdb *RedisDB) PostDocuments(w http.ResponseWriter, r *http.Request) {
	var todos models.Todos
	body, err := ioutil.ReadAll(io.LimitReader(r.Body, 1048576))
	if err != nil {
		panic(err)
	}
	if err := r.Body.Close(); err != nil {
		panic(err)
	}
	if err := json.Unmarshal(body, &todos); err != nil {
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		w.WriteHeader(422) // unprocessable entity
		if err := json.NewEncoder(w).Encode(err); err != nil {
			panic(err)
		}
	}

	models.CreateDocument(rdb.redisSearch, rdb.autoCompleter, todos)

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusCreated)
	if err := json.NewEncoder(w).Encode(todos); err != nil {
		panic(err)
	}

}

func (rdb *RedisDB) DeleteDocument(w http.ResponseWriter, r *http.Request) {
	term := mux.Vars(r)["documentName"]

	err := models.DeleteDocument(rdb.redisSearch, term)
	if err != nil {
		panic(err)
	}
	fmt.Fprintln(w, "Document successfully deleted")
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

func (rdb *RedisDB) Search(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	term := mux.Vars(r)["term"]

	var HighlightedFields []string
	HighlightedFields = append(HighlightedFields, "title")

	// Searching with limit and sorting
	docs, total, err := rdb.redisSearch.Search(redisearch.NewQuery(term).
		Limit(0, 5).Highlight(HighlightedFields, "<b>", "</b>"))
	if err != nil {
		panic(err)
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
		auto, err = rdb.autoCompleter.SuggestOpts(term, opts)
		if err != nil {
			panic(err)
		}

		docs, total, err = rdb.redisSearch.Search(redisearch.NewQuery(auto[0].Payload).
			Limit(0, 5).Highlight(HighlightedFields, "<b>", "</b>"))
		if err != nil {
			panic(err)
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
