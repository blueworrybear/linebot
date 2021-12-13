package api

import (
	"fmt"
	"log"

	"github.com/blueworrybear/lineBot/config"
	"github.com/gin-gonic/gin"
	"github.com/line/line-bot-sdk-go/v7/linebot"
)

type Router struct {
	Config *config.Config
}

func HandleHello() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.String(200, "hi!")
	}
}

func handleFollowEvent(c *gin.Context, event *linebot.Event) error {
	return nil
}

func handleEvents(c *gin.Context, event *linebot.Event) error {
	switch event.Type {
	case linebot.EventTypeFollow:
		return handleFollowEvent(c, event)
	default:
		return fmt.Errorf("unknown event type: %s", event.Type)
	}
}

func HandleWebhooks() gin.HandlerFunc {
	return func(c *gin.Context) {
		linebot := MustGetLineBotFrom(c)
		events, err := linebot.ParseRequest(c.Request)
		if err != nil {
			c.String(500, err.Error())
			return
		}
		log.Println("get webhooks")
		for _, event := range events {
			if err := handleEvents(c, event); err != nil {
				c.String(500, err.Error())
				return
			}
		}
	}
}

func (r *Router) RegisterRoutes(e *gin.Engine) {
	e.GET("/hello", HandleHello())
	e.POST("/webhooks", UseLineBot(r.Config), HandleWebhooks())
}
