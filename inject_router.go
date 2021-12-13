package main

import (
	"github.com/blueworrybear/lineBot/config"
	"github.com/blueworrybear/lineBot/routers"
	"github.com/google/wire"
)

var routerSet = wire.NewSet(
	provideRouter,
)

func provideRouter(Config *config.Config) *routers.Routers {
	return &routers.Routers{
		Config: Config,
	}
}
