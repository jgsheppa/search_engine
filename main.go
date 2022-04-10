package main

import (
	"fmt"
	"github.com/gorilla/mux"
	"github.com/jgsheppa/search_engine/controllers"
	"github.com/jgsheppa/search_engine/models"
	redis_conn "github.com/jgsheppa/search_engine/redis"
	"github.com/jgsheppa/search_engine/swagger"
	"github.com/spf13/viper"
	"log"
	"net/http"
	"os"

	httpSwagger "github.com/swaggo/http-swagger"
)

// @title BestPracticer Search Engine
// @version 1.0
// @description This is a search engine built with Redisearch

// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html

// @host localhost:3000
// @BasePath /
// @query.collection.format multi
// @schemes http https
func main() {
	port, baseUrl, httpProtocol := SetupConfig()

	var pool = redis_conn.NewPool()

	client := pool.Get()
	defer client.Close()

	//c, autocomplete, err := models.CreateIndex(pool)
	//if err != nil {
	//	log.Println(err)
	//}

	c, autocomplete, err := models.CreateIndex(pool, "guide")
	if err != nil {
		log.Println(err)
	}
	//_, _, err = models.CreateIndex(pool, "guide")
	//if err != nil {
	//	log.Println(err)
	//}

	swagger.SwaggerInfo.Title = "BestPracticer Search Engine"
	swagger.SwaggerInfo.Description = "This is a search engine built with Redisearch"
	swagger.SwaggerInfo.Version = "1.0"
	swagger.SwaggerInfo.Host = baseUrl
	swagger.SwaggerInfo.BasePath = "/"
	swagger.SwaggerInfo.Schemes = []string{httpProtocol}

	r := mux.NewRouter()
	searchController := controllers.NewArticle(*pool, *c, *autocomplete)

	// Swagger routes
	r.PathPrefix("/swagger/").Handler(httpSwagger.WrapHandler)
	r.PathPrefix("/swagger/").Handler(http.StripPrefix("/swagger", http.FileServer(http.Dir("swagger"))))
	r.HandleFunc("/", controllers.DisplayAPIRoutes).Methods("GET")
	// Document routes
	r.HandleFunc("/api/document/article", searchController.PostDocuments).Methods("POST")
	r.HandleFunc("/api/document/guide/{url}", searchController.PostGuideDocuments).Methods("POST")
	r.HandleFunc("/api/document/delete/{documentName}", searchController.DeleteDocument).Methods("DELETE")
	// Index routes
	r.HandleFunc("/api/index/delete/", searchController.DropIndex).Methods("DELETE")
	r.HandleFunc("/api/index/create/", searchController.CreateIndex).Methods("POST")
	// Search routes
	r.HandleFunc("/api/search/{term}/{service}", searchController.Search).Methods("GET")
	r.HandleFunc("/api/search/{term}/{service}", searchController.Search).
		Queries("limit", "{limit}").
		Methods("GET")
	r.HandleFunc("/api/search/{term}/{service}", searchController.Search).
		Queries("ascending", "{ascending}").
		Methods("GET")
	r.HandleFunc("/api/search/{term}/{service}", searchController.Search).
		Queries("sort", "{sort}").
		Methods("GET")
	r.HandleFunc("/api/search/{term}/{service}", searchController.Search).
		Queries("sort", "{sort}").
		Queries("ascending", "{ascending}").Methods("GET")
	r.HandleFunc("/api/search/{term}/{service}", searchController.Search).
		Queries("sort", "{sort}").
		Queries("limit", "{limit}").Methods("GET")
	r.HandleFunc("/api/search/{term}/{service}", searchController.Search).
		Queries("limit", "{limit}").
		Queries("ascending", "{ascending}").Methods("GET")
	r.HandleFunc("/api/search/{term}/{service}", searchController.Search).
		Queries("sort", "{sort}").
		Queries("ascending", "{ascending}").
		Queries("limit", "{limit}").
		Methods("GET")

	// HandlerFunc converts notFound to the correct type
	r.NotFoundHandler = http.HandlerFunc(controllers.NotFound)

	fmt.Println("Starting the development server on port " + port)
	err = http.ListenAndServe(":"+port, r)
	if err != nil {
		log.Fatal(err)
	}

}

func SetupConfig() (portNum, baseURL, http string) {
	port := "3000"
	baseUrl := "localhost:3000"
	httpProtocol := "http"

	if os.Getenv("IS_PROD") == "true" {
		viper.AutomaticEnv()
		port = os.Getenv("PORT")
		baseUrl = os.Getenv("BASE_URL")
		httpProtocol = "https"
		return port, baseUrl, httpProtocol
	}
	config := "config"
	viper.SetConfigName(config)
	viper.AddConfigPath(".")
	err := viper.ReadInConfig()

	if err != nil {
		log.Fatalf("Error while reading config file %s", err)
	}
	return port, baseUrl, httpProtocol
}
