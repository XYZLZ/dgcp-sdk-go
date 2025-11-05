package client

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"time"
)

const sdkVersion = "1.0.0"

type BaseClient struct {
	config     *SDKConfig
	endpoint   APIEndpoint
	baseURL    string
	httpClient *http.Client
}

func NewBaseClient(config *SDKConfig, endpoint APIEndpoint) *BaseClient {
	return &BaseClient{
		config:   config,
		endpoint: endpoint,
		baseURL:  GetBaseURL(endpoint),
		httpClient: &http.Client{
			Timeout: config.Timeout,
		},
	}
}

func (c *BaseClient) GetEndpoint() APIEndpoint {
	return c.endpoint
}

func (c *BaseClient) GetBaseURL() string {
	return c.baseURL
}

func (c *BaseClient) Request(ctx context.Context, method, path string, body interface{}, result interface{}) error {
	url := c.baseURL + path
	var reqBody io.Reader

	if body != nil {
		switch v := body.(type) {
		case bytes.Buffer:
			reqBody = bytes.NewBuffer(v.Bytes())
		case io.Reader:
			reqBody = v
		case string:
			reqBody = bytes.NewBufferString(v)
		case []byte:
			reqBody = bytes.NewBuffer(v)
		default:
			jsonData, err := json.Marshal(body)
			if err != nil {
				return fmt.Errorf("error marshaling request body: %w", err)
			}
			reqBody = bytes.NewBuffer(jsonData)
		}
	}

	req, err := http.NewRequestWithContext(ctx, method, url, reqBody)
	if err != nil {
		return fmt.Errorf("error creating request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("User-Agent", fmt.Sprintf("dgcp-sdk-go/%s", sdkVersion))

	if c.config.APIKey != "" {
		req.Header.Set("X-API-Key", c.config.APIKey)
	}

	for key, value := range c.config.CustomHeader {
		req.Header.Set(key, value)
	}

	if c.config.Debug {
		log.Printf("[SDK Request] %s %s (Endpoint: %s)", method, url, c.endpoint)
	}

	resp, err := c.doWithRetry(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if c.config.Debug {
		log.Printf("[SDK Response] %d %s", resp.StatusCode, url)
	}

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("error reading response body: %w", err)
	}

	if resp.StatusCode >= 400 {
		return c.handleErrorResponse(resp.StatusCode, respBody)
	}
	if result != nil && len(respBody) > 0 {
		if err := json.Unmarshal(respBody, result); err != nil {
			return fmt.Errorf("error unmarshaling response: %w", err)
		}
	}

	return nil
}

func (c *BaseClient) doWithRetry(req *http.Request) (*http.Response, error) {
	var resp *http.Response
	var err error

	for attempt := 0; attempt <= c.config.MaxRetries; attempt++ {
		if attempt > 0 {
			time.Sleep(c.config.RetryDelay * time.Duration(attempt))
			if c.config.Debug {
				log.Printf("[SDK Retry] Attempt %d/%d", attempt, c.config.MaxRetries)
			}
		}

		resp, err = c.httpClient.Do(req)
		if err != nil {
			if attempt == c.config.MaxRetries {
				return nil, fmt.Errorf("request failed after %d attempts: %w", c.config.MaxRetries+1, err)
			}
			continue
		}

		if !c.shouldRetry(resp.StatusCode) {
			return resp, nil
		}

		resp.Body.Close()
	}

	return resp, err
}

func (c *BaseClient) shouldRetry(statusCode int) bool {
	retryableStatuses := []int{408, 429, 500, 502, 503, 504}
	for _, status := range retryableStatuses {
		if statusCode == status {
			return true
		}
	}
	return false
}

func (c *BaseClient) handleErrorResponse(statusCode int, body []byte) error {
	var errorResp struct {
		Message string                 `json:"message"`
		Details map[string]interface{} `json:"details"`
	}

	if err := json.Unmarshal(body, &errorResp); err != nil {
		return NewSDKError(INTERNAL_SERVER_ERROR, "Failed to parse error response", statusCode, map[string]interface{}{
			"body": string(body),
		})
	}

	switch statusCode {
	case 401:
		return AuthenticationError(errorResp.Message, errors.New("invalid Api key"))
	case 404:
		return NotFoundError(errorResp.Message)
	case 429:
		return RateLimitError(errorResp.Message, errors.New("rate limit exeeded"))
	case 400:
		fields := make(map[string]string)
		for k, v := range errorResp.Details {
			fields[k] = fmt.Sprint(v)
		}
		return ValidationError(errorResp.Message, fields)
	default:
		return NewSDKError(INTERNAL_SERVER_ERROR, errorResp.Message, statusCode, errorResp.Details)
	}
}

func (c *BaseClient) Get(ctx context.Context, path string, result interface{}) error {
	return c.Request(ctx, http.MethodGet, path, nil, result)
}

func (c *BaseClient) Post(ctx context.Context, path string, body interface{}, result interface{}) error {
	return c.Request(ctx, http.MethodPost, path, body, result)
}

func (c *BaseClient) Patch(ctx context.Context, path string, body interface{}, result interface{}) error {
	return c.Request(ctx, http.MethodPatch, path, body, result)
}

func (c *BaseClient) Put(ctx context.Context, path string, body interface{}, result interface{}) error {
	return c.Request(ctx, http.MethodPut, path, body, result)
}

func (c *BaseClient) Delete(ctx context.Context, path string) error {
	return c.Request(ctx, http.MethodDelete, path, nil, nil)
}
