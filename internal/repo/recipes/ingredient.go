package recipes

import (
	"context"
	"fmt"

	"github.com/eifzed/makaji/internal/entity/recipes"
	"github.com/eifzed/makaji/lib/utility"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"gopkg.in/mgo.v2/bson"
)

func (conn *recipesConn) InsertIngredients(ctx context.Context, ingredients []recipes.Ingredient) error {
	session := getSessionFromContext(ctx)

	err := mongo.WithSession(ctx, session, func(sc mongo.SessionContext) error {
		var collection *mongo.Collection
		if session != nil {
			collection = getCollectionFromSession(session, "ingredients")
		} else {
			collection = conn.DB.Collection("ingredients")
		}
		result, err := collection.InsertMany(ctx, utility.ConvertSliceToSliceOfInterface(ingredients))
		if err != nil {
			return err
		}
		for i := range ingredients {
			if i >= len(result.InsertedIDs) {
				break
			}
			ingredients[i].ID = result.InsertedIDs[i].(primitive.ObjectID).Hex()
		}
		return nil
	})
	if err != nil {
		return err
	}
	return nil
}

func (conn *recipesConn) GetIngredients(ctx context.Context, params recipes.GetIngredientsRequest) (ingredients []recipes.Ingredient, err error) {
	var filter bson.M

	if params.IngredientID != "" {
		filter = bson.M{"_id": params.IngredientID}
	} else if params.Keyword != "" {
		if params.IsExact {
			filter = bson.M{"$or": []bson.M{
				{"name": bson.M{"$regex": fmt.Sprintf("^%s$", params.Keyword), "$options": "i"}},
				{"alternative_names": bson.M{"$regex": fmt.Sprintf("^%s$", params.Keyword), "$options": "i"}},
			}}
		} else {
			filter = bson.M{"$or": []bson.M{
				{"name": bson.M{"$regex": params.Keyword, "$options": "i"}},
				{"alternative_names": bson.M{"$regex": params.Keyword, "$options": "i"}},
			}}
		}

	}

	options := options.Find().SetSkip(int64((params.Page - 1) * params.Limit)).SetLimit(int64(params.Limit))

	cursor, err := conn.DB.Collection("ingredients").Find(ctx, filter, options)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	ingredients = []recipes.Ingredient{}

	for cursor.Next(ctx) {
		var ingredient recipes.Ingredient
		if err = cursor.Decode(&ingredient); err != nil {
			return
		}
		ingredients = append(ingredients, ingredient)
	}
	return ingredients, nil
}
