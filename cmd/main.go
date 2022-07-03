package main

import (
	"fmt"
	"log"
	"os"

	"github.com/eifzed/joona/internal/config"
)

func main() {
	fmt.Println("hello world")
	secrete := config.GetSecretes()
	if secrete == nil {
		log.Fatal("failed to get secretes")
		return
	}
	cfg, err := config.GetConfig()
	if err != nil {
		log.Fatal(err)
	}
	cfg.Secretes = secrete
	modules := newModules(modules{
		Config: cfg,
	})
	router := getRoute(modules)
	port := os.Getenv("PORT")
	if port != "" {
		cfg.Server.HTTP.Address = port
	}
	err = ListenAndServe(cfg.Server.HTTP.Address, router)
	if err != nil {
		log.Println("application exited with error: ", err)
	}
}
