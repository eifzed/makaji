package recipes

import (
	"context"

	"github.com/eifzed/joona/internal/entity/recipes"
)

func (uc *recipesUC) CreateRecipe(ctx context.Context, params recipes.Recipe) (result recipes.GenericPostResponse, err error) {
	err = uc.recipesDB.InsertRecipe(ctx, &params)
	if err != nil {
		return
	}
	result.ID = params.ID
	return
}

func (uc *recipesUC) GetRecipes(ctx context.Context, params recipes.GetRecipeParams) (result []recipes.Recipe, err error) {
	result, err = uc.recipesDB.GetRecipes(ctx, params)
	if err != nil {
		return
	}
	return
}
