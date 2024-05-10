package usecase

import (
	"context"

	"github.com/eifzed/joona/internal/entity/recipes"
)

type RecipesUCInterface interface {
	RegisterIngredient(ctx context.Context, params recipes.Ingredient) (err error)
	GetIngredients(ctx context.Context, params recipes.GetIngredientsRequest) (ingredients []recipes.Ingredient, err error)

	CreateRecipe(ctx context.Context, params recipes.Recipe) (result recipes.GenericPostResponse, err error)
	GetRecipes(ctx context.Context, params recipes.GetRecipeParams) (result recipes.GetRecipeListResponse, err error)
}
