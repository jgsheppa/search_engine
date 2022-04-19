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

	sortBy := "city"
	if len(sort) > 0 {
		sortBy = sort
	}

	isAscending := true
	if ascending == "false" {
		isAscending = false
	}

	highlighted := []string{"city"}

	result, err := rdb.s.SearchAndSuggest(isAscending, queryLimit, highlighted, term, sortBy)
	if err != nil {
		json.NewEncoder(w).Encode(err)
		json.NewEncoder(w).Encode(models.ValidationError)
		return
	}
	json.NewEncoder(w).Encode(result)
	return
}

// GeoSearch godoc
// @Summary Search Redisearch documents
// @Tags Search
// @ID country geo search
// @Param longitude query string true "Longitude"
// @Param latitude query string true "Latitude"
// @Param radius query string true "Radius"
// @Param limit query int false "Limit number of results"
// @Produce json
// @Success 200 {object} swagger.SwaggerSearchResponse "Ok"
// @Failure 404 {object} models.ApiError "Not Found"
// @Failure 500 {object} models.ApiError "Server Error"
// @Router /api/search/geo [get]
func (rdb *RedisDB) GeoSearch(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	longitude := r.FormValue("longitude")
	latitude := r.FormValue("latitude")
	radius := r.FormValue("radius")
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

	docs := rdb.s.GeoSearch(longitude, latitude, radius, queryLimit)

	json.NewEncoder(w).Encode(docs)
	return
}
