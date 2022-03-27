package controllers

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
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

type Todo struct {
	Document string `json:"document"`
	Name string `json:"name"`
	ID int `json:"id"`
	Person string `json:"person"`
	Age int `json:"age"`

}

type Field struct {
	Name string `json:"name"`
	Type string `json:"type"`
}

type Todos []Todo

func(a *Article) PostDocuments(w http.ResponseWriter, r *http.Request) {
	var todos Todos
	body, err := ioutil.ReadAll(io.LimitReader(r.Body, 1048576))
	if err != nil {
		panic(err)
	}
	if err := r.Body.Close(); err != nil {
		panic(err)
	}
	if err := json.Unmarshal(body, &todos); err != nil {
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		w.WriteHeader(422) // unprocessable entity
		if err := json.NewEncoder(w).Encode(err); err != nil {
			panic(err)
		}
	}

	index, err := a.redisSearch.Info()
	if err != nil {
		panic(err)
	}
	fmt.Println("index!", index)
	
	
		a.CreateDocument(todos)
	
	
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusCreated)
	if err := json.NewEncoder(w).Encode(todos); err != nil {
			panic(err)
	}
	
}

func(a *Article) DeleteDocument(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	term := mux.Vars(r)["document"]
	err := a.redisSearch.DeleteDocument(term)
	if err != nil {
			panic(err)
	}
	fmt.Fprintln(w, "Document successfully deleted")
}

func(a *Article) PostField(w http.ResponseWriter, r *http.Request) {
	var fields []Field
	body, err := ioutil.ReadAll(io.LimitReader(r.Body, 1048576))
	if err != nil {
		panic(err)
	}
	if err := r.Body.Close(); err != nil {
		panic(err)
	}
	if err := json.Unmarshal(body, &fields); err != nil {
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		w.WriteHeader(422) // unprocessable entity
		if err := json.NewEncoder(w).Encode(err); err != nil {
			panic(err)
		}
	}

	for _, field := range fields {
		a.AddField(field)
	}
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusCreated)
	if err := json.NewEncoder(w).Encode(fields); err != nil {
			panic(err)
	}
}

// TODO: Create an endpoint to add a field where
// the user can dictate the field options: name, weight, etc.
func(a *Article) AddField(field Field) {
	switch field.Type {
	case "text":
		a.redisSearch.AddField(redisearch.NewTextFieldOptions(field.Name, redisearch.TextFieldOptions{Weight: 5.0, Sortable: true}))
	case "number":
		a.redisSearch.AddField(redisearch.NewNumericField(field.Name))
	default:
		fmt.Println("No Redis field associated with this type")
	}
	
}


func(a *Article) CreateDocument(todos Todos) {
	var documents []redisearch.Document

	for _, todo := range todos {
		doc := redisearch.NewDocument(todo.Document, 1.0)
		doc.Set("name", todo.Name).
			Set("id", todo.ID).
			Set("date", time.Now().Unix())
		documents = append(documents, doc)
	}
		
	// Index the document. The API accepts multiple documents at a time
	if err := a.redisSearch.Index(documents...); err != nil {
		log.Fatal(err)
	}
	
}


// TODO: Highlighting text would be really cool
func(a *Article) Search(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	term := mux.Vars(r)["term"]

	// Searching with limit and sorting
	docs, _, err := a.redisSearch.Search(redisearch.NewQuery(term).
	Limit(0, 5))

	if err != nil {
		panic(err)
	}

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

