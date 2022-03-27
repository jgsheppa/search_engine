package main

import (
	"fmt"
	"net/http"

	"github.com/RediSearch/redisearch-go/redisearch"
	"github.com/gorilla/mux"
	"github.com/jgsheppa/search_engine/controllers"
	"github.com/jgsheppa/search_engine/redis"
)

func NotFound(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusNotFound)
	fmt.Fprintf(w, "Route not found")
}

var pool = redis.NewPool()

type App struct {
	Router *mux.Router
}

func main() {

	r := mux.NewRouter()
	port := "3000"

	client := pool.Get()
	defer client.Close()

	// Create a client. By default a client is schemaless
	// unless a schema is provided when creating the index
	c := redisearch.NewClientFromPool(pool, "todo_index2")


	// Create a schema
	sc := redisearch.NewSchema(redisearch.DefaultOptions).
	AddField(redisearch.NewNumericField("date")).
	AddField(redisearch.NewNumericField("id")).
	AddField(redisearch.NewTextFieldOptions("name", redisearch.TextFieldOptions{Weight: 5.0, Sortable: true}))

	c.Drop()

	// Create the index with the given schema
	err := c.CreateIndex(sc) 
	if err != nil {
		fmt.Println("Index already exists")
	}

	searchController := controllers.NewArticle(client, *c)

	r.HandleFunc("/documents", searchController.PostDocuments).Methods("POST")
	r.HandleFunc("/documents/{document}", searchController.DeleteDocument).Methods("POST")
	r.HandleFunc("/field", searchController.PostField).Methods("POST")
	r.HandleFunc("/search/{term}", searchController.Search).Methods("GET")

	// HandlerFunc converts notFound to the correct type
	r.NotFoundHandler = http.HandlerFunc(NotFound)

	fmt.Println("Starting the development server on port" + port)
	http.ListenAndServe(":" + port, r)
}