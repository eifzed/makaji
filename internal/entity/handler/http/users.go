package http

import "net/http"

type UsersHandler interface {
	RegisterNewAccount(w http.ResponseWriter, r *http.Request)
	LoginUser(w http.ResponseWriter, r *http.Request)
	GetSelfUser(w http.ResponseWriter, r *http.Request)
	UpdateSelfUser(w http.ResponseWriter, r *http.Request)
	GetUserByID(w http.ResponseWriter, r *http.Request)
}
