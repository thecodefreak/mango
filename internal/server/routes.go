package server

import (
	"net/http"

	"github.com/go-chi/chi/v5"

	chimiddleware "github.com/go-chi/chi/v5/middleware"
	"github.com/thecodefreak/mango/internal/handlers"
)

func NewRouter() http.Handler {
	r := chi.NewRouter()
	r.Use(chimiddleware.RequestID)
	r.Use(chimiddleware.Logger)
	r.Use(chimiddleware.Recoverer)

	r.Route("/api", func(r chi.Router) {
		
	})

	h := handlers.NewStaticPageHandler()
	r.Get("/*", h.Get)

	// API Routes

	return r
}