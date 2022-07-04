package main

import (
	"fmt"
	"log"
	"os"

	"github.com/eifzed/joona/internal/config"
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

	fmt.Println(db)
	modules := newModules(modules{
		Config: cfg,
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
