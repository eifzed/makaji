package main

import (
	"fmt"
	"log"

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
	ListenAndServe(cfg.Server.HTTP.Address, router)
}
