package controllers

type Route struct {
	Route   string `json:"route"`
	Method  string `json:"method"`
	Note    string `json:"note"`
	Example string `json:"example"`
}

type Routes struct {
	Index         Route `json:"index"`
	Search        Route `json:"search"`
	SearchAndSort Route `json:"search_and_sort"`
	PostDocument  Route `json:"post_document"`
	Delete        Route `json:"delete"`
}

func CreateRouteMap() Routes {
	index := Route{
		Route:  "/",
		Method: "GET",
		Note:   "Map of the API endpoints",
	}

	search := Route{
		Route:  "/search/{term}",
		Method: "GET",
		Note:   "Endpoint to search for articles by term or author. For this demo, try 'product', 'pro', or 'Sherif'",
	}

	searchAndSort := Route{
		Route:  "/search/{term}/{sortBy}",
		Method: "GET",
		Note:   "Similar to the search endpoint, but with the ability to sort by field. Ex: /search/engineer/author",
	}

	postDoc := Route{
		Route:   "/documents",
		Method:  "POST",
		Note:    "Endpoint for posting multiple documents to Redis store.",
		Example: "curl -H \"Content-Type: application/json\" -d '[{\"document\": \"todo6\", \"author\":\"Birgitta Böckeler\", \"title\": \"On Pair Programming\", \"url\": \"https://martinfowler.com/articles/on-pair-programming.html\", \"id\": 6}]' https://bp-search-engine.herokuapp.com/documents",
	}

	delete := Route{
		Route:   "document/delete/{document}",
		Method:  "DELETE",
		Note:    "Endpoint to delete documents.",
		Example: "curl -X \"DELETE\" https://bp-search-engine.herokuapp.com/document/delete/todo6",
	}

	routes := Routes{
		Index:         index,
		Search:        search,
		SearchAndSort: searchAndSort,
		PostDocument:  postDoc,
		Delete:        delete,
	}

	return routes
}
