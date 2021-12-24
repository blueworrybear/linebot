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

func checkSelect(t *testing.T, s *chatService, chats []string, keyword string, expect int) {
	arr := make([]*core.Chat, len(chats))
	for i, chat := range chats {
		arr[i] = &core.Chat{
			ID: i,
			Keywords: []string{chat},
		}
	}
	ans, err := s.SelectWithKeyword(arr, keyword)
	if err != nil {
		t.Fatal(err)
	}
	if ans.ID != expect {
		t.Fatal("select wrong ID", chats, keyword)
	}
}

func TestSelectWithKeyword(t *testing.T) {
	s := &chatService{}

	checkSelect(t, s, []string{"快樂",	"聖誕快樂"}, "熊熊，聖誕快樂", 1)
	checkSelect(t, s, []string{"hi",	"哈囉"}, "哈哈", 1)
	checkSelect(t, s, []string{"",	"快樂", "愛"}, "愛", 2)
}
