package mahoraga

import (
	"context"
	"fmt"

	sdkClient "github.com/XYZLZ/dgcp-sdk-go/client"
	"github.com/XYZLZ/dgcp-sdk-go/models"
	mahoModels "github.com/XYZLZ/dgcp-sdk-go/models/mahoraga"
)

type AppsResource struct {
	*sdkClient.BaseClient
}

func NewAppsResource(config *sdkClient.SDKConfig) *AppsResource {
	return &AppsResource{
		BaseClient: sdkClient.NewBaseClient(config, sdkClient.Mahoraga),
	}
}

// List returns a list of apps for a given user.
//
// The endpoint will return a slice of App structs containing the app data.
//
// The endpoint will return an error if the request fails.
// The error will contain the status code of the response and the body of the response.
func (r *AppsResource) List(ctx context.Context, userId string, opts ...models.CallOption) (*mahoModels.MahoragaResponse[[]mahoModels.App], error) {
	var result mahoModels.MahoragaResponse[[]mahoModels.App]
	callOpts := &models.CallOptions{}

	for _, opt := range opts {
		opt(callOpts)
	}

	path := "/apps/get-user?userId=" + userId

	err := r.BaseClient.Get(ctx, path, &result, nil, callOpts)
	return &result, err
}

// Create creates a new app.
//
// The endpoint will return a slice of App structs containing the app data.
//
// The endpoint will return an error if the request fails.
// The error will contain the status code of the response and the body of the response.
func (r *AppsResource) Create(ctx context.Context, app *mahoModels.App, opts ...models.CallOption) (*mahoModels.MahoragaResponse[mahoModels.App], error) {
	var res mahoModels.MahoragaResponse[mahoModels.App]
	callOpts := &models.CallOptions{}

	for _, opt := range opts {
		opt(callOpts)
	}

	app.CreatedAt = nil
	app.Id = ""
	app.UpdatedAt = nil
	err := r.BaseClient.Post(ctx, "/apps/create", app, &res, nil, callOpts)
	return &res, err
}

// Update updates an app.
//
// The endpoint will return a slice of App structs containing the app data.
//
// The endpoint will return an error if the request fails.
// The error will contain the status code of the response and the body of the response.
func (r *AppsResource) Update(ctx context.Context, app mahoModels.App, opts ...models.CallOption) (*mahoModels.MahoragaResponse[mahoModels.App], error) {
	var user mahoModels.MahoragaResponse[mahoModels.App]
	callOpts := &models.CallOptions{}

	for _, opt := range opts {
		opt(callOpts)
	}

	err := r.Put(ctx, fmt.Sprintf("/apps/update/%s", app.Id), app, &user, nil, callOpts)
	return &user, err
}

// GetSettings returns the settings of an app.
//
// The endpoint will return a slice of Settings structs containing the app settings.
//
// The endpoint will return an error if the request fails.
// The error will contain the status code of the response and the body of the response.
func (r *AppsResource) GetSettings(ctx context.Context, appId string, opts ...models.CallOption) (*mahoModels.MahoragaResponse[mahoModels.AppSettings], error) {
	var settings mahoModels.MahoragaResponse[mahoModels.AppSettings]
	callOpts := &models.CallOptions{}
	for _, opt := range opts {
		opt(callOpts)
	}
	err := r.Get(ctx, fmt.Sprintf("/apps/get-settings?appId=%s", appId), &settings, nil, callOpts)
	return &settings, err
}

// CreateSettings creates a new settings for an app.
//
// The endpoint will return a slice of Settings structs containing the app settings.
//
// The endpoint will return an error if the request fails.
// The error will contain the status code of the response and the body of the response.
func (r *AppsResource) CreateSettings(ctx context.Context, settings mahoModels.AppSettings, opts ...models.CallOption) (*mahoModels.MahoragaResponse[mahoModels.AppSettings], error) {
	var res mahoModels.MahoragaResponse[mahoModels.AppSettings]
	callOpts := &models.CallOptions{}
	for _, opt := range opts {
		opt(callOpts)
	}
	settings.Id = nil
	settings.CreatedAt = nil
	settings.UpdatedAt = nil
	settings.State = nil

	err := r.Post(ctx, "/apps/create-settings", settings, &res, nil, callOpts)
	return &res, err
}

// UpdateSettings updates an existing settings for an app.
//
// The endpoint will return a slice of Settings structs containing the app settings.
//
// The endpoint will return an error if the request fails.
// The error will contain the status code of the response and the body of the response.
func (r *AppsResource) UpdateSettings(ctx context.Context, settings mahoModels.AppSettings, opts ...models.CallOption) (*mahoModels.MahoragaResponse[mahoModels.AppSettings], error) {
	var user mahoModels.MahoragaResponse[mahoModels.AppSettings]
	callOpts := &models.CallOptions{}
	for _, opt := range opts {
		opt(callOpts)
	}
	err := r.Put(ctx, fmt.Sprintf("/apps/update-settings/%d", settings.Id), settings, &user, nil, callOpts)
	return &user, err
}
