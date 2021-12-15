package models

import (
	"github.com/blueworrybear/lineBot/core"
	"github.com/lib/pq"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type question struct {
	gorm.Model
	ID int `gorm:"uniqueIndex;not null"`
	ChatID int
	Chat *chat
	Answer string
	Options pq.StringArray `gorm:"type:text[]"`
	OkResponse []*chat `gorm:"many2many:ok_response"`
	ErrorResponse []*chat `gorm:"many2many:error_response"`
	NextQuestions []*question `gorm:"many2many:next_questions"`
}

type questionStore struct {
	db core.DatabaseService
}

func NewQuestionStore(db core.DatabaseService) core.QuestionStore {
	return &questionStore{db:db}
}

func (s *questionStore) Create(q *core.Question) error {
	x := s.db.Session()
	m, err := toQuestion(q)
	if err != nil {
		return err
	}
	err = x.Clauses(clause.OnConflict{DoNothing: true}).Create(m).Error
	if err != nil {
		return err
	}
	q.ID = m.ID
	return nil
}

func (s *questionStore) Find(q *core.Question) (*core.Question, error) {
	x := s.db.Session()
	target := &question{}
	if err :=	x.Preload(
		"Chat",
	).Preload(
		"OkResponse",
	).Preload(
		"ErrorResponse",
	).Preload(
		"NextQuestions",
	).Find(target, &question{ID: q.ID}).Error; err != nil {
		return nil, err
	}
	return target.toCore(), nil
}

func toQuestion(c *core.Question) (*question, error){
	mChat, err := toChat(c.Chat, "")
	if err != nil {
		return nil, err
	}

	okResponse, err := sliceToChat(c.OkResponse, "")
	if err != nil {
		return nil, err
	}
	errorResponse, err := sliceToChat(c.ErrorResponse, "")
	if err != nil {
		return nil, err
	}
	nextQuestions := []*question{}
	if (c.NextQuestions != nil) {
		nextQuestions = make([]*question, len(c.NextQuestions))
		for i, n := range c.NextQuestions {
			nextQuestions[i], err = toQuestion(n)
			if err != nil {
				return nil, err
			}
		}
	}
	return &question{
		ID: c.ID,
		Chat: mChat,
		Answer: c.Answer,
		Options: c.Options,
		OkResponse: okResponse,
		ErrorResponse: errorResponse,
		NextQuestions: nextQuestions,
	}, nil
}

func (m *question) toCore() *core.Question {
	nextQuestions := make([]*core.Question, len(m.NextQuestions))
	for i, n := range m.NextQuestions {
		nextQuestions[i] = n.toCore()
	}
	var coreChat *core.Chat
	if m.Chat != nil {
		coreChat = m.Chat.toCore()
	}
	return &core.Question{
		ID: m.ID,
		Chat: coreChat,
		Answer: m.Answer,
		Options: m.Options,
		OkResponse: chatSlice(m.OkResponse).toCore(),
		ErrorResponse: chatSlice(m.ErrorResponse).toCore(),
		NextQuestions: nextQuestions,
	}
}
