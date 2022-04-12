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
	"net/http"
)

func NewArticle(pool redis.Pool, redisSearch redisearch.Client, autoCompleter redisearch.Autocompleter) *RedisDB {
	return &RedisDB{
		Pool:          pool,
		redisSearch:   redisSearch,
		autoCompleter: autoCompleter,
	}
}

type RedisDB struct {
	Pool redis.Pool
	// Redisearch client
	redisSearch redisearch.Client
	// Redis autocomplete
	autoCompleter redisearch.Autocompleter
} // @name RedisDB

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
// @Failure 422 {object} ApiError
// @Router /api/document/article [post]
func (rdb *RedisDB) PostDocuments(w http.ResponseWriter, r *http.Request) {
	var articles models.Articles
	body, err := ioutil.ReadAll(io.LimitReader(r.Body, 1048576))
	if err != nil {
		json.NewEncoder(w).Encode(largePayloadError)

	}
	if err := r.Body.Close(); err != nil {
		json.NewEncoder(w).Encode(serverError)
	}
	if err := json.Unmarshal(body, &articles); err != nil {
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		w.WriteHeader(422) // unprocessable entity
		if err := json.NewEncoder(w).Encode(validationError); err != nil {
			fmt.Fprintln(w, validationError.Description)
		}
	}

	err = models.CreateDocument(rdb.redisSearch, rdb.autoCompleter, articles)
	if err != nil {
		json.NewEncoder(w).Encode(validationError)
	}

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusCreated)
	if err := json.NewEncoder(w).Encode(articles); err != nil {
		json.NewEncoder(w).Encode(validationError)
	}
	fmt.Fprintln(w, "Document successfully uploaded")
}

// DeleteDocument godoc
// @Summary Delete documents from Redisearch
// @Tags Document
// @Param documentName path string true "search term"
// @ID documentName
// @Success 200 {string} string "Ok"
// @Failure 404 {object} ApiError
// @Router /api/document/delete/{documentName} [delete]
func (rdb *RedisDB) DeleteDocument(w http.ResponseWriter, r *http.Request) {
	term := mux.Vars(r)["documentName"]

	err := models.DeleteDocument(rdb.redisSearch, term)
	if err != nil {
		json.NewEncoder(w).Encode(notFoundError)
		return
	}
	w.WriteHeader(http.StatusOK)
	fmt.Fprintln(w, http.StatusOK, "Document successfully deleted")
}
