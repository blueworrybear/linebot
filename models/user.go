package models

import (
	"fmt"
	"time"

	"github.com/blueworrybear/lineBot/core"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type user struct {
	gorm.Model
	ID string `gorm:"uniqueIndex;not null"`
	Name string
	State int `gorm:"index"`
	Questions []*question `gorm:"many2many:user_questions"`
	ReplyTag string `gorm:"default:'basic'"`
	Role string `gorm:"index;default:'normal'"`
	LastRequestTime time.Time
}

type userStore struct {
	db core.DatabaseService
}

func NewUserStore(db core.DatabaseService) core.UserStore {
	return &userStore{db:db}
}

func (store *userStore) Create(c *core.User) error {
	x := store.db.Session()
	u := &user{
		ID: c.ID,
		Name: c.Name,
		State: int(c.State),
	}
	return x.Create(u).Error
}

func (store *userStore) Find(c *core.User) (*core.User, error) {
	x := store.db.Session()
	u := &user{}
	if err := x.Preload(
		"Questions",
	).Preload(
		"Questions.Chat",
	).Preload(
		"Questions.OkResponse",
	).Preload(
		"Questions.ErrorResponse",
	).Preload(
		"Questions.NextQuestions",
	).Where(&user{ID: c.ID}).First(u).Error; err != nil {
		return nil, err
	}
	return u.toCore(), nil
}

func (store *userStore) FindForRequest(c *core.User) (*core.User, error) {
	x := store.db.Session()
	tx := x.Begin()
	defer tx.Rollback()

	u := &user{}

	if err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).Where(&user{ID: c.ID}).First(u).Error; err != nil {
		return nil, err
	}

	now := time.Now()

	if now.Sub(u.LastRequestTime).Seconds() < core.USER_REQUEST_THROTTLE_SECONDS {
		return nil, fmt.Errorf("too fast")
	}
	u.LastRequestTime = now
	if err := tx.Model(&u).Updates(u).Error; err != nil {
		return nil, err
	}
	return u.toCore(), tx.Commit().Error
}

func (store *userStore) All() ([]*core.User, error) {
	x := store.db.Session()
	var users []*user
	if err := x.Find(&users).Error; err != nil {
		return nil, err
	}
	results := make([]*core.User, len(users))
	for i, u := range users {
		results[i] = u.toCore()
	}
	return results, nil
}

func (store *userStore) VIPs() ([]*core.User, error) {
	return store.All()
}

func (store *userStore) Update(c *core.User) error {
	x := store.db.Session()
	err := x.Model(&user{ID: c.ID}).Update("state", int(c.State)).Error
	if err != nil {
		return err
	}
	if c.Question != nil {
		q, err := toQuestion(c.Question)
		if err != nil {
			return err
		}
		return x.Model(&user{ID: c.ID}).Association("Questions").Replace([]*question{q})
	}
	return x.Model(&user{ID: c.ID}).Association("Questions").Clear()
}

func (store *userStore) SetRequestTime(c *core.User, t time.Time) error {
	x := store.db.Session()
	c.LastRequestTime = t
	return x.Model(&user{ID: c.ID}).Updates(&user{LastRequestTime: t}).Error
}

func (m *user) toCore() *core.User {
	var q *core.Question
	if m.Questions != nil && len(m.Questions) > 0 {
		q = m.Questions[0].toCore()
	}
	return &core.User{
		ID: m.ID,
		Name: m.Name,
		State: core.UserState(m.State),
		Question: q,
		Role: core.UserRole(m.Role),
		ReplyTag: m.ReplyTag,
		LastRequestTime: m.LastRequestTime,
	}
}
