package users

import (
	"context"

	"github.com/eifzed/joona/internal/entity/users"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type UsersDBInterface interface {
	CheckUserExistByEmail(ctx context.Context, email string) (bool, error)
	CheckUserExistByUsername(ctx context.Context, username string) (bool, error)
	GetUserDetailExistByUsernameOrEmail(ctx context.Context, usernameOrEmail string) (*users.UserDetail, error)
	InsertUser(ctx context.Context, userDetail *users.UserDetail) error
	UpdateUserByID(ctx context.Context, userID primitive.ObjectID, userDetail *users.UserProfile) (err error)
	GetUserProfileByEmail(ctx context.Context, email string) (user users.UserProfile, err error)
	GetUserProfileByID(ctx context.Context, id primitive.ObjectID) (user users.UserProfile, err error)
}
