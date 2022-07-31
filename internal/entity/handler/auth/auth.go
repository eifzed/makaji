package auth

import "net/http"

type AuthModuleInterface interface {
	AuthHandler(next http.Handler) http.Handler
}
