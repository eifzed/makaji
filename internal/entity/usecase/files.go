package usecase

import (
	"context"

	"github.com/eifzed/joona/internal/entity/files"
)

type FileUCInterface interface {
	UploadFile(ctx context.Context, param files.UploadFileRequest) (response files.UploadFileResponse, err error)
}
