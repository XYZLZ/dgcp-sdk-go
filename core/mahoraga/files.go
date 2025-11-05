package mahoraga

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"mime/multipart"

	sdkClient "github.com/XYZLZ/dgcp-sdk-go/client"
	"github.com/XYZLZ/dgcp-sdk-go/models"
	mahoModels "github.com/XYZLZ/dgcp-sdk-go/models/mahoraga"
)

type FilesResource struct {
	*sdkClient.BaseClient
}

func NewFilesResource(config *sdkClient.SDKConfig) *FilesResource {
	return &FilesResource{
		BaseClient: sdkClient.NewBaseClient(config, sdkClient.Mahoraga),
	}
}

// List returns a list of files.
//
// If params is nil, the endpoint will return all files.
// If params is not nil, the endpoint will return a paginated list of files.
// The Page and Limit fields of params are used to paginate the list.
// The endpoint will return a slice of FilesInfo structs:
//
// The endpoint will return an error if the request fails.
// The error will contain the status code of the response and the body of the response.
func (r *FilesResource) List(ctx context.Context, params *models.PaginationRequest) (*mahoModels.MahhoragaPaginatedResponse[[]mahoModels.FilesInfo], error) {
	var (
		result mahoModels.MahhoragaPaginatedResponse[[]mahoModels.FilesInfo]
		path   = "/files"
	)

	if params != nil {
		path = fmt.Sprintf("%s?page=%d&limit=%d", path, params.Page, params.Limit)
	}

	err := r.Get(ctx, path, &result)
	return &result, err
}

// Upload allows you to upload files to the server.
//
// The endpoint will return a slice of FilesInfo structs:
//
// The endpoint will return an error if the request fails.
// The error will contain the status code of the response and the body of the response.
func (r *FilesResource) Upload(ctx context.Context, files []*multipart.FileHeader) (*mahoModels.MahoragaResponse[[]mahoModels.FilesInfo], error) {
	var (
		result mahoModels.MahoragaResponse[[]mahoModels.FilesInfo]
		path   = "/files/upload"
		b      bytes.Buffer
		writer = multipart.NewWriter(&b)
	)

	for _, fileHeader := range files {
		file, err := fileHeader.Open()
		if err != nil {
			return nil, err
		}
		defer file.Close()

		part, err := writer.CreateFormFile("files", fileHeader.Filename)
		if err != nil {
			return nil, fmt.Errorf("failed to create form file: %v", err)
		}

		_, err = io.Copy(part, file)
		if err != nil {
			return nil, fmt.Errorf("failed to copy file data: %v", err)
		}
	}

	if err := writer.Close(); err != nil {
		return nil, fmt.Errorf("failed to close multipart writer: %v", err)
	}

	err := r.Post(ctx, path, &b, &result)
	return &result, err
}

// Download returns the contents of a file.
//
// The endpoint will return a byte array containing the contents of the file.
//
// The endpoint will return an error if the request fails.
// The error will contain the status code of the response and the body of the response.
func (r *FilesResource) Download(ctx context.Context, fileId string) ([]byte, error) {
	var result []byte

	err := r.Get(ctx, fmt.Sprintf("/files/download/%s", fileId), &result)
	return result, err
}

func (r *FilesResource) Delete(ctx context.Context, fileId string) error {
	path := fmt.Sprintf("/files/delete/%s", fileId)
	return r.BaseClient.Delete(ctx, path)
}
