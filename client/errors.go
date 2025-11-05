package client

import (
	"fmt"
	"time"
)

type SDKError struct {
	Message    string
	Code       string
	StatusCode int
	Details    map[string]interface{}
}

func (e *SDKError) Error() string {
	return fmt.Sprintf("SDK Error [%d]: %s", e.StatusCode, e.Message)
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
		StatusCode: 401,
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
		StatusCode: 429,
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
		StatusCode: 400,
		Details: map[string]interface{}{
			"fields": fields,
		},
	}
}

func ApiError(message string, cause error) *SDKError {
	return &SDKError{
		Code:       API_ERROR,
		Message:    "API request failed: " + message,
		StatusCode: 500,
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
		StatusCode: 404,
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
		StatusCode: 401,
		Details: map[string]interface{}{
			"message": message,
		},
	}
}
