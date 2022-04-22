package middleware

import (
	"encoding/json"
	"fmt"
	"github.com/jgsheppa/search_engine/context"
	"github.com/jgsheppa/search_engine/models"
	"net/http"
)

type User struct {
	models.UserService
}

func (mw *User) Apply(next http.Handler) http.HandlerFunc {
	return mw.ApplyFn(next.ServeHTTP)
}

func (mw *User) ApplyFn(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// if the user is logged in then pass the user for the navbar
		cookie, err := r.Cookie("remember_token")
		if err != nil {
			next(w, r)
			return
		}
		user, err := mw.UserService.ByRemember(cookie.Value)
		if err != nil {
			next(w, r)
			return
		}
		fmt.Println("user mw", user)

		ctx := r.Context()
		ctx = context.WithUser(ctx, user)
		r = r.WithContext(ctx)
		next(w, r)
	}
}

// Assumes that user has already been run,
// otherwise it will not work correctly
type RequireUser struct {
	User
}

// Assumes that user has already been run,
// otherwise it will not work correctly
func (mw *RequireUser) ApplyFn(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		user := context.User(r.Context())
		if user == nil {
			json.NewEncoder(w).Encode(models.AuthError)
			return
		}
		next(w, r)
	}
}

// Assumes that user has already been run,
// otherwise it will not work correctly
func (mw *RequireUser) Apply(next http.Handler) http.HandlerFunc {
	return mw.ApplyFn(next.ServeHTTP)
}
