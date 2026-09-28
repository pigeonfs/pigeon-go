package pigeon

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"
)

const (
	version     = "0.1.0"
	userAgent   = "pigeon-go/" + version
	contentType = "application/json"
)

var defaultBaseURL = envOr("PIGEON_BASE_URL", "http://localhost:4005/")

var defaultHTTPClient = &http.Client{Timeout: time.Minute}

// Client talks to the Pigeon HTTP API. Layout follows resend-go:
// client.Emails.Send, client.Domains, client.Contacts.
type Client struct {
	client    *http.Client
	ApiKey    string
	BaseURL   *url.URL
	UserAgent string
	Emails    *EmailsSvc
	Domains   *DomainsSvc
	Contacts  *ContactsSvc
	Automations *AutomationsSvc
}

func NewClient(apiKey string) *Client {
	return NewCustomClient(defaultHTTPClient, strings.TrimSpace(apiKey))
}

func NewCustomClient(httpClient *http.Client, apiKey string) *Client {
	if httpClient == nil {
		httpClient = defaultHTTPClient
	}
	if !strings.HasSuffix(defaultBaseURL, "/") {
		defaultBaseURL += "/"
	}
	baseURL, _ := url.Parse(defaultBaseURL)
	c := &Client{client: httpClient, BaseURL: baseURL, UserAgent: userAgent, ApiKey: apiKey}
	c.Emails = &EmailsSvc{client: c}
	c.Domains = &DomainsSvc{client: c}
	c.Contacts = &ContactsSvc{client: c}
	c.Automations = &AutomationsSvc{client: c}
	return c
}

func (c *Client) NewRequest(ctx context.Context, method, path string, params any) (*http.Request, error) {
	u, err := c.BaseURL.Parse(strings.TrimPrefix(path, "/"))
	if err != nil {
		return nil, err
	}
	req, err := http.NewRequestWithContext(ctx, method, u.String(), nil)
	if err != nil {
		return nil, err
	}
	if params != nil {
		buf := new(bytes.Buffer)
		if err := json.NewEncoder(buf).Encode(params); err != nil {
			return nil, err
		}
		req.Body = io.NopCloser(buf)
		req.Header.Set("Content-Type", contentType)
	}
	req.Header.Set("Accept", contentType)
	req.Header.Set("User-Agent", c.UserAgent)
	req.Header.Set("Authorization", "Bearer "+c.ApiKey)
	return req, nil
}

func (c *Client) Perform(req *http.Request, ret any) (*http.Response, error) {
	resp, err := c.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return nil, handleError(resp)
	}
	if ret != nil && resp.StatusCode != http.StatusNoContent && resp.Body != nil {
		if err := json.NewDecoder(resp.Body).Decode(ret); err != nil && !errors.Is(err, io.EOF) {
			return nil, err
		}
	}
	return resp, nil
}

type APIError struct {
	StatusCode int    `json:"statusCode"`
	Name       string `json:"name"`
	Message    string `json:"message"`
}

func (e *APIError) Error() string {
	if e.Message != "" {
		return fmt.Sprintf("[ERROR]: %s", e.Message)
	}
	return fmt.Sprintf("[ERROR]: HTTP %d", e.StatusCode)
}

func handleError(resp *http.Response) error {
	err := &APIError{StatusCode: resp.StatusCode}
	_ = json.NewDecoder(resp.Body).Decode(err)
	if err.Message == "" {
		err.Message = resp.Status
	}
	return err
}

func envOr(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		if !strings.HasSuffix(v, "/") {
			v += "/"
		}
		return v
	}
	return fallback
}
