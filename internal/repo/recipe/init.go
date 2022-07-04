package recipe

import (
	"errors"

	"go.mongodb.org/mongo-driver/mongo"
)

type recipeDB struct {
	DB *mongo.Database
}

type OrderDBOption struct {
	DB *mongo.Database
}

func GetNewOrderDB(option *OrderDBOption) (recipeDB, error) {
	db := recipeDB{}
	if option == nil || option.DB == nil {
		return db, errors.New("no DB config")
	}
	db.DB = option.DB
	return db, nil
}
