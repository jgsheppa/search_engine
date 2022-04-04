package main

import (
	"fmt"
	"github.com/RediSearch/redisearch-go/redisearch"
	"github.com/gorilla/mux"
	"github.com/jgsheppa/search_engine/controllers"
	redis_conn "github.com/jgsheppa/search_engine/redis"
	"github.com/spf13/viper"
	"log"
	"net/http"
	"os"
)

func NotFound(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusNotFound)
	fmt.Fprintf(w, "Route not found")
}

func main() {
	port := "3000"
	if os.Getenv("IS_PROD") == "true" {
		viper.AutomaticEnv()
		port = os.Getenv("PORT")
	} else {
		config := "config"
		viper.SetConfigName(config)
		viper.AddConfigPath(".")
		err := viper.ReadInConfig()

		if err != nil {
			log.Fatalf("Error while reading config file %s", err)
		}
	}

	var pool = redis_conn.NewPool()

	r := mux.NewRouter()

	client := pool.Get()
	defer client.Close()

	// Create a client. By default a client is schemaless
	// unless a schema is provided when creating the index

	c := redisearch.NewClientFromPool(pool, "bpArticles")
	autocomplete := redisearch.NewAutocompleterFromPool(pool, "bpArticles")
	//c.Drop()

	// Create a schema
	sc := CreateSchema()

	// Create the index with the given schema
	err := c.CreateIndex(sc)
	if err != nil {
		fmt.Println("Index already exists")
	}

	searchController := controllers.NewArticle(client, *c, *autocomplete)

	r.HandleFunc("/", controllers.DisplayAPIRoutes).Methods("GET")
	r.HandleFunc("/documents", searchController.PostDocuments).Methods("POST")
	r.HandleFunc("/document/delete/{documentName}", searchController.DeleteDocument).Methods("DELETE")
	r.HandleFunc("/search/{term}", searchController.Search).Methods("GET")
	r.HandleFunc("/search/{term}/{sortBy}", searchController.SearchAndSort).Methods("GET")

	// HandlerFunc converts notFound to the correct type
	r.NotFoundHandler = http.HandlerFunc(NotFound)

	fmt.Println("Starting the development server on port" + port)
	http.ListenAndServe(":"+port, r)
}

// CreateSchema is used to create the schema for your Redisearch documents,
// which will allow you to add your data in the form of these documents
func CreateSchema() *redisearch.Schema {
	sc := redisearch.NewSchema(redisearch.DefaultOptions).
		AddField(redisearch.NewNumericField("date")).
		AddField(redisearch.NewNumericField("id")).
		AddField(redisearch.NewTextFieldOptions("author", redisearch.TextFieldOptions{Weight: 5.0, Sortable: true})).
		AddField(redisearch.NewTextFieldOptions("title", redisearch.TextFieldOptions{Weight: 5.0, Sortable: true})).
		AddField(redisearch.NewTextFieldOptions("url", redisearch.TextFieldOptions{Weight: 5.0, Sortable: true}))
	return sc
}
