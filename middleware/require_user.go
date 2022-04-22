package middleware

import (
	"encoding/json"
	"github.com/jgsheppa/search_engine/context"
	"github.com/jgsheppa/search_engine/models"
	"net/http"
)

type RequireUser struct {
	models.UserService
}

// ApplyFn assumes that user has already been run,
// otherwise it will not work correctly
func (mw *RequireUser) ApplyFn(next http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		user := context.User(r.Context())
		if user == nil {
			json.NewEncoder(w).Encode(models.AuthError)
			w.WriteHeader(models.AuthError.HttpStatus)
			return
		}
		next(w, r)
	})
}

// Apply assumes that user has already been run,
// otherwise it will not work correctly
func (mw *RequireUser) Apply(next http.Handler) http.HandlerFunc {
	return mw.ApplyFn(next.ServeHTTP)
}
