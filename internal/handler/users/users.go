package users

import (
	"net/http"

	"github.com/eifzed/joona/internal/entity/users"
	"github.com/eifzed/joona/internal/handler/auth"
	"github.com/eifzed/joona/lib/common/commonerr"
	"github.com/go-chi/chi"
	"github.com/leebenson/conform"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func (h *UsersHandler) RegisterNewAccount(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	registerData := users.UserRegistration{}
	err := bindingBind(r, &registerData)
	if err != nil {
		err = commonerr.ErrorBadRequest("invalid_params", "invalid registration params")
		commonwriterRespondError(ctx, w, err)
		return
	}
	conform.Strings(registerData)
	if err := registerData.ValidateInput(); err != nil {
		commonwriterRespondError(ctx, w, err)
		return
	}

	auth, err := h.UsersUC.RegisterNewUser(ctx, registerData)
	if err != nil {
		commonwriterRespondError(ctx, w, err)
		return
	}
	commonwriterRespondOKWithData(ctx, w, auth)
}

func (h *UsersHandler) LoginUser(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	loginData := users.UserLogin{}
	err := bindingBind(r, &loginData)
	if err != nil {
		err = commonerr.ErrorBadRequest("invalid_params", "invalid login params")
		commonwriterRespondError(ctx, w, err)
		return
	}

	auth, err := h.UsersUC.LoginUser(ctx, loginData)
	if err != nil {
		commonwriterRespondError(ctx, w, err)
		return
	}
	commonwriterRespondOKWithData(ctx, w, auth)
}

func (h *UsersHandler) GetSelfUser(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	userContext, ok := auth.GetUserDetailFromContext(ctx)
	if !ok {
		err := commonerr.DefaultUnauthorized
		commonwriterRespondError(ctx, w, err)
		return
	}

	id, err := primitive.ObjectIDFromHex(userContext.UserID)
	if err != nil {
		err = commonerr.ErrorNotFound("user")
		commonwriterRespondError(ctx, w, err)
		return
	}

	auth, err := h.UsersUC.GetUserByID(ctx, id)
	if err != nil {
		commonwriterRespondError(ctx, w, err)
		return
	}
	commonwriterRespondOKWithData(ctx, w, auth)
}

func (h *UsersHandler) GetUserByID(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	userID := chi.URLParam(r, "id")

	oid, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		err = commonerr.ErrorNotFound("user")
		commonwriterRespondError(ctx, w, err)
		return
	}

	auth, err := h.UsersUC.GetUserByID(ctx, oid)
	if err != nil {
		commonwriterRespondError(ctx, w, err)
		return
	}
	commonwriterRespondOKWithData(ctx, w, auth)
}

func (h *UsersHandler) UpdateSelfUser(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	updateData := users.UserProfile{}
	err := bindingBind(r, &updateData)
	if err != nil {
		err = commonerr.ErrorBadRequest("update_self_user", "invalid update params")
		commonwriterRespondError(ctx, w, err)
		return
	}

	auth, err := h.UsersUC.UpdateSelfUser(ctx, updateData)
	if err != nil {
		commonwriterRespondError(ctx, w, err)
		return
	}
	commonwriterRespondOKWithData(ctx, w, auth)
}

func (h *UsersHandler) GetUserList(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	params := users.GenericFilterParams{
		Limit: 10,
		Page:  1,
	}
	err := bindingBind(r, &params)
	if err != nil {
		err = commonerr.ErrorBadRequest("invalid_params", "invalid params")
		commonwriterRespondError(ctx, w, err)
		return
	}

	auth, err := h.UsersUC.GetUserList(ctx, params)
	if err != nil {
		commonwriterRespondError(ctx, w, err)
		return
	}
	commonwriterRespondOKWithData(ctx, w, auth)
}
