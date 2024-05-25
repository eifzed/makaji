package external

import (
	"context"

	"github.com/eifzed/makaji/internal/entity/files"
)

type BlobInterface interface {
	UploadFile(ctx context.Context, param files.UploadFileRequest) (url string, err error)
}
