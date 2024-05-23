package db

import (
	"context"

	"github.com/eifzed/joona/internal/entity/recipes"
	"github.com/eifzed/joona/internal/entity/users"
)

type ElasticsearchInterface interface {
	GetRecipeList(ctx context.Context, params recipes.GetRecipeParams) (result recipes.GetRecipeListResponse, err error)
	InsertRecipe(ctx context.Context, data *recipes.ReceipeItem) (err error)
	UpdateRecipe(ctx context.Context, id string, data *recipes.ReceipeItem) (err error)
	InsertUser(ctx context.Context, data *users.UserItem) (err error)
	UpdateUser(ctx context.Context, id string, data *users.UserItem) (err error)
}
