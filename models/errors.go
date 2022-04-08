package models

import "strings"

const (
	ErrCreateIndex           modelError   = "Unable to create index. Hint: it likely exists."
	ErrDeleteIndex           modelError   = "Index successfully deleted"
	ErrRememberTokenTooShort privateError = "models: remember token must be at least 32 bytes"
	ErrRememberRequired      privateError = "models: remember token required"
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
