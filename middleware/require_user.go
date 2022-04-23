package middleware

import (
	"encoding/json"
	"github.com/jgsheppa/search_engine/context"
	"github.com/jgsheppa/search_engine/models"
	"net/http"
	"strings"
)

type User struct {
	models.UserService
}

func (mw *User) Apply(next http.Handler) http.HandlerFunc {
	return mw.ApplyFn(next.ServeHTTP)
}

func (mw *User) ApplyFn(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		auth := r.Header.Get("Authorization")
		bearer := strings.Split(auth, " ")
		var token string

		if len(bearer) > 1 {
			token = bearer[1]

			user, _ := mw.UserService.ByRemember(token)

			ctx := r.Context()
			ctx = context.WithUser(ctx, user)
			r = r.WithContext(ctx)
			next(w, r)
		} else {
			// if the user is logged in then pass the user for the navbar
			cookie, err := r.Cookie("remember_token")
			if err != nil {
				json.NewEncoder(w).Encode(models.AuthError)
				return
			}
			user, err := mw.UserService.ByRemember(cookie.Value)
			if err != nil {
				json.NewEncoder(w).Encode(models.AuthError)
				return
			}

			ctx := r.Context()
			ctx = context.WithUser(ctx, user)
			r = r.WithContext(ctx)
			next(w, r)
		}
	}
}
