package upload

import (
	"context"
	"io"
)

type OSS interface {
	// UploadFile save the file and returns the file URL, file name, and error
	// @param category is the file category
	UploadFile(ctx context.Context, category string, reader io.Reader, fileName string) (string, string, error)
	// DeleteFile delete the file and return an error if any, input key is the file name
	DeleteFile(ctx context.Context, category string, key string) error
}
