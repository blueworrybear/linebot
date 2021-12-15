package main

import (
	"github.com/blueworrybear/lineBot/core"
	"github.com/blueworrybear/lineBot/models"
	"github.com/google/wire"
	"gorm.io/gorm"
)

var storeSet = wire.NewSet(
	provideDatabaseService,
	provideUserStore,
	provideChatStore,
	provideQuestionStore,
)

func provideDatabaseService(db *gorm.DB) core.DatabaseService {
	return models.NewDatabaseService(db)
}

func provideUserStore(db core.DatabaseService) core.UserStore{
	return models.NewUserStore(db)
}

func provideChatStore(db core.DatabaseService) core.ChatStore {
	return models.NewChatStore(db)
}

func provideQuestionStore(db core.DatabaseService) core.QuestionStore {
	return models.NewQuestionStore(db)
}
