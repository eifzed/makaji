package recipe

// func (db *recipeDB) GetRecipe(ctx context.Context, filter recipe.GetRecipeFilter) (*recipe.Recipe, error) {
// 	aggPipeline := []bson.M{}
// 	if !filter.ID.IsZero() {
// 		aggPipeline = append(aggPipeline, bson.M{"$match": bson.M{"_id": filter.ID}})
// 	}
// 	if filter.Tag != "" {
// 		aggPipeline = append(aggPipeline, bson.M{"$contain": bson.M{"tags": filter.Tag}})
// 	}
// 	if filter.Price != nil {

// 	}
// }
