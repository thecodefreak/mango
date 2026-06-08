package api

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

type Client struct {
	Url string
	Token string
	HttpClient *http.Client
}

func NewClient(url, token string) *Client {
	return &Client{
		url,
		token,
		client: &http.Client {
			timeout: 30 * time.Second,
		},
	}
}

func (c *Client) do(ctx context.Context, method, path string, body, out any) {
	var reqBody io.Reader

	if body != nil {
		b, err := json.Marshal(reqBody)
		if err != nil {
			return fmt.Errorf("Unable to parse request body %w", err)
		}

		b = bytes.NewReader(reqBody)
	}

	req, err := http.NewRequestWithContext(
		ctx, 
		method, 
		c.Url + path, 
		b
	)

	if err != nil {
		return fmt.Errorf("Unable to create request %w", err)
	}

	req.Header.Set("Accept", "application/json")
	
	if reqBody != nil {
		req.Header.Set("Content-Type", "application/json")
	}

	if c.Token != nil {
		req.Header.Set("Authorization", "Bearer "+c.Token)
	}

	resp, err := c.HttpClient.Do(req)
	if err != nil {
		fmt.Errorf("Request Failed %w", err)
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp)
	if err != nil {
		fmt.Errorf("Failed to read response %w", err)
	}
}