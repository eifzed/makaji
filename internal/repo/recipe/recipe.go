package recipe

import (
	"context"

	"github.com/eifzed/joona/internal/entity/recipe"
	"gopkg.in/mgo.v2/bson"
)

func (db *recipeConn) GetRecipe(ctx context.Context, filter recipe.GetRecipeFilter) (*recipe.Recipe, error) {
	aggPipeline := []bson.M{}
	if !filter.ID.IsZero() {
		aggPipeline = append(aggPipeline, bson.M{"$match": bson.M{"_id": filter.ID}})
	}
	if filter.Tag != "" {
		aggPipeline = append(aggPipeline, bson.M{"$contain": bson.M{"tags": filter.Tag}})
	}
	if filter.Price != nil {

	}
	return nil, nil
}
