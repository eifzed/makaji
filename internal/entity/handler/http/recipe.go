package http

import "net/http"

type RecipesHandler interface {
	RegisterIngredients(w http.ResponseWriter, r *http.Request)
	GetIngredients(w http.ResponseWriter, r *http.Request)
}
