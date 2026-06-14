package handlers

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/go-chi/chi/v5"
	"github.com/thecodefreak/mango/internal/helpers"
)

var storagePath = "storage/static-pages/"

type StaticPageHandler struct{}

type StaticPageContentFiles struct {
	Path    string `json:"path"`
	Content string `json:"content"`
	Hash    string `json:"hash"`
}

type StaticPageContent struct {
	PagePath string                   `json:"page_path"`
	Files    []StaticPageContentFiles `json:"files"`
}

func NewStaticPageHandler(p string) *StaticPageHandler {
	storagePath = p
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

func (h *StaticPageHandler) CreateOrUpdate(w http.ResponseWriter, r *http.Request) {
	ct := r.Header.Get("Content-Type")
	if ct != "" {
		mediaType := strings.ToLower(strings.TrimSpace(strings.Split(ct, ";")[0]))
		if mediaType != "application/json" {
			fmt.Printf("Unsupported Media Type: %s\n", mediaType)
			http.Error(w, "Unsupported Media Type", http.StatusUnsupportedMediaType)
			return
		}
	}

	decodeBody := json.NewDecoder(r.Body)
	decodeBody.DisallowUnknownFields()
	var spc StaticPageContent

	if err := decodeBody.Decode(&spc); err != nil {
		http.Error(w, "Invalid JSON body", http.StatusBadRequest)
		return
	}

	pageDir := storagePath + spc.PagePath
	_, err := os.Stat(pageDir)
	if os.IsNotExist(err) {
		err = os.MkdirAll(pageDir, os.ModePerm)
		if err != nil {
			http.Error(w, "Failed to create directory for static page", http.StatusInternalServerError)
			return
		}
	}

	for _, item := range spc.Files {
		itemPath := pageDir + "/"
		if strings.HasPrefix(item.Path, "/") {
			itemPath += item.Path[1:]
		} else {
			itemPath += item.Path
		}
		itemBasePath := filepath.Dir(itemPath)
		if itemBasePath != "" && !helpers.IsFileExist(itemBasePath) {
			fmt.Printf("Creating directory %s...\n", itemBasePath)
			err = os.MkdirAll(itemBasePath, os.ModePerm)
			if err != nil {
				http.Error(w, "Failed to create directory for static page", http.StatusInternalServerError)
				return
			}
		}

		if helpers.IsFileExist(itemPath) {
			hash, err := helpers.FileChecksum(itemPath)
			if err != nil {
				http.Error(w, "Failed to calculate file checksum", http.StatusInternalServerError)
				return
			}
			if hash == item.Hash {
				continue
			}
		}

		decodedFile, err := base64.StdEncoding.DecodeString(item.Content)
		if err != nil {
			http.Error(w, "Failed to decode file content", http.StatusBadRequest)
			return
		}

		fmt.Printf("Writing file %s...\n", itemPath)
		err = os.WriteFile(itemPath, []byte(decodedFile), os.ModePerm)
		if err != nil {
			http.Error(w, "Failed to create files", http.StatusInternalServerError)
			fmt.Printf("Failed to write file %s: %v\n", itemPath, err)
			return
		}
	}
}
