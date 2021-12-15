package service

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"sort"
	"time"

	cloudtasks "cloud.google.com/go/cloudtasks/apiv2"
	"github.com/adrg/strutil"
	"github.com/adrg/strutil/metrics"
	"github.com/blueworrybear/lineBot/config"
	"github.com/blueworrybear/lineBot/core"
	"github.com/line/line-bot-sdk-go/v7/linebot"
	taskspb "google.golang.org/genproto/googleapis/cloud/tasks/v2"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type chatService struct {
	userStore core.UserStore
	chatStore core.ChatStore
	cfg *config.Config
}

type relyCandidate struct {
	Chat *core.Chat
	Similarity float64
}

func NewChatService(cfg *config.Config, userStore core.UserStore, chatStore core.ChatStore) core.ChatService {
	return &chatService{
		cfg: cfg,
		userStore: userStore,
		chatStore: chatStore,
	}
}

func (s *chatService) Reply(bot *linebot.Client, event *linebot.Event) error {
	var query string
	switch v := event.Message.(type) {
	case *linebot.TextMessage:
		query = v.Text
	}
	chats, err := s.chatStore.FindWithTag("basic")
	if err != nil {
		return err
	}
	if len(chats) <= 0 {
		return nil
	}
	replies := make([]*relyCandidate, len(chats))
	for i, chat := range chats {
		log.Printf("reply: 43: %v", chat.Message)
		replies[i] = newReplyCandidate(chat, query)
	}
	sort.Slice(replies, func(i, j int) bool {
		return replies[i].Similarity > replies[j].Similarity
	})

	// Select close responses
	ptr := 0
	for _, reply := range replies {
		if replies[0].Similarity - reply.Similarity <= 0.1 {
			ptr++
		}
	}

	source := rand.NewSource(time.Now().Unix())
	r := rand.New(source)
	chat := replies[r.Intn(ptr)].Chat

	if _, err := bot.ReplyMessage(event.ReplyToken, chat.Message).Do(); err != nil {
		return err
	}
	user, err := s.userStore.Find(&core.User{ID: event.Source.UserID})
	if err != nil {
		return err
	}
	for _, nxt := range chat.NextChats {
		if err := s.Push(bot, user, nxt); err != nil {
			return err
		}
	}
	return nil
}

func (s *chatService) PushNow(bot *linebot.Client, user *core.User, chat *core.Chat) error {
	log.Println("PushNow", chat.ID, chat.Message)
	if _, err := bot.PushMessage(user.ID, chat.Message).Do(); err != nil {
		return err
	}
	for _, nxt := range chat.NextChats {
		if err := s.Push(bot, user, nxt); err != nil {
			return err
		}
	}
	return nil
}

func (s *chatService) Push(bot *linebot.Client, user *core.User, chat *core.Chat) error {
	chat, _ = s.chatStore.Find(chat)
	if chat.Delay <= 0 {
		return s.PushNow(bot, user, chat)
	}

	ctx := context.Background()
	client, err := cloudtasks.NewClient(ctx)
	if err != nil {
		return err
	}
	defer client.Close()

	now := time.Now().UTC()
	req := &taskspb.CreateTaskRequest{
		Parent: s.cfg.Queue.Path(),
		Task: &taskspb.Task{
			ScheduleTime: timestamppb.New(now.Add(time.Duration(chat.Delay) * time.Second)),
			MessageType: &taskspb.Task_HttpRequest{
				HttpRequest: &taskspb.HttpRequest{
					HttpMethod: taskspb.HttpMethod_POST,
					Url: fmt.Sprintf("%s/api/chats/push", s.cfg.Server.Host),
				},
			},
		},
	}

	body := &core.PushChatReq{
		ID: chat.ID,
		UserID: user.ID,
	}

	rawBody, err := json.Marshal(body)

	if err != nil {
		return err
	}
	req.Task.GetHttpRequest().Body = rawBody
	if _, err := client.CreateTask(ctx, req); err != nil {
		return err
	}
	return nil
}

func newReplyCandidate(chat *core.Chat, query string) *relyCandidate {
	r := &relyCandidate{
		Chat: chat,
	}
	r.UpdateSimilarity(query)
	return r
}

func (r *relyCandidate) UpdateSimilarity(query string) {
	for _, key := range r.Chat.Keywords {
		if similarity := strutil.Similarity(key, query, metrics.NewHamming()); similarity > r.Similarity {
			r.Similarity = similarity
		}
	}
}
