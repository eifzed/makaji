package users

import (
	"github.com/eifzed/joona/internal/config"
	"github.com/eifzed/joona/internal/entity/usecase/users"
	"github.com/eifzed/joona/lib/common/commonwriter"
	"github.com/eifzed/joona/lib/common/handler"
)

type UsersHandler struct {
	UsersUC users.UsersUCInterface
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
