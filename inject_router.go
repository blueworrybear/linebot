package main

import (
	"github.com/blueworrybear/lineBot/config"
	"github.com/blueworrybear/lineBot/core"
	"github.com/blueworrybear/lineBot/routers"
	"github.com/google/wire"
)

var routerSet = wire.NewSet(
	provideRouter,
)

func provideRouter(
	Config *config.Config,
	UserStore core.UserStore,
	ChatStore core.ChatStore,
	QuestionStore core.QuestionStore,
	ChatService core.ChatService,
	QuestionService core.QuestionService,
) *routers.Routers {
	return &routers.Routers{
		Config: Config,
		UserStore: UserStore,
		ChatStore: ChatStore,
		QuestionStore: QuestionStore,
		ChatService: ChatService,
		QuestionService: QuestionService,
	}
}
