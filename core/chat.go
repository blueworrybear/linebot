package core

import (
	"time"

	"github.com/line/line-bot-sdk-go/v7/linebot"
)

type Chat struct {
	ID int `json:"id"`
	Delay int `json:"delay"`
	Message linebot.SendingMessage
	Keywords []string
	NextChats []*Chat `json:"nextChats"`
	LastAccessTime time.Time
}

type PushChatReq struct {
	ID int `json:"id"`
	UserID string `json:"user"`
}

type ChatStore interface {
	Create(chat *Chat) error
	CreateWithTag(chat *Chat, tag string) error
	Find(chat *Chat) (*Chat, error)
	FindWithTag(tag string) ([]*Chat, error)
	UpdateLastAccess(chat *Chat) error
}

type ChatService interface {
	Reply(bot *linebot.Client, event *linebot.Event) error
	Push(bot *linebot.Client, user *User, chat *Chat) error
	PushNow(bot *linebot.Client, user *User, chat *Chat) error
	SelectWithRandom(chats []*Chat) (*Chat, error)
	SelectWithKeyword(chats []*Chat, keyword string) (*Chat, error)
}
