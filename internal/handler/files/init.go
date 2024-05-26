package files

import (
	"github.com/eifzed/makaji/internal/config"
	"github.com/eifzed/makaji/internal/entity/usecase"
	"github.com/eifzed/makaji/lib/common/commonwriter"
	"github.com/eifzed/makaji/lib/common/handler"
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
