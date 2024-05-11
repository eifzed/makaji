package external

import (
	"context"

	"github.com/eifzed/joona/internal/entity/files"
)

type BlobInterface interface {
	UploadFile(ctx context.Context, param files.UploadFileRequest) (url string, err error)
}
