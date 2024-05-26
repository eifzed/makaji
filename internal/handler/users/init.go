package users

import (
	"github.com/eifzed/makaji/internal/config"
	"github.com/eifzed/makaji/internal/entity/usecase"
	"github.com/eifzed/makaji/lib/common/commonwriter"
	"github.com/eifzed/makaji/lib/common/handler"
)

type UsersHandler struct {
	UsersUC usecase.UsersUCInterface
	Config  *config.Config
}

func NewUsersHandler(orderHandler *UsersHandler) *UsersHandler {
	return orderHandler
}

var (
	bindingBind                   = handler.Bind
	commonwriterRespondError      = commonwriter.RespondError
	commonwriterRespondOK         = commonwriter.RespondOK
	commonwriterRespondOKWithData = commonwriter.RespondOKWithData
)
