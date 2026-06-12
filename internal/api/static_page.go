package api

import (
	"context"
	"net/http"

	"github.com/thecodefreak/mango/internal/handlers"
)

func (c *Client) StaticPageCreate(r handlers.StaticPageContent) error {
	return c.do(
		context.Background(),
		http.MethodPost,
		"/static-page",
		r,
		nil,
	)
}