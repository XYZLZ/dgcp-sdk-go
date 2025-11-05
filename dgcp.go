package dgcp

import (
	"time"

	"github.com/XYZLZ/dgcp-sdk-go/client"
	mahoragaResourse "github.com/XYZLZ/dgcp-sdk-go/core/mahoraga"
)

type Mahoraga struct {
	Apps  *mahoragaResourse.AppsResource
	Files *mahoragaResourse.FilesResource
	Auth  *mahoragaResourse.LoginResource
}

type DGCP struct {
	config *client.SDKConfig

	Mahoraga Mahoraga
}

// New returns a new DGCP instance with the given API key and options.
// The options are applied to the SDKConfig instance that is used to create
// the DGCP resources.
// The returned DGCP instance contains the DGCP resources
// that are used to interact with the Mahoraga API.
// The options are applied in the order they are given, so later options can
// override earlier ones.
// If no options are given, the default SDKConfig is used.
func New(apiKey string, opts ...Option) *DGCP {
	config := client.DefaultSDKConfig(apiKey)

	// apply options
	for _, opt := range opts {
		opt(config)
	}

	return &DGCP{
		config: config,
		Mahoraga: Mahoraga{
			Apps:  mahoragaResourse.NewAppsResource(config),
			Files: mahoragaResourse.NewFilesResource(config),
			Auth:  mahoragaResourse.NewLoginResource(config),
		},
	}
}

type Option func(*client.SDKConfig)

// WithTimeout sets the timeout
func WithTimeout(timeout time.Duration) Option {
	return func(c *client.SDKConfig) {
		c.Timeout = timeout
	}
}

// WithMaxRetries sets the maximum number of retries
func WithMaxRetries(retries int) Option {
	return func(c *client.SDKConfig) {
		c.MaxRetries = retries
	}
}

// WithDebug allows you to enable debug mode
func WithDebug(debug bool) Option {
	return func(c *client.SDKConfig) {
		c.Debug = debug
	}
}

// WithCustomHeader adds a custom header
func WithCustomHeader(key, value string) Option {
	return func(c *client.SDKConfig) {
		c.CustomHeader[key] = value
	}
}

// GetEndpointInfo returns a map of API endpoints
func (d *DGCP) GetEndpointInfo() map[string]map[string]client.APIEndpoint {
	return map[string]map[string]client.APIEndpoint{
		"mahoraga": {
			"apps":  d.Mahoraga.Apps.GetEndpoint(),
			"files": d.Mahoraga.Files.GetEndpoint(),
			"auth":  d.Mahoraga.Auth.GetEndpoint(),
		},
	}
}
