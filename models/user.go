package models

import (
	"github.com/blueworrybear/lineBot/core"
	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	ID string `gorm:"uniqueIndex;not null"`
}

type UserStore struct {
	DB core.DatabaseService
}
