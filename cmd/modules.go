package main

import (
	"github.com/eifzed/joona/internal/config"
	"github.com/eifzed/joona/internal/entity/handler/auth"
	"github.com/eifzed/joona/internal/entity/handler/http"
)

type modules struct {
	httpHandler *http.HttpHandler
	Config      *config.Config
	AuthModule  auth.AuthModuleInterface
	// LogModule   handler.LogModuleInterface
}

func newModules(mod modules) *modules {
	return &mod
}
