package users

import "net/http"

type UsersHandler interface {
	RegisterNewAccount(w http.ResponseWriter, r *http.Request)
	LoginUser(w http.ResponseWriter, r *http.Request)
}
