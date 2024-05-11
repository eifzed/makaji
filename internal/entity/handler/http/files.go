package http

import "net/http"

type FileHandler interface {
	UploadFile(w http.ResponseWriter, r *http.Request)
}
