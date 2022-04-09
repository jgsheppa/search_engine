package controllers

import (
	"encoding/json"
	"net/http"
)

type ApiError struct {
	HttpStatus  int    `json:"-"`
	Title       string `json:"title"`
	Description string `json:"description"`
}

func NotFound(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusNotFound)
	json.NewEncoder(w).Encode(notFoundError)
}

// Make our struct implement the error interface
func (e *ApiError) Error() string {
	return e.Description
}

var serverError = ApiError{
	HttpStatus:  500,
	Title:       "Server error",
	Description: "We're sorry, something went wrong on our side",
}

var validationError = ApiError{
	HttpStatus:  422,
	Title:       "Validation error",
	Description: "Please refer to the swagger documentation for the correct input data format",
}

var urlParamError = ApiError{
	HttpStatus:  400,
	Title:       "Wrong parameters",
	Description: "Please enter correct parameters or check the swagger documentation",
}

var notFoundError = ApiError{
	HttpStatus:  404,
	Title:       "Not found",
	Description: "That resource doesn't exist",
}

var largePayloadError = ApiError{
	HttpStatus:  413,
	Title:       "Payload too large",
	Description: "The payload is too large - please reduce its size.",
}
