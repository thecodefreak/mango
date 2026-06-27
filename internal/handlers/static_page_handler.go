package handlers

import (
	"encoding/base64"
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

type StaticPageCheckFiles struct {
	Path string `json:"path"`
	Hash string `json:"hash"`
}

type StaticPageCheck struct {
	PagePath string                 `json:"page_path"`
	Files    []StaticPageCheckFiles `json:"files"`
}

type StaticPageRequired struct {
	PagePath string   `json:"page_path"`
	Files    []string `json:"files"`
}

func NewStaticPageHandler(p string) *StaticPageHandler {
	storagePath = filepath.Clean(p) + "/"
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
	err := helpers.ValidateRequestMediaType(r, "application/json")
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnsupportedMediaType)
		return
	}

	var spc StaticPageContent
	err = helpers.GetJsonBody(r, &spc)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	pageDir := storagePath + spc.PagePath
	_, err = os.Stat(pageDir)
	if os.IsNotExist(err) {
		err = os.MkdirAll(pageDir, os.ModePerm)
		if err != nil {
			http.Error(w, "Failed to create directory for static page", http.StatusInternalServerError)
			return
		}
	}

	basePath := pageDir + "/"
	for _, item := range spc.Files {
		itemPath := basePath + helpers.PathWoSlash(item.Path, true)
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

func (h *StaticPageHandler) Check(w http.ResponseWriter, r *http.Request) {
	err := helpers.ValidateRequestMediaType(r, "application/json")
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnsupportedMediaType)
		return
	}

	var spc StaticPageCheck
	var spr StaticPageRequired
	err = helpers.GetJsonBody(r, &spc)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	pagePath := storagePath + spc.PagePath
	spr.PagePath = spc.PagePath

	if !helpers.IsFileExist(pagePath) {
		for _, item := range spc.Files {
			spr.Files = append(spr.Files, item.Path)
		}
		helpers.JsonResponse(w, http.StatusOK, spr)
		return
	}

	for _, item := range spc.Files {
		itemPath := pagePath + "/" + helpers.PathWoSlash(item.Path, true)
		itemBasePath := filepath.Dir(itemPath)

		if !helpers.IsFileExist(itemBasePath) {
			spr.Files = append(spr.Files, item.Path)
			continue
		}

		if !helpers.IsFileExist(itemPath) {
			spr.Files = append(spr.Files, item.Path)
			continue
		}

		hash, err := helpers.FileChecksum(itemPath)
		if err != nil {
			http.Error(w, "Failed to calculate file checksum", http.StatusInternalServerError)
			return
		}

		if hash != item.Hash {
			spr.Files = append(spr.Files, item.Path)
		}
	}

	helpers.JsonResponse(w, http.StatusOK, spr)
}
