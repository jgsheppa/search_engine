package controllers

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"github.com/jgsheppa/search_engine/models"
	"net/http"
	"strconv"
)

type SuggestOptions struct {
	Num          int
	Fuzzy        bool
	WithPayloads bool
	WithScores   bool
}

// Search godoc
// @Summary Search Redisearch documents
// @Tags Search
// @ID article search
// @Param term path string true "Search by keyword"
// @Param sort query string false "Sort by field"
// @Param ascending query boolean false "Ascending?"
// @Param limit query int false "Limit number of results"
// @Param geo query boolean false "Is it a geo query"
// @Produce json
// @Success 200 {object} swagger.SwaggerSearchResponse "Ok"
// @Failure 404 {object} models.ApiError "Not Found"
// @Failure 500 {object} models.ApiError "Server Error"
// @Router /api/search/{term} [get]
func (rdb *RedisDB) Search(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	//Allow CORS here By * or specific origin
	w.Header().Set("Access-Control-Allow-Origin", "http://localhost:3000")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type,access-control-allow-origin, access-control-allow-headers")

	term := mux.Vars(r)["term"]
	sort := r.FormValue("sort")
	ascending := r.FormValue("ascending")
	limit := r.FormValue("limit")

	queryLimit := 5
	if len(limit) > 0 {
		limitAsInt, err := strconv.Atoi(limit)
		if err != nil {
			json.NewEncoder(w).Encode(err)
			json.NewEncoder(w).Encode(models.ValidationError)
			return
		}
		queryLimit = limitAsInt
	}

	sortBy := "name"
	if len(sort) > 0 {
		sortBy = sort
	}

	isAscending := true
	if ascending == "false" {
		isAscending = false
	}

	highlighted := []string{"name"}

	result, err := rdb.s.SearchAndSuggest(isAscending, queryLimit, highlighted, term, sortBy)
	if err != nil {
		json.NewEncoder(w).Encode(err)
		json.NewEncoder(w).Encode(models.ValidationError)
		return
	}
	json.NewEncoder(w).Encode(result)
	return
}
