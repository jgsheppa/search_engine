package controllers

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/RediSearch/redisearch-go/redisearch"
	"github.com/gomodule/redigo/redis"
	"github.com/gorilla/mux"
)

func NewArticle(redisClient redis.Conn, redisSearch redisearch.Client) *Article {
	return &Article{
		redisClient: redisClient,
		redisSearch: redisSearch,
	}
}

type Article struct {
	redisClient redis.Conn 
	redisSearch redisearch.Client
}

func(a *Article) CreateArticle(w http.ResponseWriter, r *http.Request) {
	document := mux.Vars(r)["document"]
	title := mux.Vars(r)["title"]

	doc := redisearch.NewDocument(document, 1.0)
	doc.Set("title", title).
		Set("body", "foo bar").
		Set("date", time.Now().Unix())

	// Index the document. The API accepts multiple documents at a time
	if err := a.redisSearch.Index([]redisearch.Document{doc}...); err != nil {
		log.Fatal(err)
		fmt.Fprintf(w, err.Error())
	}
	fmt.Println("DOC", doc)


	err := json.NewEncoder(w).Encode(doc.Properties)
	if err != nil {
		log.Println("failed to encode response")
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func(a *Article) Search(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	term := mux.Vars(r)["term"]

	// Searching with limit and sorting
	docs, _, err := a.redisSearch.Search(redisearch.NewQuery(term).
	Limit(0, 5).
	SetReturnFields("title"))

	if err != nil {
		panic(err)
	}

	fmt.Println("DOCS", docs)

	response := []map[string]interface{}{}
	for _, doc := range docs {
		response = append(response, doc.Properties)
	}

	err = json.NewEncoder(w).Encode(response)
	if err != nil {
		log.Println("failed to encode response")
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func(a *Article) Delete(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	_, err := a.redisClient.Do("SET", "mykey", "Delete in Redis!")
	if err != nil {
			panic(err)
	}

	value, err := a.redisClient.Do("DEL", "mykey")
	if err != nil {
			panic(err)
	}

	fmt.Fprintf(w, "%s \n", value)
}