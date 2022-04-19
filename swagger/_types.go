package swagger

type SwaggerSuggestion struct {
	Term    string  `json:"term" example:"Pair"`
	Score   float64 `json:"score" example:"70.70"`
	Payload string  `json:"payload" example:"Pair"`
	Incr    bool    `json:"incr" example:"false"`
}

type SwaggerResponse struct {
	// Author of article
	Country string `json:"country" example:"Austria"`
	// URL of doc
	Name string `json:"name" example:"Wien"`
	// Topics of article
	Location string `json:"location" example:"48.20849,16.37208"`
}

type SwaggerSearchResponse struct {
	Suggestion []SwaggerSuggestion `json:"suggestion"`
	Response   []SwaggerResponse   `json:"response"`
	Total      int                 `json:"total" example:"1"`
}
