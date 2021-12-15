package models

import (
	"github.com/blueworrybear/lineBot/core"
	"gorm.io/gorm"
)

var (
	tables []interface{}
)
type databaseService struct {
	db *gorm.DB
}

func init() {
	tables = append(tables, &user{}, &chat{}, &question{})
}

func NewDatabaseService(db *gorm.DB) core.DatabaseService {
	return &databaseService{
		db: db,
	}
}

func (store *databaseService) Session() *gorm.DB {
	return store.db.Session(&gorm.Session{})
}

func (store *databaseService) Migrate() error {
	return migrate(store.db)
}

func migrate(db *gorm.DB) error {
	return db.AutoMigrate(tables...)
}
