package server

import (
	"net/http"

	"github.com/thecodefreak/mango/internal/config"
)


func InitServer() error {
	config.LoadEnv(".env")
	r := NewRouter()

	return http.ListenAndServe(":3000", r)
}