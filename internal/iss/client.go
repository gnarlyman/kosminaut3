package iss

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
)

const endpoint = "https://api.wheretheiss.at/v1/satellites/25544"

type Client struct {
	http *http.Client
}

func NewClient(httpClient *http.Client) *Client {
	return &Client{http: httpClient}
}

func (c *Client) Fetch(ctx context.Context) (Position, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, endpoint, nil)
	if err != nil {
		return Position{}, err
	}
	resp, err := c.http.Do(req)
	if err != nil {
		return Position{}, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return Position{}, fmt.Errorf("wheretheiss.at: status %d", resp.StatusCode)
	}

	var pos Position
	if err := json.NewDecoder(resp.Body).Decode(&pos); err != nil {
		return Position{}, fmt.Errorf("decode: %w", err)
	}
	return pos, nil
}
