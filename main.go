package main

import (
	"fmt"
	"github.com/RediSearch/redisearch-go/redisearch"
	"github.com/jgsheppa/search_engine/controllers"
	redis_conn "github.com/jgsheppa/search_engine/redis"
	"github.com/spf13/viper"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

func NotFound(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusNotFound)
	fmt.Fprintf(w, "Route not found")
}

type App struct {
	Router *mux.Router
}

func main() {
	viper.SetConfigName("config")
	viper.AddConfigPath(".")
	err := viper.ReadInConfig()

	if err != nil {
		log.Fatalf("Error while reading config file %s", err)
	}

	var pool = redis_conn.NewPool()

	r := mux.NewRouter()
	port := "3000"

	client := pool.Get()
	defer client.Close()

	// Create a client. By default a client is schemaless
	// unless a schema is provided when creating the index

	c := redisearch.NewClientFromPool(pool, "bpArticles")
	auto := redisearch.NewAutocompleterFromPool(pool, "bpArticles")

	terms := redisearch.Suggestion{Term: "prod", Score: 1, Payload: "product", Incr: true}
	auto.AddTerms(terms)
	//c.Drop()

	// Create a schema
	sc := redisearch.NewSchema(redisearch.DefaultOptions).
		AddField(redisearch.NewNumericField("date")).
		AddField(redisearch.NewNumericField("id")).
		AddField(redisearch.NewTextFieldOptions("author", redisearch.TextFieldOptions{Weight: 5.0, Sortable: true})).
		AddField(redisearch.NewTextFieldOptions("title", redisearch.TextFieldOptions{Weight: 5.0, Sortable: true})).
		AddField(redisearch.NewTextFieldOptions("url", redisearch.TextFieldOptions{Weight: 5.0, Sortable: true}))

	// Create the index with the given schema
	err = c.CreateIndex(sc)
	if err != nil {
		fmt.Println("Index already exists")
	}

	searchController := controllers.NewArticle(client, *c)

	r.HandleFunc("/documents", searchController.PostDocuments).Methods("POST")
	r.HandleFunc("/documents/{documentName}", searchController.DeleteDocument).Methods("DELETE")
	r.HandleFunc("/search/{term}", searchController.Search).Methods("GET")

	// HandlerFunc converts notFound to the correct type
	r.NotFoundHandler = http.HandlerFunc(NotFound)

	fmt.Println("Starting the development server on port" + port)
	http.ListenAndServe(":"+port, r)
}
