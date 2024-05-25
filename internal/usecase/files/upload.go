package files

import (
	"context"

	"github.com/eifzed/makaji/internal/entity/files"
	"github.com/pkg/errors"
)

func (uc *fileUC) UploadFile(ctx context.Context, param files.UploadFileRequest) (resp files.UploadFileResponse, err error) {
	url, err := uc.blob.UploadFile(ctx, param)
	if err != nil {
		err = errors.Wrap(err, "UploadFile")
		return
	}
	resp = files.UploadFileResponse{
		URL: url,
	}
	return
}
