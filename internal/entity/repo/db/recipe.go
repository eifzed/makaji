package db

import (
	"context"

	"github.com/eifzed/makaji/internal/entity/recipes"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type RecipesDBInterface interface {
	// ingredients
	InsertIngredients(ctx context.Context, ingredients []recipes.Ingredient) error
	GetIngredients(ctx context.Context, params recipes.GetIngredientsRequest) ([]recipes.Ingredient, error)

	// recipes
	InsertRecipe(ctx context.Context, recipe *recipes.Recipe) error
	GetRecipes(ctx context.Context, filter recipes.GetRecipeParams) (result []recipes.Recipe, err error)
	UpdateRecipeByID(ctx context.Context, id string, recipe *recipes.Recipe) error
	GetRecipeByID(ctx context.Context, id primitive.ObjectID) (recipe recipes.Recipe, err error)
}
