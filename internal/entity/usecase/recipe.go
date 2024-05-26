package usecase

import (
	"context"

	"github.com/eifzed/makaji/internal/entity/recipes"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type RecipesUCInterface interface {
	RegisterIngredient(ctx context.Context, params recipes.Ingredient) (err error)
	GetIngredients(ctx context.Context, params recipes.GetIngredientsRequest) (ingredients []recipes.Ingredient, err error)

	CreateRecipe(ctx context.Context, params recipes.Recipe) (result recipes.GenericPostResponse, err error)
	UpdateRecipe(ctx context.Context, recipeID primitive.ObjectID, params recipes.Recipe) (result recipes.GenericPostResponse, err error)

	GetRecipes(ctx context.Context, params recipes.GetRecipeParams) (result recipes.GetRecipeListResponse, err error)
	GetRecipeDetailByID(ctx context.Context, id primitive.ObjectID) (result recipes.Recipe, err error)
}
