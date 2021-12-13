package routers

import (
	"github.com/blueworrybear/lineBot/config"
	"github.com/blueworrybear/lineBot/routers/api"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

type Routers struct{
	Config *config.Config
}

func (r *Routers) RegisterRoutes(e *gin.Engine) {
	e.Use(cors.Default())

	apiRoute := &api.Router{
		Config: r.Config,
	}
	apiRoute.RegisterRoutes(e)
}
