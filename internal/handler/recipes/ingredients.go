package recipes

import (
	"net/http"

	"github.com/eifzed/makaji/internal/entity/recipes"
	"github.com/eifzed/makaji/lib/common/commonerr"
)

func (h *RecipesHandler) RegisterIngredients(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	ingredient := recipes.Ingredient{}
	err := bindingBind(r, &ingredient)
	if err != nil {
		err = commonerr.ErrorBadRequest("invalid_params", "invalid registration params")
		commonwriterRespondError(ctx, w, err)
		return
	}
	if err := ingredient.ValidateInput(); err != nil {
		commonwriterRespondError(ctx, w, err)
		return
	}

	err = h.RecipesUC.RegisterIngredient(ctx, ingredient)
	if err != nil {
		commonwriterRespondError(ctx, w, err)
		return
	}
	commonwriterRespondOK(ctx, w)
}

func (h *RecipesHandler) GetIngredients(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	ingredients := recipes.GetIngredientsRequest{
		GenericFilterParams: recipes.GenericFilterParams{
			Limit: 10,
			Page:  1,
		},
	}
	err := bindingBind(r, &ingredients)
	if err != nil {
		err = commonerr.ErrorBadRequest("invalid_params", "invalid registration params")
		commonwriterRespondError(ctx, w, err)
		return
	}

	result, err := h.RecipesUC.GetIngredients(ctx, ingredients)
	if err != nil {
		commonwriterRespondError(ctx, w, err)
		return
	}
	commonwriterRespondOKWithData(ctx, w, result)
}
