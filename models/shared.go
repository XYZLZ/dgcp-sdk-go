package models

import "net/http"

type PaginationRequest struct {
	Page  *int
	Limit *int
}

type ResponseMetadata struct {
	Headers    http.Header
	StatusCode int
	RequestID  string
}

// CallOption es una función que modifica las opciones de llamada
type CallOption func(*CallOptions)

type CallOptions struct {
	IncludeMetadata bool
	Metadata        *ResponseMetadata
}

func WithMetadata(metadata *ResponseMetadata) CallOption {
	return func(opts *CallOptions) {
		opts.IncludeMetadata = true
		opts.Metadata = metadata
	}
}

type ResponseWithMetadata[T any] struct {
	Data     T
	Metadata *ResponseMetadata `json:"-"`
}
