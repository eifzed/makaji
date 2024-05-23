package users

import (
	"context"
	"time"

	"github.com/eifzed/joona/lib/utility/hash"
	"github.com/eifzed/joona/lib/utility/jwt"
	"go.mongodb.org/mongo-driver/bson/primitive"

	"github.com/eifzed/joona/internal/constant"
	"github.com/eifzed/joona/internal/entity/users"
	"github.com/eifzed/joona/internal/handler/auth"
	"github.com/eifzed/joona/lib/common/commonerr"
	"github.com/pkg/errors"
)

func (uc *usersUC) RegisterNewUser(ctx context.Context, userDetail users.UserRegistration) (*users.UserAuth, error) {
	isExist, err := uc.usersDB.CheckUserExistByEmail(ctx, userDetail.Email)
	if err != nil {
		return nil, errors.Wrap(err, "RegisterNewUser.CheckUserExistByEmail")
	}
	if isExist {
		return nil, commonerr.ErrorAlreadyExist("email_exist", "user with similar email already exist")
	}
	isExist, err = uc.usersDB.CheckUserExistByUsername(ctx, userDetail.Username)
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
	ctx, err = uc.tx.Start(ctx)
	defer uc.tx.Finish(ctx, &err)

	userData := users.UserDetail{
		Username: userDetail.Username,
		Email:    userDetail.Email,
		Password: passwordHashed,
		FullName: userDetail.FullName,
	}
	err = uc.usersDB.InsertUser(ctx, &userData)
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
		UserID:         userData.UserID.Hex(),
		Email:          userData.Email,
		PasswordHashed: userData.Password,
		FullName:       userData.FullName,
		Username:       userData.Username,
	}
	token, err := jwt.GenerateToken(userPayload, uc.config.Secrets.Data.JWTCertificate.PrivateKey, constant.MinutesInOneDay)
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
	userDetail, err := uc.usersDB.GetUserDetailExistByUsernameOrEmail(ctx, loginData.UsernameOrEmail)
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

func (uc *usersUC) UpdateSelfUser(ctx context.Context, updateUserDetail users.UserProfile) (data users.UserProfile, err error) {
	user, ok := auth.GetUserDetailFromContext(ctx)
	if !ok {
		err = commonerr.DefaultUnauthorized
		return
	}

	oid, err := primitive.ObjectIDFromHex(user.UserID)
	if err != nil {
		err = commonerr.ErrorBadRequest("update_user", "invalid user ID")
		return
	}

	updateUserDetail.UserID = oid

	err = uc.usersDB.UpdateUserByID(ctx, oid, &updateUserDetail)
	if err != nil {
		err = errors.Wrap(err, "UpdateSelfUser.UpdateUserByID")
		return
	}

	err = uc.elastic.UpdateUser(ctx, user.UserID, &users.UserItem{
		UserID:         user.UserID,
		Username:       updateUserDetail.Username,
		FullName:       updateUserDetail.FullName,
		ProfilePicture: updateUserDetail.ProfilePicture,
		Nationality:    updateUserDetail.Nationality,
		Bio:            updateUserDetail.Bio,
	})
	if err != nil {
		err = errors.Wrap(err, "UpdateSelfUser.UpdateUser")
		return
	}
	return updateUserDetail, nil
}

func (uc *usersUC) GetUserByID(ctx context.Context, id primitive.ObjectID) (data users.UserProfile, err error) {
	data, err = uc.usersDB.GetUserProfileByID(ctx, id)
	if err != nil {
		err = errors.Wrap(err, "GetSelfUser.GetUserProfileByID")
		return
	}

	return
}
