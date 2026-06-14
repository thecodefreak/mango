package server

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	chimiddleware "github.com/go-chi/chi/v5/middleware"
	"github.com/thecodefreak/mango/internal/config"
	"github.com/thecodefreak/mango/internal/handlers"
	"github.com/thecodefreak/mango/internal/server/middleware"
)

func NewRouter(c *config.ServerConfig) http.Handler {
	r := chi.NewRouter()
	r.Use(chimiddleware.RequestID)
	r.Use(chimiddleware.Logger)
	r.Use(chimiddleware.Recoverer)

	sh := handlers.NewStaticPageHandler(c.DocumentRoot)

	r.With(middleware.ValidateToken(c.ServerToken)).Route("/api", func(r chi.Router) {
		r.Post("/static-page", sh.CreateOrUpdate)
	})

	r.Get("/*", sh.Get)

	// API Routes

	return r
}
