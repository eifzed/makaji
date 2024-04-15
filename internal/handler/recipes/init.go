package recipes

import (
	"github.com/eifzed/joona/internal/config"
	"github.com/eifzed/joona/internal/entity/usecase"
	"github.com/eifzed/joona/lib/common/commonwriter"
	"github.com/eifzed/joona/lib/common/handler"
)

type RecipesHandler struct {
	RecipesUC usecase.RecipesUCInterface
	Config    *config.Config
}

func NewRecipesHandler(recipesHandler *RecipesHandler) *RecipesHandler {
	return recipesHandler
}

var (
	bindingBind                   = handler.Bind
	commonwriterRespondError      = commonwriter.RespondError
	commonwriterRespondOK         = commonwriter.RespondOK
	commonwriterRespondOKWithData = commonwriter.RespondOKWithData
)
