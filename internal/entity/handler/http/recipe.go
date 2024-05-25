package http

import "net/http"

type RecipesHandler interface {
	// ingredients
	RegisterIngredients(w http.ResponseWriter, r *http.Request)
	GetIngredients(w http.ResponseWriter, r *http.Request)

	// recipes
	CreateRecipe(w http.ResponseWriter, r *http.Request)
	GetRecipes(w http.ResponseWriter, r *http.Request)
	GetRecipeDetailByID(w http.ResponseWriter, r *http.Request)
}
