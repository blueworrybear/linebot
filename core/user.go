package core

import "time"

//go:generate mockgen -package mock -destination ../mock/user_mock.go . UserStore

type UserState int
type UserRole string

const (
	USER_REQUEST_THROTTLE_SECONDS = 5
)

const (
	UserStateIdle UserState = iota
	UserStateWaiting
	UserStateListening
)

const (
	UserRoleAdmin UserRole = "admin"
	UserRoleVIP UserRole = "vip"
	UserRoleNormal UserRole = "normal"
	UserRoleInactive UserRole = "inactive"
)

type User struct {
	ID string
	Name string
	State UserState
	Question *Question
	Role UserRole
	ReplyTag string
	LastRequestTime time.Time
}

type UserStore interface{
	Create(user *User) error
	Find(user *User) (*User, error)
	FindForRequest(user *User) (*User, error)
	Update(user *User) error
	All() ([]*User, error)
	VIPs() ([]*User, error)
}
