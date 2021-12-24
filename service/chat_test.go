package service

import (
	"testing"

	"github.com/adrg/strutil"
	"github.com/adrg/strutil/metrics"
	"github.com/blueworrybear/lineBot/core"
)

func TestSimilarity(t *testing.T) {
	similarity := strutil.Similarity("對", "對", metrics.NewHamming())
	t.Log(similarity)
	similarity = strutil.Similarity("不對", "對", metrics.NewHamming())
	t.Log(similarity)
}

func TestSelectWithKeyword(t *testing.T) {
	s := &chatService{}

	chats := []*core.Chat {
		{
			ID: 0,
			Keywords: []string{
				"快樂",
			},
		},
		{
			ID: 1,
			Keywords: []string{
				"聖誕快樂",
			},
		},
	}

	chat, err := s.SelectWithKeyword(chats, "熊熊，聖誕快樂")
	if err != nil {
		t.Fatal(err)
	}
	if chat.ID != 1 {
		t.Fatal("select wrong ID")
	}

	chats = []*core.Chat {
		{
			ID: 0,
			Keywords: []string{
				"hi",
			},
		},
		{
			ID: 1,
			Keywords: []string{
				"哈囉",
			},
		},
	}

	chat, err = s.SelectWithKeyword(chats, "哈哈")
	if err != nil {
		t.Fatal(err)
	}
	if chat.ID != 1 {
		t.Fatal("select wrong ID")
	}
}
