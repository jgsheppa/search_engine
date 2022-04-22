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

func NewDocuments(redisService models.Services) *RedisDB {
	return &RedisDB{
		Pool:          redisService.Pool,
		redisSearch:   redisService.Redisearch,
		autoCompleter: redisService.Autocomplete,
		s:             redisService,
	}
}

type RedisDB struct {
	Pool redis.Pool
	// Redisearch client
	redisSearch redisearch.Client
	// Redis autocomplete
	autoCompleter redisearch.Autocompleter
	s             models.Services
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
// @Param Body body models.Documents true "The body to create a Redis document for an article"
// @Success 201 {object} models.Documents
// @Success 401 {object} models.ApiError "Unauthorized"
// @Failure 422 {object} models.ApiError
// @Router /api/document [post]
func (rdb *RedisDB) PostDocuments(w http.ResponseWriter, r *http.Request) {
	var documents models.Documents
	body, err := ioutil.ReadAll(io.LimitReader(r.Body, 1048576))
	if err != nil {
		json.NewEncoder(w).Encode(models.LargePayloadError)
		w.WriteHeader(models.LargePayloadError.HttpStatus)
		return
	}
	if err := r.Body.Close(); err != nil {
		json.NewEncoder(w).Encode(err)
		return
	}
	if err := json.Unmarshal(body, &documents); err != nil {
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		w.WriteHeader(422) // unprocessable entity
		json.NewEncoder(w).Encode(1)
		json.NewEncoder(w).Encode(models.ValidationError)
		w.WriteHeader(models.ValidationError.HttpStatus)
		return
	}

	err = rdb.s.CreateDocument(documents)
	if err != nil {
		json.NewEncoder(w).Encode(err)
	}

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	if err := json.NewEncoder(w).Encode(documents); err != nil {
		json.NewEncoder(w).Encode(models.ValidationError)
		w.WriteHeader(models.ValidationError.HttpStatus)
		return
	}
	fmt.Fprintln(w, "Document successfully uploaded")
}

// DeleteDocument godoc
// @Summary Delete documents from Redisearch
// @Tags Document
// @Param documentName path string true "search term"
// @ID documentName
// @Success 200 {string} string "Ok"
// @Success 401 {object} models.ApiError "Unauthorized"
// @Failure 404 {object} models.ApiError
// @Router /api/document/delete/{documentName} [delete]
func (rdb *RedisDB) DeleteDocument(w http.ResponseWriter, r *http.Request) {
	term := mux.Vars(r)["documentName"]

	err := rdb.s.DeleteDocument(term)
	if err != nil {
		json.NewEncoder(w).Encode(models.NotFoundError)
		w.WriteHeader(models.NotFoundError.HttpStatus)
		return
	}
	w.WriteHeader(http.StatusOK)
	fmt.Fprintln(w, http.StatusOK, "Document successfully deleted")
}
