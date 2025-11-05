package client

type EndpointConfig struct {
	Endpoint APIEndpoint
	BaseURL  string
}

func NewEndpointConfig(endpoint APIEndpoint) *EndpointConfig {
	return &EndpointConfig{
		Endpoint: endpoint,
		BaseURL:  BaseURLs[endpoint],
	}
}

func GetBaseURL(endpoint APIEndpoint) string {
	if url, ok := BaseURLs[endpoint]; ok {
		return url
	}
	return ""
}
