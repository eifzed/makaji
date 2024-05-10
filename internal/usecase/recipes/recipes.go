package recipes

import (
	"context"

	"github.com/eifzed/joona/internal/entity/recipes"
	"github.com/eifzed/joona/internal/handler/auth"
	"github.com/eifzed/joona/lib/common/commonerr"
	"github.com/pkg/errors"
)

func (uc *recipesUC) CreateRecipe(ctx context.Context, params recipes.Recipe) (result recipes.GenericPostResponse, err error) {
	user, exist := auth.GetUserDetailFromContext(ctx)
	if !exist {
		err = commonerr.ErrorUnauthorized("you need to be logged in to post a recipe")
		return
	}
	params.CreatorName = user.FullName
	params.CreatorUsername = user.Username
	params.CreatorID = user.UserID

	ctx, err = uc.tx.Start(ctx)
	defer uc.tx.Finish(ctx, &err)

	err = uc.recipesDB.InsertRecipe(ctx, &params)
	if err != nil {
		err = errors.Wrap(err, "recipesDB.InsertRecipe")
		return
	}
	result.ID = params.ID
	err = uc.elastic.InsertRecipe(ctx, &recipes.ReceipeItem{
		ID:                params.ID,
		Name:              params.Name,
		Description:       params.Description,
		PriceEstimation:   params.PriceEstimation,
		CountryOrigin:     params.CountryOrigin,
		TimeToCookMinutes: params.TimeToCookMinutes,
		CalorieCount:      params.CalorieCount,
		Difficulty:        params.Difficulty,
		Tags:              params.Tags,
		Tools:             params.Tools,
		CreatorName:       user.FullName,
		CreatorUsername:   user.Username,
		CreatorID:         user.UserID,
	})
	if err != nil {
		err = errors.Wrap(err, "elastic.InsertRecipe")
		return
	}

	return
}

func (uc *recipesUC) UpdateRecipe(ctx context.Context, params recipes.Recipe) (result recipes.GenericPostResponse, err error) {
	user, exist := auth.GetUserDetailFromContext(ctx)
	if !exist {
		err = commonerr.ErrorUnauthorized("you need to be logged in to edit a recipe")
		return
	}
	params.CreatorName = user.FullName
	params.CreatorUsername = user.Username
	params.CreatorID = user.UserID

	ctx, err = uc.tx.Start(ctx)
	defer uc.tx.Finish(ctx, &err)

	err = uc.recipesDB.UpdateRecipeByID(ctx, params.ID, &params)
	if err != nil {
		err = errors.Wrap(err, "recipesDB.UpdateRecipeByID")
		return
	}
	result.ID = params.ID
	err = uc.elastic.UpdateRecipe(ctx, params.ID, &recipes.ReceipeItem{
		ID:                params.ID,
		Name:              params.Name,
		Description:       params.Description,
		PriceEstimation:   params.PriceEstimation,
		CountryOrigin:     params.CountryOrigin,
		TimeToCookMinutes: params.TimeToCookMinutes,
		CalorieCount:      params.CalorieCount,
		Difficulty:        params.Difficulty,
		Tags:              params.Tags,
		Tools:             params.Tools,
		CreatorName:       user.FullName,
		CreatorUsername:   user.Username,
		CreatorID:         user.UserID,
	})
	if err != nil {
		err = errors.Wrap(err, "elastic.UpdateRecipe")
		return
	}

	return
}

func (uc *recipesUC) GetRecipes(ctx context.Context, params recipes.GetRecipeParams) (result recipes.GetRecipeListResponse, err error) {
	result, err = uc.elastic.GetRecipeList(ctx, params)
	if err != nil {
		return
	}
	return
}
