package recipes

import (
	"context"

	"github.com/eifzed/joona/internal/entity/recipes"
	"github.com/eifzed/joona/lib/common/commonerr"
	"github.com/pkg/errors"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"gopkg.in/mgo.v2/bson"
)

func (conn *recipesConn) GetRecipes(ctx context.Context, filter recipes.GetRecipeParams) (result []recipes.Recipe, err error) {
	aggPipeline := bson.M{}
	if filter.ID != "" {
		oid, oErr := primitive.ObjectIDFromHex(filter.ID)
		if oErr != nil {
			// ignore error
			return
		}
		aggPipeline = bson.M{"_id": oid}
	} else {
		calorieFilter := bson.M{}
		if filter.CalorieMin > 0 {
			calorieFilter["$gte"] = filter.CalorieMin
		}
		if filter.CalorieMax > 0 {
			calorieFilter["$lte"] = filter.CalorieMax
		}
		if len(calorieFilter) > 0 {
			aggPipeline["calorie_count"] = calorieFilter
		}

		priceFilter := bson.M{}
		if filter.PriceMin > 0 {
			priceFilter["$gte"] = filter.PriceMin
		}
		if filter.PriceMax > 0 {
			priceFilter["$lte"] = filter.PriceMax
		}
		if len(priceFilter) > 0 {
			aggPipeline["price_estimation"] = priceFilter
		}

		if filter.Keyword != "" {
			aggPipeline["$or"] = []bson.M{
				{"name": bson.M{"$regex": filter.Keyword, "$options": "i"}},
				{"tags": bson.M{"$regex": filter.Keyword, "$options": "i"}},
			}
		}
	}
	options := options.Find().SetSkip(int64((filter.Page - 1) * filter.Limit)).SetLimit(int64(filter.Limit))

	cursor, err := conn.DB.Collection("recipes").Find(ctx, aggPipeline, options)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	for cursor.Next(ctx) {
		var recipe recipes.Recipe
		if err = cursor.Decode(&recipe); err != nil {
			return
		}
		result = append(result, recipe)
	}

	return
}

func (conn *recipesConn) GetRecipeByID(ctx context.Context, id primitive.ObjectID) (recipe recipes.Recipe, err error) {
	result := conn.DB.Collection("recipes").FindOne(ctx, bson.M{"_id": id})

	if result.Err() != nil {
		err = errors.Wrap(result.Err(), "GetRecipeByID.FindOne")
		return
	}
	err = result.Decode(&recipe)
	if err != nil {
		err = errors.Wrap(err, "GetRecipeByID.Decode")
		return
	}
	return
}

func (conn *recipesConn) InsertRecipe(ctx context.Context, recipe *recipes.Recipe) error {
	session := getSessionFromContext(ctx)

	err := mongo.WithSession(ctx, session, func(sc mongo.SessionContext) error {
		var collection *mongo.Collection
		if session != nil {
			collection = getCollectionFromSession(session, "recipes")
		} else {
			collection = conn.DB.Collection("recipes")
		}
		result, err := collection.InsertOne(ctx, recipe)
		if err != nil {
			return err
		}
		recipe.ID = result.InsertedID.(primitive.ObjectID).Hex()
		return nil
	})
	if err != nil {
		return err
	}
	return nil
}

func (conn *recipesConn) UpdateRecipeByID(ctx context.Context, id string, recipe *recipes.Recipe) error {
	session := getSessionFromContext(ctx)

	err := mongo.WithSession(ctx, session, func(sc mongo.SessionContext) error {
		var collection *mongo.Collection
		if session != nil {
			collection = getCollectionFromSession(session, "recipes")
		} else {
			collection = conn.DB.Collection("recipes")
		}
		oid, err := primitive.ObjectIDFromHex(id)
		if err != nil {
			return commonerr.InvalidObjectID
		}

		result, err := collection.UpdateByID(ctx, oid, recipe)
		if err != nil {
			return err
		}
		recipe.ID = result.UpsertedID.(primitive.ObjectID).Hex()
		return nil
	})
	if err != nil {
		return err
	}
	return nil
}
