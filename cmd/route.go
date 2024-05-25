package main

import (
	"fmt"
	"net/http"

	// "github.com/eifzed/joona/lib/utility/urlpath"
	"github.com/eifzed/joona/lib/utility/urlpath"
	"github.com/go-chi/chi"
)

func getRoute(m *modules) *chi.Mux {
	router := chi.NewRouter()
	path := urlpath.New("")
	router.Get("/ping", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(fmt.Sprintf("{\"CommitHash\": \"%s\"}", CommitHash)))
	})
	router.Route("/v1", func(v1 chi.Router) {
		v1.Group(func(authRoute chi.Router) {
			authRoute.Use(m.AuthModule.AuthHandler)
			path.Group("/users", func(usersRoute urlpath.Routes) {
				authRoute.Post(usersRoute.URL("/login"), m.httpHandler.UsersHandler.LoginUser)
				authRoute.Post(usersRoute.URL("/register"), m.httpHandler.UsersHandler.RegisterNewAccount)
				authRoute.Get(usersRoute.URL("/{id}"), m.httpHandler.UsersHandler.GetUserByID)
				authRoute.Get(usersRoute.URL("/self"), m.httpHandler.UsersHandler.GetSelfUser)
				authRoute.Get(usersRoute.URL(""), m.httpHandler.UsersHandler.GetUserList)
				authRoute.Put(usersRoute.URL("/self"), m.httpHandler.UsersHandler.UpdateSelfUser)
			})

			path.Group("/ingredients", func(ingredientsRoute urlpath.Routes) {
				authRoute.Post(ingredientsRoute.URL("/"), m.httpHandler.RecipesHandler.RegisterIngredients)
				authRoute.Get(ingredientsRoute.URL("/"), m.httpHandler.RecipesHandler.GetIngredients)
			})

			path.Group("/recipes", func(recipesRoute urlpath.Routes) {
				authRoute.Post(recipesRoute.URL("/"), m.httpHandler.RecipesHandler.CreateRecipe)
				authRoute.Get(recipesRoute.URL("/"), m.httpHandler.RecipesHandler.GetRecipes)
				authRoute.Get(recipesRoute.URL("/{id}"), m.httpHandler.RecipesHandler.GetRecipeDetailByID)
			})

			path.Group("/files", func(filesRoute urlpath.Routes) {
				authRoute.Post(filesRoute.URL("/"), m.httpHandler.FileHandler.UploadFile)
			})

		})

	})

	return router
}
