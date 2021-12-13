package core

import "gorm.io/gorm"

type DatabaseService interface {
	Session() *gorm.DB
	Migrate() error
}
