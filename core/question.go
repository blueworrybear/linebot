package core

import "github.com/line/line-bot-sdk-go/v7/linebot"

type QuestionCheckMethod string

const (
	QuestionCheckMethodEqual QuestionCheckMethod = "equal"
)

type Question struct {
	ID int `json:"id"`
	Chat *Chat `json:"chat"`
	Answer string `json:"answer"`
	Options []string `json:"options"`
	OkResponse []*Chat `json:"ok"`
	ErrorResponse []*Chat `json:"error"`
	NextQuestions []*Question `json:"next"`
}

type QuestionStore interface {
	Create(q *Question) error
	Find(q *Question) (*Question, error)
}

type QuestionService interface {
	Ask(bot *linebot.Client, u *User, q *Question) error
	Answer(bot *linebot.Client, u *User, event *linebot.Event) error
}
