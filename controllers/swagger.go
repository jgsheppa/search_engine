package controllers

import "net/http"

func DisplayAPIRoutes(w http.ResponseWriter, r *http.Request) {
	http.Redirect(w, r, "/swagger/index.html", http.StatusFound)
}
