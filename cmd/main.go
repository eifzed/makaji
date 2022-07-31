package main

import (
	"log"
	"os"

	"github.com/eifzed/joona/internal/config"
	"github.com/eifzed/joona/internal/entity/handler/http"
	"github.com/eifzed/joona/internal/handler/auth"
	usersHttpHandler "github.com/eifzed/joona/internal/handler/users"
	"github.com/eifzed/joona/internal/repo/users"
	usersUsecase "github.com/eifzed/joona/internal/usecase/users"
	"github.com/eifzed/joona/lib/database/mongodb/transactions"
)

func main() {
	secret := config.GetSecrets()
	if secret == nil {
		log.Fatal("failed to get secrets")
	}
	cfg, err := config.GetConfig()
	if err != nil {
		log.Fatal(err)
	}
	cfg.Secrets = secret
	client, err := getDBConnection(cfg)
	if err != nil {
		log.Fatal(err)
	}
	db := client.Database("joona-db")

	usersDB, err := users.GetUsersDB(&users.UsersDBOption{
		DB: db,
	})
	if err != nil {
		log.Fatal(err)
	}
	tx := transactions.GetNewMongoDBTransaction(&transactions.Options{
		Client: client,
	})
	usersUC := usersUsecase.GetNewUsersUC(&usersUsecase.Options{
		UsersDB: &usersDB,
		TX:      tx,
	})
	usersHadler := usersHttpHandler.NewUsersHandler(&usersHttpHandler.UsersHandler{
		UsersUC: usersUC,
		Config:  cfg,
	})
	authModule := auth.NewAuthModule(&auth.AuthModule{
		JWTCertificate: cfg.Secrets.Data.JWTCertificate,
		RouteRoles:     cfg.RouteRoles,
		Cfg:            cfg,
	})
	modules := newModules(modules{
		Config: cfg,
		httpHandler: &http.HttpHandler{
			UsersHandler: usersHadler,
		},
		AuthModule: authModule,
	})
	router := getRoute(modules)
	port := os.Getenv("PORT")
	if port != "" {
		cfg.Server.HTTP.Address = ":" + port
	}
	err = ListenAndServe(cfg.Server.HTTP.Address, router)
	if err != nil {
		log.Println("application exited with error: ", err)
	}
}
