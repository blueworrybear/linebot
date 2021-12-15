package api

import (
	"log"

	"github.com/gin-gonic/gin"
	"github.com/line/line-bot-sdk-go/v7/linebot"
)

const (
	keyLineBot = "linebot"
	keyLineEvents = "lineEvents"
)

func WithLineBot(c *gin.Context, linebot *linebot.Client) {
	c.Set(keyLineBot, linebot)
}

func WithEvents(c *gin.Context, events []*linebot.Event) {
	c.Set(keyLineEvents, events)
}

func MustGetLineBotFrom(c *gin.Context) *linebot.Client {
	data := c.MustGet(keyLineBot)
	linebot, ok := data.(*linebot.Client)
	if !ok {
		log.Panic("linebot not found in current context")
	}
	return linebot
}

func MustGetEventsFrom(c *gin.Context) []*linebot.Event {
	data := c.MustGet(keyLineEvents)
	events, ok := data.([]*linebot.Event)
	if !ok {
		log.Panic("events not found in current context")
	}
	return events
}
