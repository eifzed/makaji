package usecase

import (
	"context"

	"github.com/eifzed/makaji/internal/entity/users"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type UsersUCInterface interface {
	RegisterNewUser(context context.Context, userDetail users.UserRegistration) (*users.UserAuth, error)
	LoginUser(ctx context.Context, loginData users.UserLogin) (*users.UserAuth, error)
	GetUserByID(ctx context.Context, id primitive.ObjectID) (data users.UserProfile, err error)
	UpdateSelfUser(ctx context.Context, updateUserDetail users.UserProfile) (data users.UserProfile, err error)
}
