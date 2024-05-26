package recipes

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/eifzed/makaji/internal/entity/recipes"
	"github.com/eifzed/makaji/internal/handler/auth"
	"github.com/eifzed/makaji/internal/repo/redis"
	"github.com/eifzed/makaji/lib/common"
	"github.com/eifzed/makaji/lib/common/commonerr"
	"github.com/pkg/errors"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func (uc *recipesUC) CreateRecipe(ctx context.Context, params recipes.Recipe) (result recipes.GenericPostResponse, err error) {
	user, exist := auth.GetUserDetailFromContext(ctx)
	if !exist {
		err = commonerr.ErrorUnauthorized("you need to be logged in to post a recipe")
		return
	}
	oid, err := primitive.ObjectIDFromHex(user.UserID)
	if err != nil {
		err = commonerr.ErrorBadRequest("create_recipe", "invalid user id")
		return
	}

	params.CreatorID = oid

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

func (uc *recipesUC) UpdateRecipe(ctx context.Context, recipeID primitive.ObjectID, params recipes.Recipe) (result recipes.GenericPostResponse, err error) {
	user, exist := auth.GetUserDetailFromContext(ctx)
	if !exist {
		err = commonerr.ErrorUnauthorized("you need to be logged in to edit a recipe")
		return
	}
	oid, err := primitive.ObjectIDFromHex(user.UserID)
	if err != nil {
		err = commonerr.ErrorBadRequest("update_recipe", "invalid user id")
		return
	}

	recipe, err := uc.recipesDB.GetRecipeByID(ctx, recipeID)
	if err != nil {
		err = errors.Wrap(err, "recipesDB.GetRecipeByID")
		return
	}
	if recipe.CreatorID.Hex() != user.UserID {
		err = commonerr.ErrorUnauthorized("you are not the creator of this recipe")
		return
	}

	params.CreatorID = oid
	ctx, err = uc.tx.Start(ctx)
	defer uc.tx.Finish(ctx, &err)

	err = uc.recipesDB.UpdateRecipeByID(ctx, recipeID, &params)
	if err != nil {
		err = errors.Wrap(err, "recipesDB.UpdateRecipeByID")
		return
	}
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
	result.ID = recipe.ID

	return
}

const (
	recipeListKey = "recipe-list-%s"
)

func (uc *recipesUC) GetRecipes(ctx context.Context, params recipes.GetRecipeParams) (result recipes.GetRecipeListResponse, err error) {
	paramsB, err := json.Marshal(params)
	if err != nil {
		err = errors.Wrap(err, "json.Marshal")
		return
	}

	hash := common.ComputeSHA256(paramsB)
	key := fmt.Sprintf(recipeListKey, hash)

	cachedList, err := uc.redis.Get(key)
	if err != nil && !errors.Is(err, redis.ErrKeyNotFound) {
		err = errors.Wrap(err, "redis.Get")
		return
	} else if err == nil {
		jErr := json.Unmarshal([]byte(cachedList), &result)
		if jErr == nil {
			return
		}
	}
	result, err = uc.elastic.GetRecipeList(ctx, params)
	if err != nil {
		err = errors.Wrap(err, "GetRecipeList")
		return
	}
	resultB, err := json.Marshal(result)
	if err != nil {
		err = errors.Wrap(err, "json.Marshal")
		return
	}

	recipeListCacheSecond := 120
	if uc.config.CacheExpire.RecipeListSecond > 0 {
		recipeListCacheSecond = uc.config.CacheExpire.RecipeListSecond
	}

	_, err = uc.redis.SetWithExpire(key, string(resultB), recipeListCacheSecond)
	if err != nil {
		err = errors.Wrap(err, "redis.SetWithExpire")
		return
	}
	return
}

func (uc *recipesUC) GetRecipeDetailByID(ctx context.Context, id primitive.ObjectID) (result recipes.Recipe, err error) {
	result, err = uc.recipesDB.GetRecipeByID(ctx, id)
	if err != nil {
		err = errors.Wrap(err, "GetRecipeList")
		return
	}
	return
}
