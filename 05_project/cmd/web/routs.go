package main

import (
	"net/http"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"

)

func (app *Config) routes() http.Handler {
	// create router
	mux := chi.NewRouter()

	// setup middleware
	mux.Use(middleware.Recoverer)

	// define application routs
	mux.Get("/", app.HomePage)
	return mux
}