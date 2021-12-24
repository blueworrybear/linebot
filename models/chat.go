package models

import (
	"encoding/json"
	"log"
	"time"

	"github.com/blueworrybear/lineBot/core"
	"github.com/lib/pq"
	"github.com/line/line-bot-sdk-go/v7/linebot"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type chat struct {
	gorm.Model
	ID int `gorm:"uniqueIndex;not null"`
	MessageType string
	// TextMessage
	Text string
	Emojis []byte
	// StickerMessage
	StickerPackage string
	StickerID string
	// ImageMessage
	ImageURL string
	// Meta
	Delay int
	Tag string `gorm:"index;size:256"`
	Keywords pq.StringArray `gorm:"type:text[]"`
	NextChats []*chat `gorm:"many2many:next_chats"`
	LastAccessTime time.Time
}

type chatSlice []*chat

type chatStore struct {
	db core.DatabaseService
}

func NewChatStore(db core.DatabaseService) core.ChatStore {
	return &chatStore{db: db}
}

func (store *chatStore) Create(c *core.Chat) error {
	return store.createWithTag(c, "")
}

func (store *chatStore) CreateWithTag(c *core.Chat, tag string) error {
	return store.createWithTag(c, tag)
}

func (store *chatStore) createWithTag(c *core.Chat, tag string) error {
	x := store.db.Session()
	m, err := toChat(c, tag)
	if err != nil {
		return err
	}
	err = x.Clauses(clause.OnConflict{DoNothing: true}).Create(m).Error
	c.ID = m.ID
	return err
}

func (store *chatStore) Find(c *core.Chat) (*core.Chat, error) {
	x := store.db.Session()
	u := &chat{}
	if err := x.Preload("NextChats").First(u, &chat{ID: c.ID}).Error; err != nil {
		return nil, err
	}
	return u.toCore(), nil
}

func (store *chatStore) FindWithTag(tag string) ([]*core.Chat, error) {
	var results []*chat
	x := store.db.Session()
	if err := x.Preload("NextChats").Find(&results, &chat{Tag: tag}).Error; err != nil {
		return nil, err
	}
	coreResults := make([]*core.Chat, len(results))
	for i, result := range results {
		coreResults[i] = result.toCore()
	}
	return coreResults, nil
}

func (store *chatStore) UpdateLastAccess(c *core.Chat) error {
	x := store.db.Session()
	return x.Model(&chat{ID: c.ID}).Updates(&chat{LastAccessTime: time.Now()}).Error
}

func toChat(c *core.Chat, tag string) (*chat, error) {
	if c == nil {
		return nil, nil
	}
	var err error
	nextChats := []*chat{}
	if c.NextChats != nil {
		nextChats = make([]*chat, len(c.NextChats))
		for i, nxt := range c.NextChats {
			if nextChats[i], err = toChat(nxt, tag); err != nil {
				return nil, err
			}
		}
	}

	target := &chat{
		ID: c.ID,
		Delay: c.Delay,
		NextChats: nextChats,
		Keywords: c.Keywords,
		Tag: tag,
	}

	switch v := c.Message.(type) {
	case *linebot.TextMessage:
		target.Text = v.Text
		target.MessageType = string(v.Type())
		if v.Emojis != nil {
			emojis, _ := json.Marshal(v.Emojis)
			target.Emojis = emojis
		}
	case *linebot.StickerMessage:
		target.StickerPackage = v.PackageID
		target.StickerID = v.StickerID
		target.MessageType = string(v.Type())
	case *linebot.ImageMessage:
		target.ImageURL = v.OriginalContentURL
		target.MessageType = string(v.Type())
	default:
		log.Println("unknown message type")
	}

	return target, nil
}

func sliceToChat(slice []*core.Chat, tag string) ([]*chat, error) {
	if slice == nil {
		return []*chat{}, nil
	}
	var err error
	target := make([]*chat, len(slice))
	for i, n := range slice {
		target[i], err = toChat(n, tag)
		if err != nil {
			return nil, err
		}
	}
	return target, nil
}

func (m *chat) toCore() *core.Chat {
	var message linebot.SendingMessage
	switch linebot.MessageType(m.MessageType) {
	case linebot.MessageTypeText:
		message = linebot.NewTextMessage(m.Text)
		if m.Emojis != nil && len(m.Emojis) > 0 {
			var emojis []*linebot.Emoji
			if err := json.Unmarshal(m.Emojis, &emojis); err != nil {
				log.Fatal(err)
			}
			for _, emoji := range emojis {
				message.AddEmoji(emoji)
			}
		}
	case linebot.MessageTypeSticker:
		message = linebot.NewStickerMessage(m.StickerPackage, m.StickerID)
	case linebot.MessageTypeImage:
		message = linebot.NewImageMessage(m.ImageURL, m.ImageURL)
	default:
		log.Fatalf("unknown message %s", m.MessageType)
	}
	nextChats := make([]*core.Chat, len(m.NextChats))
	for i, nxt := range m.NextChats {
		nextChats[i] = nxt.toCore()
	}
	return &core.Chat{
		ID: m.ID,
		Delay: m.Delay,
		Message: message,
		Keywords: m.Keywords,
		NextChats: nextChats,
		LastAccessTime: m.LastAccessTime,
	}
}

func (m chatSlice) toCore() []*core.Chat {
	target := make([]*core.Chat, len(m))
	for i, n := range m {
		target[i] = n.toCore()
	}
	return target
}
