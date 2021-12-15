package models

import (
	"testing"

	"github.com/blueworrybear/lineBot/core"
	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	"github.com/line/line-bot-sdk-go/v7/linebot"
)

func messageTransformer() cmp.Option {
	return  cmp.Transformer("Chat", func (in linebot.SendingMessage) string {
		switch v := in.(type) {
		case *linebot.TextMessage:
			return v.Text
		}
		return ""
	})
}

func TestChatCreate(t *testing.T) {
	ctrl, db := getDatabaseService(t)
	defer ctrl.Finish()
	store := NewChatStore(db)

	m := &core.Chat{
		Message: linebot.NewTextMessage("test"),
		NextChats: []*core.Chat{
			{
				Message: linebot.NewTextMessage("test2"),
			},
		},
	}
	if err := store.Create(m); err != nil {
		t.Fatal(err)
	}
	if err := store.Create(m); err != nil {
		t.Fatal(err)
	}
	m2 := &core.Chat{
		Message: linebot.NewTextMessage("test3"),
		NextChats: []*core.Chat{m},
	}
	if err := store.Create(m2); err != nil {
		t.Fatal(err)
	}
}

func TestChatFind(t *testing.T) {
	ctrl, db := getDatabaseService(t)
	defer ctrl.Finish()
	store := NewChatStore(db)
	m := &core.Chat{
		Message: linebot.NewTextMessage("test"),
		NextChats: []*core.Chat{
			{
				Message: linebot.NewTextMessage("test2"),
				NextChats: []*core.Chat{},
			},
		},
	}
	if err := store.Create(m); err != nil {
		t.Fatal(err)
	}
	result, err := store.Find(m)
	if err != nil {
		t.Fatal(err)
	}
	if m := cmp.Diff(m, result, cmpopts.IgnoreFields(core.Chat{}, "ID"), messageTransformer()); m != "" {
		t.Fatal(m)
	}
}

func TestFindWithTag(t *testing.T) {
	ctrl, db := getDatabaseService(t)
	defer ctrl.Finish()
	store := NewChatStore(db)

	chats := []*core.Chat {
		{
			Message: linebot.NewTextMessage("tag1"),
		},
		{
			Message: linebot.NewTextMessage("tag2"),
		},
	}
	for _, c := range chats {
		if err := store.CreateWithTag(c, "tag"); err != nil {
			t.Fatal(err)
		}
	}
	result, err := store.FindWithTag("tag")
	if err != nil {
		t.Fatal(err)
	}
	if m := cmp.Diff(chats, result, cmpopts.IgnoreFields(core.Chat{}, "ID", "NextChats"), messageTransformer()); m != "" {
		t.Fatal(m)
	}
}
