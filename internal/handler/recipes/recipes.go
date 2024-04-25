package recipes

import (
	"net/http"

	"github.com/eifzed/joona/internal/entity/recipes"
	"github.com/eifzed/joona/lib/common/commonerr"
)

func (h *RecipesHandler) CreateRecipe(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	recipe := recipes.Recipe{}
	err := bindingBind(r, &recipe)
	if err != nil {
		err = commonerr.ErrorBadRequest("invalid_params", "invalid params")
		commonwriterRespondError(ctx, w, err)
		return
	}
	if err := recipe.ValidateInput(); err != nil {
		commonwriterRespondError(ctx, w, err)
		return
	}

	result, err := h.RecipesUC.CreateRecipe(ctx, recipe)
	if err != nil {
		commonwriterRespondError(ctx, w, err)
		return
	}
	commonwriterRespondOKWithData(ctx, w, result)
}

func (h *RecipesHandler) GetRecipes(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	params := recipes.GetRecipeParams{
		GenericFilterParams: recipes.GenericFilterParams{
			Limit: 10,
			Page:  1,
		},
	}
	err := bindingBind(r, &params)
	if err != nil {
		err = commonerr.ErrorBadRequest("invalid_params", "invalid params")
		commonwriterRespondError(ctx, w, err)
		return
	}

	result, err := h.RecipesUC.GetRecipes(ctx, params)
	if err != nil {
		commonwriterRespondError(ctx, w, err)
		return
	}
	commonwriterRespondOKWithData(ctx, w, result)
}