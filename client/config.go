package client

import "time"

type APIEndpoint string

const (
	API      APIEndpoint = "api"
	Mahoraga APIEndpoint = "mahoraga"
)

var BaseURLs = map[APIEndpoint]string{
	API:      "https://datosabiertos.dgcp.gob.do/api-dgcp/v1",
	Mahoraga: "https://mahoraga.dgcp.gob.do/api/v1",
}

type SDKConfig struct {
	APIKey       string
	Timeout      time.Duration
	MaxRetries   int
	RetryDelay   time.Duration
	Debug        bool
	CustomHeader map[string]string
}

// DefaultSDKConfig return a new SDKConfig with the default values
func DefaultSDKConfig(apiKey string) *SDKConfig {
	return &SDKConfig{
		APIKey:       apiKey,
		Timeout:      30 * time.Second,
		MaxRetries:   3,
		RetryDelay:   1 * time.Second,
		Debug:        false,
		CustomHeader: make(map[string]string),
	}
}

func (c *SDKConfig) WithTimeout(timeout time.Duration) *SDKConfig {
	c.Timeout = timeout
	return c
}

func (c *SDKConfig) WithMaxRetries(retries int) *SDKConfig {
	c.MaxRetries = retries
	return c
}

func (c *SDKConfig) WithDebug(debug bool) *SDKConfig {
	c.Debug = debug
	return c
}

func (c *SDKConfig) WithCustomHeader(key, value string) *SDKConfig {
	c.CustomHeader[key] = value
	return c
}
