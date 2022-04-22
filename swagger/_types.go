package swagger

type SwaggerSuggestion struct {
	Term    string  `json:"term" example:"Pair"`
	Score   float64 `json:"score" example:"70.70"`
	Payload string  `json:"payload" example:"Pair"`
	Incr    bool    `json:"incr" example:"false"`
}

type SwaggerResponse struct {
	// Author of article
	Link string `json:"link" example:"https://www.nba.com"`
	// URL of doc
	Name string `json:"name" example:"Michael Jordan"`
	// Topics of article
	Active string `json:"active" example:"1985-2003"`
}

type SwaggerSearchResponse struct {
	Suggestion []SwaggerSuggestion `json:"suggestion"`
	Response   []SwaggerResponse   `json:"response"`
	Total      int                 `json:"total" example:"1"`
}
