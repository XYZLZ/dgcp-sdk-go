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

type parsedErrorInfo struct {
	Message string
	Details map[string]interface{}
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

func (c *BaseClient) Request(ctx context.Context, method, path string, body interface{}, result interface{}, customHeaders *map[string]string) error {
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

	if customHeaders != nil {
		for key, value := range *customHeaders {
			req.Header.Set(key, value)
		}
	}

	if c.config.Debug {
		log.Printf("[SDK Request] %s %s (Endpoint: %s)", method, url, c.endpoint)
		for key, value := range req.Header {
			log.Printf("[SDK Request Header] %s: %s", key, value)
		}
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
		err := c.handleErrorResponse(resp.StatusCode, respBody)
		if c.config.Debug {
			log.Printf("[SDK Error] %s", err.Error())
		}
		return err
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

	if len(body) == 0 {
		return c.createErrorByStatusCode(statusCode, "No error details provided", nil)
	}

	parsedError := c.parseErrorBody(body)

	return c.createErrorByStatusCode(statusCode, parsedError.Message, parsedError.Details)
}

func (c *BaseClient) parseErrorBody(body []byte) parsedErrorInfo {
	result := parsedErrorInfo{
		Message: string(body),
		Details: make(map[string]interface{}),
	}

	// 1. estándar format: {"message": "...", "details": {...}}
	var standardError struct {
		Message string                 `json:"message"`
		Details map[string]interface{} `json:"details"`
	}
	if err := json.Unmarshal(body, &standardError); err == nil && standardError.Message != "" {
		result.Message = standardError.Message
		if standardError.Details != nil {
			result.Details = standardError.Details
		}
		return result
	}

	// 2.  alternative format: {"error": "...", "data": {...}}
	var altError1 struct {
		Error string                 `json:"error"`
		Data  map[string]interface{} `json:"data"`
	}
	if err := json.Unmarshal(body, &altError1); err == nil && altError1.Error != "" {
		result.Message = altError1.Error
		if altError1.Data != nil {
			result.Details = altError1.Data
		}
		return result
	}

	// 3. alternative format: {"msg": "...", "errors": {...}}
	var altError2 struct {
		Msg    string                 `json:"msg"`
		Errors map[string]interface{} `json:"errors"`
	}
	if err := json.Unmarshal(body, &altError2); err == nil && altError2.Msg != "" {
		result.Message = altError2.Msg
		if altError2.Errors != nil {
			result.Details = altError2.Errors
		}
		return result
	}

	// 4. format with error_description: {"error_description": "..."}
	var altError3 struct {
		ErrorDescription string `json:"error_description"`
	}
	if err := json.Unmarshal(body, &altError3); err == nil && altError3.ErrorDescription != "" {
		result.Message = altError3.ErrorDescription
		return result
	}

	// 5. String: "error message"
	var simpleString string
	if err := json.Unmarshal(body, &simpleString); err == nil && simpleString != "" {
		result.Message = simpleString
		return result
	}

	// 6. error Array: ["error1", "error2"]
	var errorArray []string
	if err := json.Unmarshal(body, &errorArray); err == nil && len(errorArray) > 0 {
		result.Message = errorArray[0] // Usar el primer error como mensaje principal
		if len(errorArray) > 1 {
			result.Details["all_errors"] = errorArray
		}
		return result
	}

	// 7. Object Array: [{"field": "email", "message": "invalid"}]
	var errorObjectArray []map[string]interface{}
	if err := json.Unmarshal(body, &errorObjectArray); err == nil && len(errorObjectArray) > 0 {
		if msg, ok := errorObjectArray[0]["message"].(string); ok {
			result.Message = msg
		} else if msg, ok := errorObjectArray[0]["error"].(string); ok {
			result.Message = msg
		}
		result.Details["validation_errors"] = errorObjectArray
		return result
	}

	// 8. Generic Object - extraer cualquier campo que parezca un mensaje
	var genericObject map[string]interface{}
	if err := json.Unmarshal(body, &genericObject); err == nil {
		// search error fields
		possibleMessageFields := []string{"message", "error", "msg", "description", "detail", "title"}
		for _, field := range possibleMessageFields {
			if msg, ok := genericObject[field].(string); ok && msg != "" {
				result.Message = msg
				delete(genericObject, field)
				break
			}
		}

		if len(genericObject) > 0 {
			result.Details = genericObject
		}
		return result
	}

	// 9 Response object used by some endpoints
	var altError4 struct {
		Code     int  `json:"code"`
		HasError bool `json:"hasError"`
		Payload  struct {
			Message string      `json:"message"`
			Content interface{} `json:"content"`
			Errors  []string    `json:"errors"`
		} `json:"payload"`
	}

	if err := json.Unmarshal(body, &altError4); err == nil && altError4.Payload.Message != "" {
		result.Message = altError4.Payload.Message
		if len(altError4.Payload.Errors) > 0 {
			result.Details["errors"] = altError4.Payload.Errors
		}
		return result
	}

	result.Message = string(body)
	result.Details["raw_body"] = string(body)

	if c.config.Debug {
		fmt.Printf("[SDK] Could not parse error body as JSON, using raw string: %s\n", string(body))
	}

	return result
}

func (c *BaseClient) createErrorByStatusCode(statusCode int, message string, details map[string]interface{}) error {
	switch statusCode {
	case http.StatusUnauthorized:
		return AuthenticationError(message, errors.New("authentication failed"))

	case http.StatusForbidden:
		if details == nil {
			details = make(map[string]interface{})
		}
		details["status"] = "forbidden"
		return NewSDKError(AUTHENTICATION_ERROR, message, statusCode, details)

	case http.StatusNotFound:
		return NotFoundError(message)

	case http.StatusTooManyRequests:
		return RateLimitError(message, errors.New("rate limit exceeded"))

	case http.StatusBadRequest:
		fields := extractValidationFields(details)
		if len(fields) > 0 {
			return ValidationError(message, fields)
		}
		return NewSDKError(VALIDATION_ERROR, message, statusCode, details)

	case http.StatusRequestTimeout:
		return TimeoutError(message, errors.New("request timeout"))

	case http.StatusUnprocessableEntity:
		fields := extractValidationFields(details)
		if len(fields) > 0 {
			return ValidationError(message, fields)
		}
		return NewSDKError(VALIDATION_ERROR, message, statusCode, details)

	case http.StatusInternalServerError,
		http.StatusBadGateway,
		http.StatusServiceUnavailable,
		http.StatusGatewayTimeout:
		return NewSDKError(INTERNAL_SERVER_ERROR, message, statusCode, details)

	default:
		if statusCode >= 500 {
			return NewSDKError(INTERNAL_SERVER_ERROR, message, statusCode, details)
		} else if statusCode >= 400 {
			return NewSDKError(VALIDATION_ERROR, message, statusCode, details)
		}
		return NewSDKError(INTERNAL_SERVER_ERROR, message, statusCode, details)
	}
}

func extractValidationFields(details map[string]interface{}) map[string]string {
	fields := make(map[string]string)

	if details == nil {
		return fields
	}

	// pattern 1: {"field_name": "error message"}
	for key, value := range details {
		if strValue, ok := value.(string); ok {
			fields[key] = strValue
		}
	}

	// pattern 2: {"errors": {"field_name": "error message"}}
	if errors, ok := details["errors"].(map[string]interface{}); ok {
		for key, value := range errors {
			if strValue, ok := value.(string); ok {
				fields[key] = strValue
			}
		}
	}

	// pattern 3: {"validation_errors": [{"field": "email", "message": "invalid"}]}
	if validationErrors, ok := details["validation_errors"].([]interface{}); ok {
		for _, item := range validationErrors {
			if errMap, ok := item.(map[string]interface{}); ok {
				field := ""
				message := ""

				if f, ok := errMap["field"].(string); ok {
					field = f
				}
				if m, ok := errMap["message"].(string); ok {
					message = m
				}

				if field != "" && message != "" {
					fields[field] = message
				}
			}
		}
	}

	// pattern 4: {"all_errors": ["error1", "error2"]}
	if allErrors, ok := details["all_errors"].([]interface{}); ok {
		for i, err := range allErrors {
			if strErr, ok := err.(string); ok {
				fields[fmt.Sprintf("error_%d", i)] = strErr
			}
		}
	}

	// pattern 5: {"payload": {"errors": ["error1", "error2"]}}
	if payload, ok := details["payload"].(map[string]interface{}); ok {
		if errors, ok := payload["errors"].([]interface{}); ok {
			for i, err := range errors {
				if strErr, ok := err.(string); ok {
					fields[fmt.Sprintf("error_%d", i)] = strErr
				}
			}
		}
	}

	return fields
}

func (c *BaseClient) Get(ctx context.Context, path string, result interface{}, customHeaders *map[string]string) error {
	return c.Request(ctx, http.MethodGet, path, nil, result, customHeaders)
}

func (c *BaseClient) Post(ctx context.Context, path string, body interface{}, result interface{}, customHeaders *map[string]string) error {
	return c.Request(ctx, http.MethodPost, path, body, result, customHeaders)
}

func (c *BaseClient) Patch(ctx context.Context, path string, body interface{}, result interface{}, customHeaders *map[string]string) error {
	return c.Request(ctx, http.MethodPatch, path, body, result, customHeaders)
}

func (c *BaseClient) Put(ctx context.Context, path string, body interface{}, result interface{}, customHeaders *map[string]string) error {
	return c.Request(ctx, http.MethodPut, path, body, result, customHeaders)
}

func (c *BaseClient) Delete(ctx context.Context, path string, customHeaders *map[string]string) error {
	return c.Request(ctx, http.MethodDelete, path, nil, nil, customHeaders)
}
