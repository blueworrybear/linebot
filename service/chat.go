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

type NoChatsError struct {}

func (e *NoChatsError) Error() string {
	return "no chats (chats size 0)"
}

type chatService struct {
	userStore core.UserStore
	chatStore core.ChatStore
	cfg *config.Config
}

type relyCandidate struct {
	Chat *core.Chat
	Similarity float64
}

type bound struct {
	low int64
	up int64
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
	default:
		query = ""
	}

	user, err := s.userStore.Find(&core.User{ID: event.Source.UserID})
	if err != nil {
		return err
	}

	chats, err := s.chatStore.FindWithTag(user.ReplyTag)
	if err != nil {
		return err
	}
	if len(chats) <= 0 {
		return nil
	}

	chat, err := s.SelectWithKeyword(chats, query)
	if err != nil {
		return err
	}

	if _, err := bot.ReplyMessage(event.ReplyToken, chat.Message).Do(); err != nil {
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

func (s *chatService) SelectWithRandom(chats []*core.Chat) (*core.Chat, error) {
	if len(chats) <= 0 {
		return nil, &NoChatsError{}
	}
	now := time.Now()
	bounds := make([]*bound,  len(chats))
	n := int64(0)
	for i, chat := range chats {
		ptr := n + minInt64(int64(now.Sub(chat.LastAccessTime).Seconds()),int64(7 * 24 * 60 * 60))
		bounds[i] = &bound{
			low: n,
			up: ptr,
		}
		n = ptr
	}
	seed := rand.NewSource(now.Unix())
	r := rand.New(seed)
	index := r.Int63n(n)
	for i, bound := range bounds {
		if bound.In(index) {
			return chats[i], nil
		}
	}
	return nil, fmt.Errorf("SelectWithRandom: unknown error")
}

func (s *chatService) SelectWithKeyword(chats []*core.Chat, keyword string) (*core.Chat, error) {
	if len(chats) <= 0 {
		return nil, &NoChatsError{}
	}
	candidates := make([]*relyCandidate, len(chats))
	for i, chat := range chats {
		candidates[i] = newReplyCandidate(chat, keyword)
	}
	sort.Slice(candidates, func(i, j int) bool {
		return candidates[i].Similarity > candidates[j].Similarity
	})
	ptr := 0
	for _, candidate := range candidates {
		if candidates[0].Similarity -  candidate.Similarity < 0.1 {
			ptr ++
		}
	}
	seed := rand.NewSource(time.Now().Unix())
	r := rand.New(seed)
	return candidates[r.Intn(ptr)].Chat, nil
}

func newReplyCandidate(chat *core.Chat, query string) *relyCandidate {
	r := &relyCandidate{
		Chat: chat,
	}
	r.UpdateSimilarity(query)
	return r
}

func (r *relyCandidate) UpdateSimilarity(query string) {
	metric := metrics.NewJaccard()
	metric.NgramSize = 1
	for _, key := range r.Chat.Keywords {
		if similarity := strutil.Similarity(key, query, metric); similarity > r.Similarity {
			r.Similarity = similarity
		}
	}
}

func (b *bound) In(n int64) bool {
	if b.low <= n && b.up > n {
		return true
	}
	return false
}
