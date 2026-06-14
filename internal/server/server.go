package server

import (
	"fmt"
	"net/http"

	"github.com/thecodefreak/mango/internal/config"
)

func InitServer(c *config.ServerConfig) error {
	r := NewRouter(c)

	fmt.Printf("Started mango server at %s\n", c.Addr)
	err := http.ListenAndServe(c.Addr, r)
	if err != nil {
		return fmt.Errorf("Unable to start server, %w", err)
	}

	return nil
}
