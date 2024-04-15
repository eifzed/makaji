package recipes

import (
	"errors"

	"github.com/eifzed/joona/lib/database/mongodb/transactions"
	"go.mongodb.org/mongo-driver/mongo"
)

type recipesConn struct {
	DB *mongo.Database
}

type RecipesDBOption struct {
	DB *mongo.Database
}

func GetRecipesDB(option *RecipesDBOption) (recipesConn, error) {
	db := recipesConn{}
	if option == nil || option.DB == nil {
		return db, errors.New("no DB config")
	}
	db.DB = option.DB
	return db, nil
}

var (
	getSessionFromContext    = transactions.GetSessionFromContext
	getCollectionFromSession = transactions.GetCollectionFromSession
)
