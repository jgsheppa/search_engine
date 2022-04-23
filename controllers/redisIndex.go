package controllers

import (
	"encoding/json"
	"github.com/jgsheppa/search_engine/models"
	"net/http"
)

// DropIndex godoc
// @Summary Delete all documents from Redisearch
// @Tags Index
// @Success 200 {string} string "Ok"
// @Success 401 {object} models.ApiError "Unauthorized"
// @Router /api/index/delete [delete]
func (rdb *RedisDB) DropIndex(w http.ResponseWriter, r *http.Request) {
	err := rdb.redisSearch.Drop()
	if err != nil {
		http.Error(w, "Cannot drop index - it does not exist", http.StatusNotFound)
		return
	}

	response := Response{
		Message: "Index successfully deleted",
	}
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}

// CreateIndex godoc
// @Summary Create Redis index for BestPracticer guides
// @Tags Index
// @Success 200 {string} string "Ok"
// @Success 401 {object} models.ApiError "Unauthorized"
// @Router /api/index/create [POST]
func (rdb *RedisDB) CreateIndex(w http.ResponseWriter, r *http.Request) {
	pool := rdb.Pool
	_, _, err := models.CreateIndex(&pool, "index")
	if err != nil {
		http.Error(w, "Index already exists", http.StatusNotFound)
		return
	}

	response := Response{
		Message: "Index successfully created",
	}
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(response)
}
