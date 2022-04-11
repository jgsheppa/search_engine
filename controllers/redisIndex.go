package controllers

import (
	"fmt"
	"github.com/jgsheppa/search_engine/models"
	"net/http"
)

// DropIndex godoc
// @Summary Delete all documents from Redisearch
// @Tags Index
// @Success 200 {string} string "Ok"
// @Router /api/index/delete/ [delete]
func (rdb *RedisDB) DropIndex(w http.ResponseWriter, r *http.Request) {
	err := rdb.redisSearch.Drop()
	if err != nil {
		http.Error(w, "Cannot drop index - it does not exist", http.StatusNotFound)
		return
	}
	w.WriteHeader(http.StatusOK)
	fmt.Fprintln(w, http.StatusOK, "Index successfully deleted")
}

// CreateIndex godoc
// @Summary Delete all documents from Redisearch
// @Tags Index
// @Success 200 {string} string "Ok"
// @Router /api/index/create/articles [POST]
func (rdb *RedisDB) CreateIndex(w http.ResponseWriter, r *http.Request) {
	pool := rdb.Pool
	_, _, err := models.CreateIndex(&pool, "articles")
	if err != nil {
		http.Error(w, "Index already exists", http.StatusNotFound)
		return
	}
	w.WriteHeader(http.StatusOK)
	fmt.Fprintln(w, http.StatusOK, "Index successfully created")
}

// CreateIndex godoc
// @Summary Delete all documents from Redisearch
// @Tags Index
// @Success 200 {string} string "Ok"
// @Router /api/index/create/ [POST]
func (rdb *RedisDB) CreateGuideIndex(w http.ResponseWriter, r *http.Request) {
	pool := rdb.Pool
	_, _, err := models.CreateIndex(&pool, "guide")
	if err != nil {
		http.Error(w, "Index already exists", http.StatusNotFound)
		return
	}
	w.WriteHeader(http.StatusOK)
	fmt.Fprintln(w, http.StatusOK, "Index successfully created")
}
