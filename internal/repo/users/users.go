package users

import (
	"context"

	"github.com/eifzed/joona/internal/entity/users"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"gopkg.in/mgo.v2/bson"
)

func (conn *UsersDB) CheckUserExistByEmail(ctx context.Context, email string) (bool, error) {
	filter := bson.M{"email": email}
	count, err := conn.DB.Collection("users").CountDocuments(ctx, filter)
	if err != nil {
		return false, err
	}
	return count > 0, nil
}

func (conn *UsersDB) CheckUserExistByUsername(ctx context.Context, username string) (bool, error) {
	filter := bson.M{"username": username}
	count, err := conn.DB.Collection("users").CountDocuments(ctx, filter)
	if err != nil {
		return false, err
	}
	return count > 0, nil
}

func (conn *UsersDB) GetUserDetailExistByUsernameOrEmail(ctx context.Context, usernameOrEmail string) (*users.UserDetail, error) {
	filter := bson.M{"$or": []bson.M{{"username": usernameOrEmail}, {"email": usernameOrEmail}}}
	result := conn.DB.Collection("users").FindOne(ctx, filter)
	userDetail := users.UserDetail{}
	err := result.Decode(&userDetail)
	if err != nil {
		return nil, err
	}
	return &userDetail, nil
}

func (conn *UsersDB) InsertUser(ctx context.Context, userDetail *users.UserDetail) error {
	session := getSessionFromContext(ctx)

	err := mongo.WithSession(ctx, session, func(sc mongo.SessionContext) error {
		var collection *mongo.Collection
		if session != nil {
			collection = getCollectionFromSession(session, "users")
		} else {
			collection = conn.DB.Collection("users")
		}
		result, err := collection.InsertOne(ctx, userDetail)
		if err != nil {
			return err
		}
		userDetail.UserID = result.InsertedID.(primitive.ObjectID).Hex()
		return nil
	})
	if err != nil {
		return err
	}
	return nil
}
