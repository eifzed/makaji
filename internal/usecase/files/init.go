package files

import (
	"github.com/eifzed/makaji/internal/config"
	"github.com/eifzed/makaji/internal/entity/repo/external"
)

type fileUC struct {
	config *config.Config
	blob   external.BlobInterface
}

type Options struct {
	Config *config.Config
	Blob   external.BlobInterface
}

func GetNewFileUC(option *Options) *fileUC {
	return &fileUC{
		config: option.Config,
		blob:   option.Blob,
	}
}
