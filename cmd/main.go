package main

import (
	"log"
	"os"

	"github.com/eifzed/joona/internal/config"
	"github.com/eifzed/joona/internal/entity/handler/http"
	"github.com/eifzed/joona/internal/handler/auth"
	recipesHttpHandler "github.com/eifzed/joona/internal/handler/recipes"
	usersHttpHandler "github.com/eifzed/joona/internal/handler/users"
	"github.com/eifzed/joona/internal/repo/elasticsearch"
	"github.com/eifzed/joona/internal/repo/recipes"
	"github.com/eifzed/joona/internal/repo/users"
	recipesUsecase "github.com/eifzed/joona/internal/usecase/recipes"
	usersUsecase "github.com/eifzed/joona/internal/usecase/users"
	"github.com/eifzed/joona/lib/database/mongodb/transactions"
)

var CommitHash string

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
	recipesDB, err := recipes.GetRecipesDB(&recipes.RecipesDBOption{
		DB: db,
	})
	if err != nil {
		log.Fatal(err)
	}
	tx := transactions.GetNewMongoDBTransaction(&transactions.Options{
		Client: client,
		DBName: "joona-db",
	})

	// elasaticsearch
	esClient, err := elasticsearch.New(elasticsearch.Option{
		CloudID: cfg.Secrets.Data.Elasticsearch.CloudID,
		APIKey:  cfg.Secrets.Data.Elasticsearch.APIKey,
		Config:  cfg,
	})

	if err != nil {
		log.Fatal(err)
	}

	usersUC := usersUsecase.GetNewUsersUC(&usersUsecase.Options{
		UsersDB: &usersDB,
		TX:      tx,
		Config:  cfg,
	})

	recipesUC := recipesUsecase.GetNewRecipesUC(&recipesUsecase.Options{
		UsersDB:   &usersDB,
		Config:    cfg,
		TX:        tx,
		RecipesDB: &recipesDB,
		Elastic:   esClient,
	})

	usersHadler := usersHttpHandler.NewUsersHandler(&usersHttpHandler.UsersHandler{
		UsersUC: usersUC,
		Config:  cfg,
	})

	recipesHandler := recipesHttpHandler.NewRecipesHandler(&recipesHttpHandler.RecipesHandler{
		RecipesUC: recipesUC,
		Config:    cfg,
	})

	authModule := auth.NewAuthModule(&auth.AuthModule{
		JWTCertificate: cfg.Secrets.Data.JWTCertificate,
		RouteRoles:     cfg.RouteRoles,
		Cfg:            cfg,
	})
	modules := newModules(modules{
		Config: cfg,
		httpHandler: &http.HttpHandler{
			UsersHandler:   usersHadler,
			RecipesHandler: recipesHandler,
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
