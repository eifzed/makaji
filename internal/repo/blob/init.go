package blob

import (
	"fmt"
	"net/url"

	"github.com/Azure/azure-storage-blob-go/azblob"
	"github.com/eifzed/makaji/internal/config"
	"github.com/pkg/errors"
)

type blob struct {
	config      *config.Config
	credentials *azblob.SharedKeyCredential
	service     *azblob.ServiceURL
	accountName string
}

type Option struct {
	AccountName string
	AccountKey  string
	Config      *config.Config
}

func New(opt Option) (*blob, error) {
	credentials, err := azblob.NewSharedKeyCredential(opt.AccountName, opt.AccountKey)
	if err != nil {
		return nil, errors.Wrap(err, "NewSharedKeyCredential")
	}
	pipeline := azblob.NewPipeline(credentials, azblob.PipelineOptions{})

	urlData, err := url.Parse(fmt.Sprintf("https://%s.blob.core.windows.net", opt.AccountName))
	if err != nil {
		return nil, errors.Wrap(err, "url.Parse")
	}

	service := azblob.NewServiceURL(*urlData, pipeline)

	return &blob{
		credentials: credentials,
		config:      opt.Config,
		service:     &service,
		accountName: opt.AccountName,
	}, nil
}
