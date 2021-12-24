package api

import (
	"log"

	"github.com/blueworrybear/lineBot/config"
	"github.com/blueworrybear/lineBot/core"
	"github.com/gin-gonic/gin"
	"github.com/line/line-bot-sdk-go/v7/linebot"
)

type Router struct {
	Config *config.Config
	UserStore core.UserStore
	ChatStore core.ChatStore
	QuestionStore core.QuestionStore
	ChatService core.ChatService
	QuestionService core.QuestionService
}

func HandleIndex() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.String(200, "ok")
	}
}

func HandleFollowEvent(userStore core.UserStore) gin.HandlerFunc {
	return func(c *gin.Context) {
		bot := MustGetLineBotFrom(c)
		events := MustGetEventsFrom(c)
		for _, event := range events {
			if event.Type != linebot.EventTypeFollow {
				continue
			}
			if event.Source.Type != linebot.EventSourceTypeUser || event.Source.UserID == ""{
				continue
			}
			user := &core.User{ID: event.Source.UserID}
			if u, err := userStore.Find(user); err == nil && u != nil {
				continue
			}
			if profile, err := bot.GetProfile(event.Source.UserID).Do(); err == nil {
				user.Name = profile.DisplayName
			}
			if err := userStore.Create(user); err != nil {
				log.Println(err)
				continue
			}
			if _, err := bot.ReplyMessage(event.ReplyToken, linebot.NewStickerMessage("11537", "52002734")).Do(); err != nil {
				log.Println(err)
				return
			}
		}
	}
}

func HandleMessageEvent(usrStore core.UserStore, chatService core.ChatService, questionService core.QuestionService) gin.HandlerFunc {
	return func(c *gin.Context) {
		bot := MustGetLineBotFrom(c)
		events := MustGetEventsFrom(c)
		for _, event := range events {
			if event.Type != linebot.EventTypeMessage {
				continue
			}
			_, err := usrStore.FindForRequest(&core.User{ID: event.Source.UserID})
			if err != nil {
				log.Println(err)
				continue
			}
			user, err := usrStore.Find(&core.User{ID: event.Source.UserID})
			if err != nil {
				log.Println(err)
				continue
			}
			if user.Role == core.UserRoleInactive {
				continue
			}
			switch user.State {
			case core.UserStateIdle:
				if err := chatService.Reply(bot, event); err != nil {
					log.Println(err)
				}
			case core.UserStateWaiting:
				if err := questionService.Answer(bot, user, event); err != nil {
					log.Println(err)
				}
			}
		}
	}
}

func HandleWebhooks() gin.HandlerFunc {
	return func(c *gin.Context) {
		defer c.Abort()
		c.String(200, "get webhooks")
	}
}

func (r *Router) RegisterRoutes(e *gin.Engine) {
	e.GET("", HandleIndex())
	{
		g := e.Group("/webhooks", UseLineBot(r.Config), ParseRequest())
		g.POST(
			"",
			HandleFollowEvent(r.UserStore),
			HandleMessageEvent(r.UserStore, r.ChatService, r.QuestionService),
			HandleWebhooks(),
		)
	}
	{
		g := e.Group("/api")
		g.POST("/chats/text", HandlePostTextChat(r.ChatStore))
		g.POST("/chats/sticker", HandlePostStickerChat(r.ChatStore))
		g.POST("/chats/image", HandlePostImageChat(r.ChatStore))
		g.GET(
			"/chats/random/:tag",
			UseLineBot(r.Config),
			HandleGetRandomChat(r.UserStore, r.ChatStore, r.ChatService),
		)
		g.POST(
			"/chats/push", UseLineBot(r.Config),
			HandlePostChatPush(r.UserStore, r.ChatStore, r.ChatService),
		)
		g.POST("/questions", HandlePostQuestions(r.QuestionStore))
		g.GET("/questions/:id",
			UseLineBot(r.Config),
			HandleGetQuestionsAsk(r.UserStore, r.QuestionStore, r.QuestionService),
		)
	}
}
