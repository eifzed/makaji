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
		v1.Group(func(users chi.Router) {
			users.Use(m.AuthModule.AuthHandler)
			path.Group("/users", func(usersRoute urlpath.Routes) {
				users.Post(usersRoute.URL("/register"), m.httpHandler.UsersHandler.RegisterNewAccount)
				users.Post(usersRoute.URL("/Login"), m.httpHandler.UsersHandler.RegisterNewAccount)
			})
		})
		// 	path.Group("/shops", func(shopsRoute urlpath.Routes) {
		// 		antre.Post(shopsRoute.URL("/register"), m.httpHandler.OrderHandler.RegisterShop)
		// 		path.Group("/products", func(shopProductRoute urlpath.Routes) {
		// 			antre.Get(shopProductRoute.URL("/{shopID}"), m.httpHandler.OrderHandler.GetCustomerOrders)
		// 		})
		// 	})
		// 	path.Group("/orders", func(orderRoute urlpath.Routes) {
		// 		antre.Post(orderRoute.URL(""), m.httpHandler.OrderHandler.RegisterOrder)
		// 		antre.Get(orderRoute.URL(""), m.httpHandler.OrderHandler.GetCustomerOrders)
		// 		antre.Get(orderRoute.URL("/{id}"), m.httpHandler.OrderHandler.GetOrderByID)
		// 	})
		// })

	})

	return router
}
