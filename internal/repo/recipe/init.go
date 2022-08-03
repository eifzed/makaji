package recipe

import (
	"errors"

	"go.mongodb.org/mongo-driver/mongo"
)

type recipeConn struct {
	DB *mongo.Database
}

type RecipeDBOption struct {
	DB *mongo.Database
}

func GetRecipeDB(option *RecipeDBOption) (recipeConn, error) {
	db := recipeConn{}
	if option == nil || option.DB == nil {
		return db, errors.New("no DB config")
	}
	db.DB = option.DB
	return db, nil
}
