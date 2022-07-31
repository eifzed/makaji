package http

import "github.com/eifzed/joona/internal/entity/handler/http/users"

type HttpHandler struct {
	UsersHandler users.UsersHandler
}
