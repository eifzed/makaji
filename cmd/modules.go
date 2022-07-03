package main

import (
	"github.com/eifzed/joona/internal/config"
	// "github.com/eifzed/antre-app/internal/handler"
)

type modules struct {
	// httpHandler *handler.HttpHandler
	Config *config.Config
	// AuthModule  handler.AuthModuleInterface
	// LogModule   handler.LogModuleInterface
}

func newModules(mod modules) *modules {
	return &mod
}
