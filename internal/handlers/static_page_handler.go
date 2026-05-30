package handlers

import (
	"net/http"
	"os"
	"strings"

	"github.com/go-chi/chi/v5"
)

type StaticPageHandler struct{}

func NewStaticPageHandler() *StaticPageHandler {
	return &StaticPageHandler{}
}

func (h *StaticPageHandler) Get(w http.ResponseWriter, r *http.Request) {
	path := chi.URLParam(r, "*")
	pagePath := "storage/static-pages/" + path

	if _, err := os.Stat(pagePath); err == nil {

		if strings.HasSuffix(pagePath, "/") {
			pagePath += "index.html"
		}

		http.ServeFile(w, r, pagePath)
	} else {
		http.NotFound(w, r)
	}
}
