package api

import (
	"log"

	"github.com/blueworrybear/lineBot/config"
	"github.com/gin-gonic/gin"
	"github.com/line/line-bot-sdk-go/v7/linebot"
)

func UseLineBot(cfg *config.Config) gin.HandlerFunc {
	return func(c *gin.Context) {
		bot, err := linebot.New(cfg.Server.ChannelSecret, cfg.Server.ChannelAccessToken)
		if err != nil {
			log.Println(err);
			c.String(500, err.Error())
			c.Abort()
			return
		}
		WithLineBot(c, bot)
	}
}
