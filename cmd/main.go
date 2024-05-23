package main

import (
	"log"
	"os"

	"github.com/eifzed/joona/internal/config"
	"github.com/eifzed/joona/internal/entity/handler/http"
	"github.com/eifzed/joona/internal/handler/auth"
	fileHttpHandler "github.com/eifzed/joona/internal/handler/files"
	recipesHttpHandler "github.com/eifzed/joona/internal/handler/recipes"
	usersHttpHandler "github.com/eifzed/joona/internal/handler/users"
	"github.com/eifzed/joona/internal/repo/blob"
	"github.com/eifzed/joona/internal/repo/elasticsearch"
	"github.com/eifzed/joona/internal/repo/recipes"
	"github.com/eifzed/joona/internal/repo/redis"
	"github.com/eifzed/joona/internal/repo/users"
	fileUsecase "github.com/eifzed/joona/internal/usecase/files"
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

	blobService, err := blob.New(blob.Option{
		AccountName: cfg.Secrets.Data.AzureBlob.AccountName,
		AccountKey:  cfg.Secrets.Data.AzureBlob.AccountKey,
		Config:      cfg,
	})
	if err != nil {
		log.Fatal(err)
	}

	redisConn := redis.New(redis.Options{
		MaxIdle:       cfg.Redis.MaxIdle,
		MaxActive:     cfg.Redis.MaxActive,
		TimeoutSecond: cfg.Redis.TimeoutSecond,
		AuthKey:       cfg.Secrets.Data.RedisAuth,
		Address:       cfg.Redis.Address,
	})

	usersUC := usersUsecase.GetNewUsersUC(&usersUsecase.Options{
		UsersDB: &usersDB,
		TX:      tx,
		Config:  cfg,
		Elastic: esClient,
	})

	recipesUC := recipesUsecase.GetNewRecipesUC(&recipesUsecase.Options{
		UsersDB:   &usersDB,
		Config:    cfg,
		TX:        tx,
		RecipesDB: &recipesDB,
		Elastic:   esClient,
		Redis:     redisConn,
	})

	fileUC := fileUsecase.GetNewFileUC(&fileUsecase.Options{
		Config: cfg,
		Blob:   blobService,
	})

	usersHadler := usersHttpHandler.NewUsersHandler(&usersHttpHandler.UsersHandler{
		UsersUC: usersUC,
		Config:  cfg,
	})

	recipesHandler := recipesHttpHandler.NewRecipesHandler(&recipesHttpHandler.RecipesHandler{
		RecipesUC: recipesUC,
		Config:    cfg,
	})

	fileHandler := fileHttpHandler.NewFileHandler(&fileHttpHandler.FileHandler{
		Config: cfg,
		FileUC: fileUC,
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
			FileHandler:    fileHandler,
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
