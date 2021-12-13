//go:build wireinject
// +build wireinject

package main

import (
	"github.com/blueworrybear/lineBot/config"
	"github.com/google/wire"
	"gorm.io/gorm"
)

func InitializeApplication(cfg *config.Config, db *gorm.DB) (application, error){
	wire.Build(
		storeSet,
		routerSet,
		newApplication,
	)
	return application{}, nil
}
