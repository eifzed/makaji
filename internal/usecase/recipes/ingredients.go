package recipes

import (
	"context"

	"github.com/eifzed/makaji/internal/entity/recipes"
	"github.com/eifzed/makaji/lib/common/commonerr"
)

func (uc *recipesUC) RegisterIngredient(ctx context.Context, params recipes.Ingredient) (err error) {
	ingredients, err := uc.recipesDB.GetIngredients(ctx, recipes.GetIngredientsRequest{GenericFilterParams: recipes.GenericFilterParams{Keyword: params.Name, Limit: 1, Page: 1}, IsExact: true})
	if err != nil {
		return
	}
	if len(ingredients) > 0 {
		return commonerr.ErrorAlreadyExist("register_ingredient", "ingredient already exists")
	}

	err = uc.recipesDB.InsertIngredients(ctx, []recipes.Ingredient{params})
	if err != nil {
		return
	}
	return nil
}

func (uc *recipesUC) GetIngredients(ctx context.Context, params recipes.GetIngredientsRequest) (ingredients []recipes.Ingredient, err error) {
	ingredients, err = uc.recipesDB.GetIngredients(ctx, params)
	if err != nil {
		return
	}

	return
}
