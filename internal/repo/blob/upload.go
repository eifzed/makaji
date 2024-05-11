package blob

import (
	"context"
	"fmt"

	"github.com/Azure/azure-storage-blob-go/azblob"
	"github.com/eifzed/joona/internal/entity/files"
	"github.com/pkg/errors"
)

const (
	urlTemplate = "https://%s.blob.core.windows.net/%s/%s"
)

func (b *blob) UploadFile(ctx context.Context, param files.UploadFileRequest) (url string, err error) {
	containerURL := b.service.NewContainerURL(param.Container)

	blobURL := containerURL.NewBlockBlobURL(param.Filename)
	_, err = azblob.UploadBufferToBlockBlob(ctx, param.File, blobURL, azblob.UploadToBlockBlobOptions{})
	if err != nil {
		err = errors.Wrap(err, "UploadBufferToBlockBlob")
	}

	url = fmt.Sprintf(urlTemplate, b.accountName, param.Container, param.Filename)
	return
}
