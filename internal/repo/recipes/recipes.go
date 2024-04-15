package recipes

import (
	"context"

	"github.com/eifzed/joona/internal/entity/recipes"
	"gopkg.in/mgo.v2/bson"
)

func (db *recipesConn) GetRecipes(ctx context.Context, filter recipes.GetRecipeFilter) (*recipes.Recipe, error) {
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
