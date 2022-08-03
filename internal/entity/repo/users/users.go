package users

import (
	"context"

	"github.com/eifzed/joona/internal/entity/users"
)

type UsersDBInterface interface {
	CheckUserExistByEmail(ctx context.Context, email string) (bool, error)
	CheckUserExistByUsername(ctx context.Context, username string) (bool, error)
	GetUserDetailExistByUsernameOrEmail(ctx context.Context, usernameOrEmail string) (*users.UserDetail, error)
	InsertUser(ctx context.Context, userDetail *users.UserDetail) error
}
