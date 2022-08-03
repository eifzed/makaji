package auth

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"strings"

	"github.com/eifzed/joona/internal/config"
	"github.com/eifzed/joona/internal/entity/users"
	"github.com/eifzed/joona/lib/common/commonerr"
	"github.com/eifzed/joona/lib/common/commonwriter"
	"github.com/eifzed/joona/lib/utility/jwt"
	"github.com/go-chi/chi"
)

type AuthModule struct {
	JWTCertificate *jwt.JWTCertificate
	RouteRoles     map[string]jwt.RouteRoles
	Cfg            *config.Config
}

// fieldInfo is getter/setter value from the Info Context
type fieldInfo struct{}

type userContext struct{}

var (
	userContextKey = userContext{}
)

type Info struct {
	UserID int64
	Type   string
	Data   map[string]interface{}
}

func NewAuthModule(module *AuthModule) *AuthModule {
	return module
}

func (m *AuthModule) AuthHandler(next http.Handler) http.Handler {
	return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		rCtx := chi.RouteContext(r.Context())
		if rCtx == nil {
			authHandlerError(ctx, rw, r, errors.New("context is not Chi context"))
			return
		}
		route := fmt.Sprintf("%s %s", rCtx.RouteMethod, rCtx.RoutePattern())
		roles := m.RouteRoles[route].Roles

		if !isPublicRoute(roles, m.Cfg.Roles.Public) {
			bearerToken := r.Header.Get("Authorization")
			jwtToken, err := GetBearerToken(bearerToken)
			if err != nil {
				authHandlerError(ctx, rw, r, err)
				return
			}
			userPayload, err := jwt.DecodeToken(jwtToken, m.JWTCertificate.PublicKey)
			if err != nil {
				authHandlerError(ctx, rw, r, err)
				return
			}

			if !isUserAuthorized(userPayload.Roles, m.RouteRoles[route].Roles) {
				authHandlerError(ctx, rw, r, jwt.ErrForbidden)
				return
			}
			ctx = m.SetKeyValueToContext(ctx, userContextKey, userPayload)
		}

		r = r.WithContext(ctx)
		next.ServeHTTP(rw, r)
	})
}

func (m *AuthModule) SetKeyValueToContext(ctx context.Context, key interface{}, value interface{}) context.Context {
	if ctx == nil {
		ctx = context.Background()
	}
	return context.WithValue(ctx, key, value)
}

func GetUserDetailFromContext(ctx context.Context) (jwt.JWTPayload, bool) {
	user, exist := ctx.Value(userContextKey).(jwt.JWTPayload)
	return user, exist
}

func isUserAuthorized(userRoles []users.UserRole, authorizedRoles []users.UserRole) bool {
	if len(userRoles) == 0 || len(authorizedRoles) == 0 {
		return false
	}
	for _, user := range userRoles {
		for _, auth := range authorizedRoles {
			if user.ID == auth.ID {
				return true
			}
		}
	}
	return false
}

func isPublicRoute(roles []users.UserRole, publicRoleID int64) bool {
	for _, role := range roles {
		if role.ID == publicRoleID {
			return true
		}
	}
	return false
}
func GetBearerToken(token string) (string, error) {
	data := strings.Split(token, "Bearer ")
	if len(data) != 2 {
		return "", jwt.ErrInvalid
	}
	return data[1], nil
}

func authHandlerError(ctx context.Context, w http.ResponseWriter, r *http.Request, err error) {
	switch err {
	case jwt.ErrInvalid:
		err := commonerr.ErrorUnauthorized(err.Error())
		commonwriter.RespondError(ctx, w, err)
	case jwt.ErrExpired:
		err := commonerr.ErrorUnauthorized(err.Error())
		commonwriter.RespondError(ctx, w, err)
	case jwt.ErrForbidden:
		err := commonerr.ErrorForbidden(err.Error())
		commonwriter.RespondError(ctx, w, err)
	default:
		commonwriter.RespondError(ctx, w, err)
	}
}
