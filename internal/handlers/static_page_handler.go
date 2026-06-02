package handlers

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strings"

	"github.com/go-chi/chi/v5"
	"github.com/thecodefreak/mango/internal/helpers"
)

var storagePath = "storage/static-pages/"

type StaticPageHandler struct{}

type StaticPageContent struct {
	PagePath string `json:"page_path"`
	Files    []struct {
		Path  string `json:"path"`
		Files []struct {
			Name    string `json:"name"`
			Content string `json:"content"`
			Hash    string `json:"hash"`
		} `json:"files"`
	} `json:"files"`
}

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
		itemPath := pageDir + "/" + item.Path
		_, err := os.Stat(itemPath)
		if os.IsExist(err) && item.Path != "" {
			os.Mkdir(item.Path, os.ModePerm)
		}

		for _, file := range item.Files {
			filePath := itemPath + file.Name
			if helpers.IsFileExist(filePath) {
				hash, err := helpers.FileChecksum(filePath)
				if err != nil {
					http.Error(w, "Failed to calculate file checksum", http.StatusInternalServerError)
					return
				}
				if hash == file.Hash {
					continue
				}
			}

			decodedFile, err := base64.StdEncoding.DecodeString(file.Content)
			if err != nil {
				http.Error(w, "Failed to decode file content", http.StatusBadRequest)
				return
			}
			err = os.WriteFile(filePath, []byte(decodedFile), os.ModePerm)
			if err != nil {
				http.Error(w, "Failed to create files", http.StatusInternalServerError)
				fmt.Printf("Failed to write file %s: %v\n", filePath, err)
				return
			}
		}
	}
}
