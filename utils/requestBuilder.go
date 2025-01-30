package Utils

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
)

// Builder pattern
type AuthedRequest struct {
	AuthToken string
	BaseURL   string
	path      string
	body      []byte
}

type AuthedRequestBuilder interface {
	SetBaseURL(string) *AuthedRequest
	SetPath(string) *AuthedRequest
	SetBodyFromStruct(interface{}) *AuthedRequest
	Build(context.Context) (*http.Request, error)
}

// NewRequestBuilder initializes a new RequestBuilder.
func NewRequestBuilder(authToken string) *AuthedRequest {
	return &AuthedRequest{
		AuthToken: authToken,
	}
}

// SetBaseURL sets the base URL for the request.
func (rb *AuthedRequest) SetBaseURL(BaseURL string) *AuthedRequest {
	rb.BaseURL = BaseURL
	return rb
}

// SetPath sets the path for the request.
func (rb *AuthedRequest) SetPath(path string) *AuthedRequest {
	rb.path = path
	return rb
}

// SetBodyFromStruct sets the request body from a struct, marshaling it to JSON.

// transform interface to ChatCompletionRequest
func (rb *AuthedRequest) SetBodyFromStruct(data interface{}) *AuthedRequest {
	body, err := json.Marshal(data)
	if err != nil {
		panic(fmt.Sprintf("failed to marshal body: %v", err)) // Using panic for debugging; replace in future with proper error handling.
	}
	rb.body = body
	return rb
}

// Build constructs the HTTP request [Method:Post].
func (rb *AuthedRequest) Build(ctx context.Context) (*http.Request, error) {
	if rb.BaseURL == "" || rb.path == "" {
		return nil, fmt.Errorf("BaseURL or path not set")
	}

	url := fmt.Sprintf("%s%s", rb.BaseURL, rb.path)
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, url, bytes.NewReader(rb.body))
	if err != nil {
		return nil, fmt.Errorf("error creating request: %w", err)
	}
	req.Header.Set("Authorization", "Bearer "+rb.AuthToken)
	req.Header.Set("Content-Type", "application/json")
	return req, nil
}

// Build constructs the HTTP request [Method:Post].
func (rb *AuthedRequest) BuildStream(ctx context.Context) (*http.Request, error) {
	if rb.BaseURL == "" || rb.path == "" {
		return nil, fmt.Errorf("BaseURL or path not set")
	}

	url := fmt.Sprintf("%s%s", rb.BaseURL, rb.path)
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, url, bytes.NewReader(rb.body))
	if err != nil {
		return nil, fmt.Errorf("error creating request: %w", err)
	}
	req.Header.Set("Authorization", "Bearer "+rb.AuthToken)
	req.Header.Set("cache-control", "no-cache")
	req.Header.Set("Content-Type", "application/json")
	return req, nil
}

// Build constructs the HTTP request [Method:Get].
func (rb *AuthedRequest) BuildGet(ctx context.Context) (*http.Request, error) {
	if rb.BaseURL == "" || rb.path == "" {
		return nil, fmt.Errorf("BaseURL or path not set")
	}

	url := fmt.Sprintf("%s%s", rb.BaseURL, rb.path)
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, bytes.NewReader(rb.body))
	if err != nil {
		return nil, fmt.Errorf("error creating request: %w", err)
	}
	req.Header.Set("Authorization", "Bearer "+rb.AuthToken)
	req.Header.Set("Content-Type", "application/json")
	return req, nil
}
