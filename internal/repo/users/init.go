package users

import (
	"errors"

	"github.com/eifzed/makaji/lib/database/mongodb/transactions"
	"go.mongodb.org/mongo-driver/mongo"
)

type UsersDB struct {
	DB *mongo.Database
}

type UsersDBOption struct {
	DB *mongo.Database
}

func GetUsersDB(option *UsersDBOption) (UsersDB, error) {
	db := UsersDB{}
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
