package swagger

type SwaggerSuggestion struct {
	Term    string  `json:"term" example:"Pair"`
	Score   float64 `json:"score" example:"70.70"`
	Payload string  `json:"payload" example:"Pair"`
	Incr    bool    `json:"incr" example:"false"`
}

type SwaggerResponse struct {
	// Author of article
	Text string `json:"author" example:"Goal Setting"`
	// URL of doc
	URL string `json:"url" example:"www.guidebook.bestpracticer.com/goal-setting"`
	// Topics of article
	Topic string `json:"topic" example:"Goal Setting"`
	Date  int    `json:"date" example:"1649762803"`
}

type SwaggerSearchResponse struct {
	Suggestion []SwaggerSuggestion `json:"suggestion"`
	Response   []SwaggerResponse   `json:"response"`
	Total      int                 `json:"total" example:"1"`
}
