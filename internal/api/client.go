package api

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"
)

type Client struct {
	Url        string
	Token      string
	HttpClient *http.Client
}

func NewClient(url, token string) *Client {
	return &Client{
		Url:   url,
		Token: token,
		HttpClient: &http.Client{
			Timeout: 30 * time.Second,
		},
	}
}

func (c *Client) do(ctx context.Context, method, path string, body, out any) error {
	var reqBody io.Reader

	if body != nil {
		b, err := json.Marshal(body)
		if err != nil {
			return fmt.Errorf("Unable to parse request body %w", err)
		}

		reqBody = bytes.NewReader(b)
	}

	if !strings.HasSuffix(c.Url, "/api") {
		c.Url = strings.TrimSuffix(c.Url, "/") + "/api"
	}
	req, err := http.NewRequestWithContext(
		ctx,
		method,
		c.Url+path,
		reqBody,
	)

	if err != nil {
		return fmt.Errorf("Unable to create request %w", err)
	}

	req.Header.Set("Accept", "application/json")

	if reqBody != nil {
		req.Header.Set("Content-Type", "application/json")
	}

	if c.Token != "" {
		req.Header.Set("Authorization", "Bearer "+c.Token)
	}

	resp, err := c.HttpClient.Do(req)
	if err != nil {
		return fmt.Errorf("Request Failed %w", err)
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("Failed to read response %w", err)
	}

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return fmt.Errorf("Request failed with status %d: %s", resp.StatusCode, string(respBody))
	}

	if out == nil {
		return nil
	}

	if len(respBody) == 0 {
		return nil
	}

	if err := json.Unmarshal(respBody, out); err != nil {
		return fmt.Errorf("Failed to parse response %w", err)
	}

	return nil
}
