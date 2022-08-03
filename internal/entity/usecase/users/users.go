package users

import (
	"context"

	"github.com/eifzed/joona/internal/entity/users"
)

type UsersUCInterface interface {
	RegisterNewUser(context context.Context, userDetail users.UserRegistration) (*users.UserAuth, error)
	LoginUser(ctx context.Context, loginData users.UserLogin) (*users.UserAuth, error)
}
