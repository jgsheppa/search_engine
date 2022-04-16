package controllers

import (
	"encoding/json"
	"fmt"
	"github.com/RediSearch/redisearch-go/redisearch"
	"github.com/gomodule/redigo/redis"
	"github.com/gorilla/mux"
	"github.com/jgsheppa/search_engine/models"
	"github.com/jgsheppa/search_engine/scraper"
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
// @Failure 422 {object} models.ApiError
// @Router /api/document [post]
func (rdb *RedisDB) PostDocuments(w http.ResponseWriter, r *http.Request) {
	var documents models.Documents
	body, err := ioutil.ReadAll(io.LimitReader(r.Body, 1048576))
	if err != nil {
		json.NewEncoder(w).Encode(models.LargePayloadError)
	}
	if err := r.Body.Close(); err != nil {
		json.NewEncoder(w).Encode(err)
	}
	if err := json.Unmarshal(body, &documents); err != nil {
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		w.WriteHeader(422) // unprocessable entity
		if err := json.NewEncoder(w).Encode(models.ValidationError); err != nil {
			fmt.Fprintln(w, models.ValidationError)
		}
	}

	err = rdb.s.CreateDocument(documents)
	if err != nil {
		json.NewEncoder(w).Encode(err)
	}

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusCreated)
	if err := json.NewEncoder(w).Encode(documents); err != nil {
		json.NewEncoder(w).Encode(models.ValidationError)
	}
	fmt.Fprintln(w, "Document successfully uploaded")
}

// PostScrapedDocuments godoc
// @Summary Post documents to Redisearch
// @Tags Document
// @Param url path string true "URL path of guidebook"
// @Param htmlTag path string true "HTML tag to scrape"
// @Param containerClass path string true "Class of container holding HTML tag"
// @Success 201 {object} models.Documents
// @Failure 422 {object} models.ApiError
// @Router /api/document/{url}/{htmlTag}/{containerClass} [post]
func (rdb *RedisDB) PostScrapedDocuments(w http.ResponseWriter, r *http.Request) {
	url := mux.Vars(r)["url"]
	htmlTag := mux.Vars(r)["htmlTag"]
	containerClass := mux.Vars(r)["containerClass"]
	documents := scraper.ScrapeWebPage(url, htmlTag, containerClass)

	err := rdb.s.CreateDocument(documents)
	if err != nil {
		json.NewEncoder(w).Encode(err)
	}

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusCreated)
	if err := json.NewEncoder(w).Encode(documents); err != nil {
		json.NewEncoder(w).Encode(models.ValidationError)
	}
	fmt.Fprintln(w, "Document successfully uploaded")
}

// DeleteDocument godoc
// @Summary Delete documents from Redisearch
// @Tags Document
// @Param documentName path string true "search term"
// @ID documentName
// @Success 200 {string} string "Ok"
// @Failure 404 {object} models.ApiError
// @Router /api/document/delete/{documentName} [delete]
func (rdb *RedisDB) DeleteDocument(w http.ResponseWriter, r *http.Request) {
	term := mux.Vars(r)["documentName"]

	err := rdb.s.DeleteDocument(term)
	if err != nil {
		json.NewEncoder(w).Encode(models.NotFoundError)
		return
	}
	w.WriteHeader(http.StatusOK)
	fmt.Fprintln(w, http.StatusOK, "Document successfully deleted")
}
