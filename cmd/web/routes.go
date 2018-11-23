package main

import (
	"net/http"

	"github.com/gorilla/mux"
)

func (app *application) routes() http.Handler {
	r := mux.NewRouter()
	r.HandleFunc("/", app.home).Methods("GET")

	// static files, js, css...
	staticDir := "./ui/static/"
	r.PathPrefix("/static/").
		Handler(http.StripPrefix("/static/", http.FileServer(http.Dir(staticDir))))

	return r
}
