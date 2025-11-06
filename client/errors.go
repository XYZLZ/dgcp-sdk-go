package client

import (
	"fmt"
	"net/http"
	"sort"
	"time"
)

type SDKError struct {
	Message    string
	Code       string
	StatusCode int
	Details    map[string]interface{}
	Err        error `json:"-"`
}

func (e *SDKError) Error() string {
	var detailsStr string

	if e.Details != nil {
		var keys []string
		for k := range e.Details {
			keys = append(keys, k)
		}
		sort.Strings(keys)
		for _, k := range keys {
			detailsStr += fmt.Sprintf("  %s: %v\n", k, e.Details[k])
		}
	}

	if e.Err != nil {
		return fmt.Sprintf("SDK Error [%d]: %s -\nDetails:\n%s\nCause:\n%s", e.StatusCode, e.Message, detailsStr, e.Err.Error())
	}
	return fmt.Sprintf("SDK Error [%d]: %s -\nDetails:\n%s", e.StatusCode, e.Message, detailsStr)
}

func (e *SDKError) Unwrap() error {
	return e.Err
}

// error codes
var (
	AUTHENTICATION_ERROR       = "AUTHENTICATION_ERROR"
	RATE_LIMIT_ERROR           = "RATE_LIMIT_ERROR"
	VALIDATION_ERROR           = "VALIDATION_ERROR"
	API_ERROR                  = "API_ERROR"
	NOT_FOUND_ERROR            = "NOT_FOUND_ERROR"
	CONFLICT_ERROR             = "CONFLICT_ERROR"
	UNPROCESSABLE_ENTITY_ERROR = "UNPROCESSABLE_ENTITY_ERROR"
	INTERNAL_SERVER_ERROR      = "INTERNAL_SERVER_ERROR"
	NETWORK_ERROR              = "NETWORK_ERROR"
	PARSE_ERROR                = "PARSE_ERROR"
	TIMEOUT_ERROR              = "TIMEOUT_ERROR"
	retryAfter                 = 1 * time.Minute
)

func NewSDKError(code string, message string, statusCode int, details map[string]interface{}) *SDKError {
	return &SDKError{
		Message:    message,
		StatusCode: statusCode,
		Details:    details,
	}
}
func AuthenticationError(message string, cause error) *SDKError {
	return &SDKError{
		Code:       AUTHENTICATION_ERROR,
		Message:    "Invalid API key: " + message,
		StatusCode: http.StatusUnauthorized,
		Details: map[string]interface{}{
			"message": message,
			"cause":   cause.Error(),
		},
	}
}

func RateLimitError(message string, cause error) *SDKError {
	return &SDKError{
		Code:       RATE_LIMIT_ERROR,
		Message:    "Rate limit exceeded: " + message,
		StatusCode: http.StatusTooManyRequests,
		Details: map[string]interface{}{
			"retryAfter": retryAfter.String(),
			"cause":      cause.Error(),
			"message":    message,
		},
	}
}

func ValidationError(message string, fields map[string]string) *SDKError {
	return &SDKError{
		Code:       VALIDATION_ERROR,
		Message:    message,
		StatusCode: http.StatusBadRequest,
		Details: map[string]interface{}{
			"fields": fields,
		},
	}
}

func ApiError(message string, cause error) *SDKError {
	return &SDKError{
		Code:       API_ERROR,
		Message:    "API request failed: " + message,
		StatusCode: http.StatusInternalServerError,
		Details: map[string]interface{}{
			"message": message,
			"cause":   cause.Error(),
		},
	}
}

func NotFoundError(message string) *SDKError {
	return &SDKError{
		Code:       NOT_FOUND_ERROR,
		Message:    message,
		StatusCode: http.StatusNotFound,
	}
}

func NetworkError(message string, cause error) *SDKError {
	return &SDKError{
		Code:       NETWORK_ERROR,
		Message:    "Network error: " + message,
		StatusCode: 0,
		Details: map[string]interface{}{
			"message": message,
			"cause":   cause.Error(),
		},
	}
}

func ApiKeyRequiredError(message string) *SDKError {
	return &SDKError{
		Code:       AUTHENTICATION_ERROR,
		Message:    message,
		StatusCode: http.StatusUnauthorized,
		Details: map[string]interface{}{
			"message": message,
		},
	}
}

func InternalServerError(message string, cause error) *SDKError {
	return &SDKError{
		Code:       INTERNAL_SERVER_ERROR,
		Message:    "Internal server error: " + message,
		StatusCode: http.StatusInternalServerError,
		Details: map[string]interface{}{
			"message": message,
			"cause":   cause.Error(),
		},
	}
}

func TimeoutError(message string, cause error) *SDKError {
	if message == "" {
		message = "Request timeout"
	}
	return &SDKError{
		Code:       TIMEOUT_ERROR,
		Message:    message,
		StatusCode: http.StatusRequestTimeout,
		Details: map[string]interface{}{
			"message": message,
			"cause":   cause.Error(),
		},
	}
}
