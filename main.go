package main

import (
	"fmt"
	"github.com/gorilla/mux"
	"github.com/jgsheppa/search_engine/controllers"
	"github.com/jgsheppa/search_engine/middleware"
	"github.com/jgsheppa/search_engine/models"
	"github.com/jgsheppa/search_engine/swagger"
	"github.com/spf13/viper"
	"log"
	"net/http"
	"os"

	httpSwagger "github.com/swaggo/http-swagger"
)

// @title Redisearch Search Engine
// @version 1.0
// @description This is a search engine built with Redisearch

// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html

// @host localhost:3001
// @BasePath /
// @query.collection.format multi
// @schemes http https

// @securityDefinitions.apikey ApiKeyAuth
// @in header
// @name Authorization
func main() {
	port, baseUrl, psqlInfo, httpProtocol := SetupConfig()

	services, err := models.NewServices(psqlInfo)
	if err != nil {
		panic(err)
	}

	services.DestructiveReset()
	services.AutoMigrate()

	emailAddress := viper.Get("ADMIN_EMAIL").(string)
	name := viper.Get("ADMIN_NAME").(string)
	password := viper.Get("ADMIN_PW").(string)

	user := models.User{
		Email:    emailAddress,
		Name:     name,
		Password: password,
		Role:     "admin",
	}

	services.User.Create(&user)

	swagger.SwaggerInfo.Title = "Redisearch Search Engine"
	swagger.SwaggerInfo.Description = "This is a search engine built with Redisearch"
	swagger.SwaggerInfo.Version = "1.0"
	swagger.SwaggerInfo.Host = baseUrl
	swagger.SwaggerInfo.BasePath = "/"
	swagger.SwaggerInfo.Schemes = []string{httpProtocol}

	userMw := middleware.User{
		UserService: services.User,
	}

	r := mux.NewRouter()
	searchController := controllers.NewDocuments(*services)
	authController := controllers.Auth(services.User)

	// Auth
	r.HandleFunc("/api/auth/login", authController.Login).Methods("POST")
	r.HandleFunc("/api/auth/logout", userMw.ApplyFn(authController.Logout)).Methods("POST")
	// Swagger routes
	r.PathPrefix("/swagger/").Handler(httpSwagger.WrapHandler)
	r.PathPrefix("/swagger/").Handler(http.StripPrefix("/swagger", http.FileServer(http.Dir("swagger"))))
	r.HandleFunc("/", controllers.DisplayAPIRoutes).Methods("GET")
	// Document routes
	r.HandleFunc("/api/document", userMw.ApplyFn(searchController.PostDocuments)).Methods("POST")
	r.HandleFunc("/api/document/delete/{documentName}", userMw.ApplyFn(searchController.DeleteDocument)).
		Methods("DELETE")
	// Index routes
	r.HandleFunc("/api/index/delete", userMw.ApplyFn(searchController.DropIndex)).Methods("DELETE")
	r.HandleFunc("/api/index/create", userMw.ApplyFn(searchController.CreateIndex)).Methods("POST")

	// Search routes
	r.HandleFunc("/api/search/{term}", searchController.Search).Methods("GET", "OPTIONS")
	r.HandleFunc("/api/search/{term}", searchController.Search).
		Queries("limit", "{limit}").
		Methods("GET", "OPTIONS")
	r.HandleFunc("/api/search/{term}", searchController.Search).
		Queries("ascending", "{ascending}").
		Methods("GET", "OPTIONS")
	r.HandleFunc("/api/search/{term}", searchController.Search).
		Queries("sort", "{sort}").
		Methods("GET", "OPTIONS")
	r.HandleFunc("/api/search/{term}", searchController.Search).
		Queries("sort", "{sort}").
		Queries("ascending", "{ascending}").Methods("GET", "OPTIONS")
	r.HandleFunc("/api/search/{term}", searchController.Search).
		Queries("sort", "{sort}").
		Queries("limit", "{limit}").Methods("GET", "OPTIONS")
	r.HandleFunc("/api/search/{term}", searchController.Search).
		Queries("limit", "{limit}").
		Queries("ascending", "{ascending}").Methods("GET", "OPTIONS")
	r.HandleFunc("/api/search/{term}", searchController.Search).
		Queries("sort", "{sort}").
		Queries("ascending", "{ascending}").
		Queries("limit", "{limit}").
		Methods("GET", "OPTIONS")

	// HandlerFunc converts notFound to the correct type
	r.NotFoundHandler = http.HandlerFunc(models.NotFound)

	fmt.Println("Starting the development server on port " + port)
	err = http.ListenAndServe(":"+port, r)
	if err != nil {
		log.Fatal(err)
	}

}

func SetupConfig() (portNum, baseURL, psqlInfo, http string) {
	port := "3001"
	baseUrl := "localhost:3001"
	httpProtocol := "http"

	if os.Getenv("IS_PROD") == "true" {
		viper.AutomaticEnv()
		port = os.Getenv("PORT")
		baseUrl = os.Getenv("BASE_URL")
		psqlInfo = os.Getenv("DATABASE_URL")
		httpProtocol = "https"
		return port, baseUrl, psqlInfo, httpProtocol
	}
	config := "config"
	viper.SetConfigName(config)
	viper.AddConfigPath(".")
	err := viper.ReadInConfig()
	if err != nil {
		log.Fatalf("Error while reading config file %s", err)
	}
	psqlInfo = viper.Get("DATABASE_URL").(string)

	return port, baseUrl, psqlInfo, httpProtocol
}
