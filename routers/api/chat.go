package api

import (
	"log"
	"math/rand"
	"time"

	"github.com/blueworrybear/lineBot/core"
	"github.com/gin-gonic/gin"
	"github.com/line/line-bot-sdk-go/v7/linebot"
)

type textChat struct {
	ID int `json:"id"`
	Tag string `json:"tag"`
	Keywords []string `json:"keywords"`
	Delay int `json:"delay"`
	NextChats []*core.Chat `json:"nextChats"`
	Text string `json:"text"`
	Emojis []*linebot.Emoji `json:"emojis"`
}

type stickerChat struct {
	ID int `json:"id"`
	Tag string `json:"tag"`
	Keywords []string `json:"keywords"`
	Delay int `json:"delay"`
	NextChats []*core.Chat `json:"nextChats"`
	StickerPackage string `json:"package"`
	StickerID string `json:"sticker"`
}

type imageChat struct {
	ID int `json:"id"`
	Tag string `json:"tag"`
	Keywords []string `json:"keywords"`
	Delay int `json:"delay"`
	NextChats []*core.Chat `json:"nextChats"`
	ImageURL string `json:"image"`
}

func (chat *textChat) toChat() *core.Chat {
	message := linebot.NewTextMessage(chat.Text)
	for _, emoji := range chat.Emojis {
		message.AddEmoji(emoji)
	}
	return &core.Chat{
		ID: chat.ID,
		Delay: chat.Delay,
		Keywords: chat.Keywords,
		Message: message,
		NextChats: chat.NextChats,
	}
}

func (chat *stickerChat) toChat() *core.Chat {
	message := linebot.NewStickerMessage(chat.StickerPackage, chat.StickerID)
	return &core.Chat{
		ID: chat.ID,
		Delay: chat.Delay,
		Keywords: chat.Keywords,
		Message: message,
		NextChats: chat.NextChats,
	}
}

func (chat *imageChat) toChat() *core.Chat {
	return &core.Chat{
		ID: chat.ID,
		Delay: chat.Delay,
		Keywords: chat.Keywords,
		Message: linebot.NewImageMessage(chat.ImageURL, chat.ImageURL),
		NextChats: chat.NextChats,
	}
}

func HandlePostTextChat(store core.ChatStore) gin.HandlerFunc {
	return func (c *gin.Context)  {
		chat := &textChat{}
		if err := c.BindJSON(chat); err != nil {
			c.String(400, err.Error())
			return
		}
		for _, nxt := range chat.NextChats {
			if m, err := store.Find(nxt); err != nil || m == nil {
				c.String(400, "next chat not found: %d", nxt.ID)
				return
			}
		}
		m := chat.toChat()
		if err := store.CreateWithTag(m, chat.Tag); err != nil {
			c.String(400, err.Error())
			return
		}
		c.JSON(200, m)
	}
}

func HandlePostStickerChat(store core.ChatStore) gin.HandlerFunc {
	return func (c *gin.Context)  {
		chat := &stickerChat{}
		if err := c.BindJSON(chat); err != nil {
			c.String(400, err.Error())
			return
		}
		for _, nxt := range chat.NextChats {
			if m, err := store.Find(nxt); err != nil || m == nil {
				c.String(400, "next chat not found: %d", nxt.ID)
				return
			}
		}
		m := chat.toChat()
		if err := store.CreateWithTag(m, chat.Tag); err != nil {
			c.String(400, err.Error())
			return
		}
		c.JSON(200, m)
	}
}

func HandlePostImageChat(store core.ChatStore) gin.HandlerFunc {
	return func (c *gin.Context)  {
		chat := &imageChat{}
		if err := c.BindJSON(chat); err != nil {
			c.String(400, err.Error())
			return
		}
		for _, nxt := range chat.NextChats {
			if m, err := store.Find(nxt); err != nil || m == nil {
				c.String(400, "next chat not found: %d", nxt.ID)
				return
			}
		}
		m := chat.toChat()
		if err := store.CreateWithTag(m, chat.Tag); err != nil {
			c.String(400, err.Error())
			return
		}
		c.JSON(200, m)
	}
}

func HandleGetRandomChat(userStore core.UserStore, chatStore core.ChatStore, service core.ChatService) gin.HandlerFunc {
	return func(c *gin.Context) {
		bot := MustGetLineBotFrom(c)
		chats, err := chatStore.FindWithTag(c.Param("tag"))
		if err != nil {
			c.String(400, err.Error())
			return
		}
		if len(chats) <= 0 {
			return
		}
		users, err := userStore.All()
		if err != nil {
			c.String(400, err.Error())
			return
		}
		s := rand.NewSource(time.Now().Unix())
		r := rand.New(s)
		chat := chats[r.Intn(len(chats))]
		for _, u := range users {
			if err := service.Push(bot, u, chat); err != nil {
				log.Println("HandleGetRandomChat", err)
			}
		}
	}
}

func HandlePostChatPush(userStore core.UserStore, chatStore core.ChatStore, chatService core.ChatService) gin.HandlerFunc {
	return func(c *gin.Context) {
		req := &core.PushChatReq{}
		if err := c.BindJSON(req); err != nil {
			c.String(400, err.Error())
			return
		}
		user, err := userStore.Find(&core.User{ID: req.UserID})
		if err != nil {
			c.String(400, err.Error())
			return
		}
		chat, err := chatStore.Find(&core.Chat{ID: req.ID})
		if err != nil {
			c.String(400, err.Error())
			return
		}
		bot := MustGetLineBotFrom(c)
		if err := chatService.PushNow(bot, user, chat); err != nil {
			c.String(400, err.Error())
			return
		}
	}
}
