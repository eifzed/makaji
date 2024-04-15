package recipes

import (
	"context"

	"github.com/eifzed/joona/internal/entity/recipes"
)

type RecipesDBInterface interface {
	InsertIngredients(ctx context.Context, ingredients []recipes.Ingredient) error
	GetIngredients(ctx context.Context, params recipes.GetIngredientsRequest) ([]recipes.Ingredient, error)
}
