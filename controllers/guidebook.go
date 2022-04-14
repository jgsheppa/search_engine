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
// @Param htmlTag path string true "HTML tag to scrape"
// @Param containerClass path string true "Class of container holding HTML tag"
// @Success 201 {object} models.Guides
// @Failure 422
// @Router /api/document/guide/{url}/{htmlTag}/{containerClass} [post]
func (rdb *RedisDB) PostGuideDocuments(w http.ResponseWriter, r *http.Request) {
	url := mux.Vars(r)["url"]
	htmlTag := mux.Vars(r)["htmlTag"]
	containerClass := mux.Vars(r)["containerClass"]
	guides := scraper.ScrapeWebPage(url, htmlTag, containerClass)

	err := models.CreateGuideDocument(rdb.redisSearch, rdb.autoCompleter, guides)
	if err != nil {
		json.NewEncoder(w).Encode(validationError)
	}

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusCreated)
	if err := json.NewEncoder(w).Encode(guides); err != nil {
		json.NewEncoder(w).Encode(validationError)
	}
	fmt.Fprintln(w, "Document successfully uploaded")
}
