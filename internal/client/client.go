package client

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"strings"
	"time"
)

const apiBasePath = "/api/v3"

type Client struct {
	endpoint   string
	apiKey     string
	httpClient *http.Client
}

func NewClient(endpoint string, apiKey string) *Client {
	return &Client{
		endpoint: strings.TrimRight(endpoint, "/") + apiBasePath,
		apiKey:   apiKey,
		httpClient: &http.Client{
			Timeout: 30 * time.Second,
		},
	}
}

func (c *Client) post(ctx context.Context, path string, payload any, target any) error {
	var body io.Reader

	if payload != nil {
		data, err := json.Marshal(payload)
		if err != nil {
			return err
		}

		body = bytes.NewReader(data)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, c.endpoint+path, body)
	if err != nil {
		return err
	}

	req.Header.Set("Authorization", "Bearer "+c.apiKey)
	req.Header.Set("Content-Type", "application/json")

	res, err := c.httpClient.Do(req)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	resBody, err := io.ReadAll(res.Body)
	if err != nil {
		return err
	}

	if res.StatusCode < 200 || res.StatusCode >= 300 {
		return fmt.Errorf("nethost API returned status %d: %s", res.StatusCode, string(resBody))
	}

	if target == nil || len(resBody) == 0 {
		return nil
	}

	return json.Unmarshal(resBody, target)
}

func stringifyID(value any) string {
	switch v := value.(type) {
	case string:
		return v
	case float64:
		return strconv.FormatInt(int64(v), 10)
	case int:
		return strconv.Itoa(v)
	case int64:
		return strconv.FormatInt(v, 10)
	default:
		return fmt.Sprintf("%v", v)
	}
}
