package main

import (
	"github.com/blueworrybear/lineBot/core"
	"github.com/blueworrybear/lineBot/models"
	"github.com/google/wire"
	"gorm.io/gorm"
)

var storeSet = wire.NewSet(
	provideDatabaseService,
)

func provideDatabaseService(db *gorm.DB) core.DatabaseService {
	return models.NewDatabaseService(db)
}
