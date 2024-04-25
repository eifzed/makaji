package recipes

import (
	"context"

	"github.com/eifzed/joona/internal/entity/recipes"
)

type RecipesDBInterface interface {
	// ingredients
	InsertIngredients(ctx context.Context, ingredients []recipes.Ingredient) error
	GetIngredients(ctx context.Context, params recipes.GetIngredientsRequest) ([]recipes.Ingredient, error)

	// recipes
	InsertRecipe(ctx context.Context, recipe *recipes.Recipe) error
	GetRecipes(ctx context.Context, filter recipes.GetRecipeParams) (result []recipes.Recipe, err error)
}
