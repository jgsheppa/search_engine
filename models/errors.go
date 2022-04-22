package models

import (
	"encoding/json"
	"net/http"
	"strings"
)

const (
	ErrNotFound              modelError   = "models: resource not found"
	ErrIDInvalid             privateError = "models: ID provided was invalid"
	ErrPasswordIncorrect     modelError   = "models: incorrect password"
	ErrRememberTokenTooShort privateError = "models: remember token must be at least 32 bytes"
	ErrRememberRequired      privateError = "models: remember token required"
	ErrEmailRequired         modelError   = "Email address is required"
	ErrEmailInvalid          modelError   = "Email address is not valid"
	ErrEmailTaken            modelError   = "models: email address is already taken"
	ErrPasswordMinLength     modelError   = "Password must be 8 characters long"
	ErrPasswordRequired      modelError   = "Password is required"
)

type modelError string

func (e modelError) Error() string {
	return string(e)
}

func (e modelError) Public() string {
	s := strings.Replace(string(e), "models: ", "", 1)
	split := strings.Split(s, " ")
	split[0] = strings.Title(split[0])
	return strings.Join(split, " ")
}

type privateError string

func (e privateError) Error() string {
	return string(e)
}

type ApiError struct {
	HttpStatus  int    `json:"-"`
	Title       string `json:"title"`
	Description string `json:"description"`
}

func NotFound(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusNotFound)
	json.NewEncoder(w).Encode(NotFoundError)
}

// Make our struct implement the error interface
func (e *ApiError) Error() string {
	return e.Description
}

var serverError = ApiError{
	HttpStatus:  http.StatusInternalServerError,
	Title:       "Server error",
	Description: "We're sorry, something went wrong on our side",
}

var AuthError = ApiError{
	HttpStatus:  http.StatusUnauthorized,
	Title:       "Unauthorized",
	Description: "Authentication required",
}

var ValidationError = ApiError{
	HttpStatus:  http.StatusUnprocessableEntity,
	Title:       "Validation error",
	Description: "Please refer to the swagger documentation for the correct input data format",
}

var urlParamError = ApiError{
	HttpStatus:  http.StatusBadRequest,
	Title:       "Wrong parameters",
	Description: "Please enter correct parameters or check the swagger documentation",
}

var NotFoundError = ApiError{
	HttpStatus:  http.StatusNotFound,
	Title:       "Not found",
	Description: "That resource doesn't exist",
}

var LargePayloadError = ApiError{
	HttpStatus:  http.StatusRequestEntityTooLarge,
	Title:       "Payload too large",
	Description: "The payload is too large - please reduce its size.",
}
