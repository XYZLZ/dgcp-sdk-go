package mahoraga

import (
	"context"

	sdkClient "github.com/XYZLZ/dgcp-sdk-go/client"
	"github.com/XYZLZ/dgcp-sdk-go/models"
	mahoModels "github.com/XYZLZ/dgcp-sdk-go/models/mahoraga"
)

type LoginResource struct {
	*sdkClient.BaseClient
}

func NewLoginResource(config *sdkClient.SDKConfig) *LoginResource {
	return &LoginResource{
		BaseClient: sdkClient.NewBaseClient(config, sdkClient.Mahoraga),
	}
}

// Login authenticates a user and returns a token to be used in the API.
//
// The endpoint will return a slice of LoginServicePayload structs containing the access token and refresh token.
//
// The endpoint will return an error if the request fails.
// The error will contain the status code of the response and the body of the response.
func (r *LoginResource) Login(ctx context.Context, credentials mahoModels.Login, opts ...models.CallOption) (*mahoModels.MahoragaResponse[mahoModels.LoginServicePayload], error) {
	var result mahoModels.MahoragaResponse[mahoModels.LoginServicePayload]
	callOpts := &models.CallOptions{}

	for _, opt := range opts {
		opt(callOpts)
	}

	path := "/auth/login"

	err := r.Post(ctx, path, credentials, &result, nil, callOpts)
	return &result, err
}
