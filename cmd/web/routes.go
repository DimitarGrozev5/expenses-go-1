package main

import (
	"net/http"

	"github.com/dimitargrozev5/expenses-go-1/pkg/config"
	"github.com/dimitargrozev5/expenses-go-1/pkg/handlers"
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
)

func routes(app *config.AppConfig) http.Handler {
	mux := chi.NewRouter()

	mux.Use(middleware.Recoverer)
	mux.Use(NoSurf)
	mux.Use(SessionLoad)

	mux.Get("/", handlers.Repo.Home)
	mux.Get("/about", handlers.Repo.About)
	mux.Get("/expenses", handlers.Repo.Expenses)
	mux.Get("/static/*", handlers.Repo.Static)

	return mux
}
