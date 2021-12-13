package api

import (
	"log"

	"github.com/gin-gonic/gin"
	"github.com/line/line-bot-sdk-go/v7/linebot"
)

const (
	keyLineBot = "linebot"
)

func WithLineBot(c *gin.Context, linebot *linebot.Client) {
	c.Set(keyLineBot, linebot)
}

func MustGetLineBotFrom(c *gin.Context) *linebot.Client {
	data := c.MustGet(keyLineBot)
	linebot, ok := data.(*linebot.Client)
	if !ok {
		log.Panic("linebot not found in current context")
	}
	return linebot
}
