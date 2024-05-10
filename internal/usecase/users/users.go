package users

import (
	"context"
	"time"

	"github.com/eifzed/joona/lib/utility/hash"
	"github.com/eifzed/joona/lib/utility/jwt"

	"github.com/eifzed/joona/internal/constant"
	"github.com/eifzed/joona/internal/entity/users"
	"github.com/eifzed/joona/lib/common/commonerr"
	"github.com/pkg/errors"
)

func (uc *usersUC) RegisterNewUser(ctx context.Context, userDetail users.UserRegistration) (*users.UserAuth, error) {
	isExist, err := uc.UsersDB.CheckUserExistByEmail(ctx, userDetail.Email)
	if err != nil {
		return nil, errors.Wrap(err, "RegisterNewUser.CheckUserExistByEmail")
	}
	if isExist {
		return nil, commonerr.ErrorAlreadyExist("email_exist", "user with similar email already exist")
	}
	isExist, err = uc.UsersDB.CheckUserExistByUsername(ctx, userDetail.Username)
	if err != nil {
		return nil, errors.Wrap(err, "RegisterNewUser.CheckUserExistByUsername")
	}
	if isExist {
		return nil, commonerr.ErrorAlreadyExist("username_exist", "username is taken")
	}

	passwordHashed, err := hash.HashPassword(userDetail.Password)
	if err != nil {
		return nil, err
	}
	ctx, err = uc.TX.Start(ctx)
	defer uc.TX.Finish(ctx, &err)

	userData := users.UserDetail{
		Username: userDetail.Username,
		Email:    userDetail.Email,
		Password: passwordHashed,
		FullName: userDetail.FullName,
	}
	err = uc.UsersDB.InsertUser(ctx, &userData)
	if err != nil {
		return nil, err
	}
	userAuth, err := uc.getUserAuth(userData)
	if err != nil {
		return nil, err
	}

	return userAuth, nil
}

func (uc *usersUC) getUserAuth(userData users.UserDetail) (*users.UserAuth, error) {
	userPayload := jwt.JWTPayload{
		UserID:         userData.UserID,
		Email:          userData.Email,
		PasswordHashed: userData.Password,
		FullName:       userData.Email,
		Username:       userData.Username,
	}
	token, err := jwt.GenerateToken(userPayload, uc.Config.Secrets.Data.JWTCertificate.PrivateKey, constant.MinutesInOneDay)
	if err != nil {
		return nil, err
	}
	tNow := time.Now()
	return &users.UserAuth{
		Token:      token,
		ValidUntil: tNow.Add(time.Duration(constant.MinutesInOneDay) * time.Minute),
	}, nil
}

func (uc *usersUC) LoginUser(ctx context.Context, loginData users.UserLogin) (*users.UserAuth, error) {
	userDetail, err := uc.UsersDB.GetUserDetailExistByUsernameOrEmail(ctx, loginData.UsernameOrEmail)
	if err != nil {
		return nil, errors.Wrap(err, "LoginUser.GetUserDetailExistByUsernameOrEmail")
	}
	if !hash.IsValidPasswordHash(loginData.Password, userDetail.Password) {
		return nil, commonerr.ErrorForbidden("invalid password")
	}
	userAuth, err := uc.getUserAuth(*userDetail)
	if err != nil {
		return nil, err
	}
	return userAuth, nil
}
