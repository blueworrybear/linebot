package service

import (
	"fmt"
	"math/rand"
	"time"

	"github.com/adrg/strutil"
	"github.com/adrg/strutil/metrics"
	"github.com/blueworrybear/lineBot/core"
	"github.com/line/line-bot-sdk-go/v7/linebot"
)

type questionService struct {
	userStore core.UserStore
	questionStore core.QuestionStore
	chatService core.ChatService
}

func NewQuestionService(userStore core.UserStore, questionStore core.QuestionStore, chatService core.ChatService) core.QuestionService {
	return &questionService{
		userStore: userStore,
		chatService: chatService,
		questionStore: questionStore,
	}
}

func (s *questionService) Ask(bot *linebot.Client, u *core.User, q *core.Question) error {
	q, _ = s.questionStore.Find(q)
	if err := s.chatService.Push(bot, u, q.Chat); err != nil {
		return err
	}
	u.State = core.UserStateWaiting
	u.Question = q
	return s.userStore.Update(u)
}

func (s *questionService) Answer(bot *linebot.Client, u *core.User, event *linebot.Event) error {
	if u.State != core.UserStateWaiting {
		return fmt.Errorf("wrong user state")
	}
	if u.Question == nil {
		return fmt.Errorf("no question")
	}
	message, ok := event.Message.(*linebot.TextMessage)
	if !ok {
		if chat := s.selectResponse(u.Question.ErrorResponse); chat != nil {
			return s.chatService.Push(bot, u, chat)
		} else {
			_, err := bot.PushMessage(u.ID, linebot.NewTextMessage("NO!")).Do()
			return err
		}
	}
	similarity := strutil.Similarity(u.Question.Answer, message.Text, metrics.NewHamming())
	if similarity > 0.8 {
		if chat := s.selectResponse(u.Question.OkResponse); chat != nil {
			if err := s.chatService.Push(bot, u, chat); err != nil {
				return err
			}
		} else {
			if _, err := bot.PushMessage(u.ID, linebot.NewTextMessage("OK!")).Do(); err != nil {
				return err
			}
		}
		if u.Question.NextQuestions == nil || len(u.Question.NextQuestions) <= 0 {
			u.State = core.UserStateIdle
			return s.userStore.Update(u)
		}
		return s.Ask(bot, u, u.Question.NextQuestions[0])
	}
	if chat := s.selectResponse(u.Question.ErrorResponse); chat != nil {
		return s.chatService.Push(bot, u, chat)
	}
	_, err := bot.PushMessage(u.ID, linebot.NewTextMessage("NO!")).Do()
	return err
}

func (s *questionService) selectResponse(chats []*core.Chat) *core.Chat {
	if len(chats) <= 0 {
		return nil
	}
	seed := rand.NewSource(time.Now().Unix())
	r := rand.New(seed)
	return chats[r.Intn(len(chats))]
}
