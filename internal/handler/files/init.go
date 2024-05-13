package files

import (
	"github.com/eifzed/joona/internal/config"
	"github.com/eifzed/joona/internal/entity/usecase"
	"github.com/eifzed/joona/lib/common/commonwriter"
	"github.com/eifzed/joona/lib/common/handler"
)

type FileHandler struct {
	Config *config.Config
	FileUC usecase.FileUCInterface
}

func NewFileHandler(recipesHandler *FileHandler) *FileHandler {
	return recipesHandler
}

var (
	bindingBind                   = handler.Bind
	commonwriterRespondError      = commonwriter.RespondError
	commonwriterRespondOK         = commonwriter.RespondOK
	commonwriterRespondOKWithData = commonwriter.RespondOKWithData
)
