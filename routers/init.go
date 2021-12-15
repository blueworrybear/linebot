package routers

import (
	"github.com/blueworrybear/lineBot/config"
	"github.com/blueworrybear/lineBot/core"
	"github.com/blueworrybear/lineBot/routers/api"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

type Routers struct{
	Config *config.Config
	UserStore core.UserStore
	QuestionStore core.QuestionStore
	ChatStore core.ChatStore
	ChatService core.ChatService
	QuestionService core.QuestionService
}

func (r *Routers) RegisterRoutes(e *gin.Engine) {
	e.Use(cors.Default())

	apiRoute := &api.Router{
		Config: r.Config,
		UserStore: r.UserStore,
		ChatStore: r.ChatStore,
		QuestionStore: r.QuestionStore,
		ChatService: r.ChatService,
		QuestionService: r.QuestionService,
	}
	apiRoute.RegisterRoutes(e)
}
