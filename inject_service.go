package main

import (
	"github.com/blueworrybear/lineBot/config"
	"github.com/blueworrybear/lineBot/core"
	"github.com/blueworrybear/lineBot/service"
	"github.com/google/wire"
)

var serviceSet = wire.NewSet(
	provideChatService,
	provideQuestionService,
)

func provideChatService(cfg *config.Config, userStore core.UserStore, chatStore core.ChatStore) core.ChatService {
	return service.NewChatService(cfg, userStore,chatStore)
}

func provideQuestionService(
	userStore core.UserStore,
	questionStore core.QuestionStore,
	chatService core.ChatService,
) core.QuestionService {
	return service.NewQuestionService(userStore, questionStore, chatService)
}
