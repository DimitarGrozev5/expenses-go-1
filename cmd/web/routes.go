package main

import (
	"net/http"

	"github.com/dimitargrozev5/expenses-go-1/internal/config"
	"github.com/dimitargrozev5/expenses-go-1/internal/handlers"
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
)

func routes(_ *config.AppConfig) http.Handler {
	mux := chi.NewRouter()

	mux.Use(middleware.Recoverer)
	mux.Use(NoSurf)
	mux.Use(SessionLoad)

	// mux.Get("/", handlers.Repo.Home)
	// mux.Get("/about", handlers.Repo.About)
	// mux.Get("/expenses", handlers.Repo.Expenses)
	// mux.Post("/expenses", handlers.Repo.PostExpenses)

	// Public routes
	mux.Group(func(r chi.Router) {
		// Handle home page
		r.Get("/", handlers.Repo.Home)

		// Handle login page
		r.Post("/login", handlers.Repo.PostLogin)

		// Serve access to static files
		r.Get("/static/*", handlers.Repo.Static)
	})

	// Private routes
	mux.Group(func(r chi.Router) {
		// Set auth middleware
		r.Use(IsAuth)

		// Handle logout
		r.Get("/logout", handlers.Repo.Logout)

		// Handle expense related routes
		r.Get("/expenses", handlers.Repo.Expenses)
		r.Post("/expenses/add", handlers.Repo.PostNewExpense)
		r.Post("/expenses/{expenseId}/edit", handlers.Repo.PostEditExpense)
		r.Post("/expenses/{expenseId}/delete", handlers.Repo.PostDeleteExpense)

		// Handle accounts related routes
		r.Get("/accounts", handlers.Repo.Accounts)
		r.Post("/accounts/add", handlers.Repo.PostNewAccount)
	})

	return mux
}
