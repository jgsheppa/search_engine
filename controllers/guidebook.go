package controllers

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/jgsheppa/search_engine/models"
	"github.com/jgsheppa/search_engine/scraper"
	"net/http"
)

// PostGuideDocuments godoc
// @Summary Post documents to Redisearch
// @Tags Document
// @Param url path string true "URL path of guidebook"
// @Success 201 {object} models.Guides
// @Failure 422
// @Router /api/document/guide/{url} [post]
func (rdb *RedisDB) PostGuideDocuments(w http.ResponseWriter, r *http.Request) {
	url := mux.Vars(r)["url"]
	guides := scraper.ScrapeGuidebookPages(url)

	models.CreateGuideDocument(rdb.redisSearch, rdb.autoCompleter, guides)

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusCreated)
	if err := json.NewEncoder(w).Encode(guides); err != nil {
		json.NewEncoder(w).Encode(validationError)
	}
	fmt.Fprintln(w, "Document successfully uploaded")
}
