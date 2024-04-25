package main

import (
	"io"
	"net/http"

	// "github.com/eifzed/joona/lib/utility/urlpath"
	"github.com/eifzed/joona/lib/utility/urlpath"
	"github.com/go-chi/chi"
)

func getRoute(m *modules) *chi.Mux {
	router := chi.NewRouter()
	path := urlpath.New("")
	router.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		io.WriteString(w, "hello inud")
	})
	router.Route("/v1", func(v1 chi.Router) {
		v1.Group(func(noAuthRoute chi.Router) {
			path.Group("/users", func(usersRoute urlpath.Routes) {
				noAuthRoute.Post(usersRoute.URL("/register"), m.httpHandler.UsersHandler.RegisterNewAccount)
			})
		})

		v1.Group(func(authRoute chi.Router) {
			authRoute.Use(m.AuthModule.AuthHandler)
			path.Group("/users", func(usersRoute urlpath.Routes) {
				authRoute.Post(usersRoute.URL("/login"), m.httpHandler.UsersHandler.LoginUser)
			})

			path.Group("/ingredients", func(ingredientsRoute urlpath.Routes) {
				authRoute.Post(ingredientsRoute.URL("/"), m.httpHandler.RecipesHandler.RegisterIngredients)
				authRoute.Get(ingredientsRoute.URL("/"), m.httpHandler.RecipesHandler.GetIngredients)
			})

			path.Group("/recipes", func(recipesRoute urlpath.Routes) {
				authRoute.Post(recipesRoute.URL("/"), m.httpHandler.RecipesHandler.CreateRecipe)
				authRoute.Get(recipesRoute.URL("/"), m.httpHandler.RecipesHandler.GetRecipes)
			})

		})

	})

	return router
}
