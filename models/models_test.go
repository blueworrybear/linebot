package models

import (
	"log"
	"os"
	"testing"

	"github.com/blueworrybear/lineBot/core"
	"github.com/blueworrybear/lineBot/mock"
	"github.com/golang/mock/gomock"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var db *gorm.DB

const (
	dbName = "test.db"
)

func OpenDatabase() *gorm.DB {
	db, err := gorm.Open(sqlite.Open(dbName), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent),
	})
	if err != nil {
		log.Fatal(err)
	}
	return db
}

func getDatabaseService(t *testing.T) (*gomock.Controller, core.DatabaseService) {
	ctrl := gomock.NewController(t)
	mockService := mock.NewMockDatabaseService(ctrl)
	mockService.EXPECT().Session().AnyTimes().Return(db.Session(&gorm.Session{}))
	return ctrl, mockService
}

func TestMain(m *testing.M) {
	if _, err := os.Stat(dbName); err == nil {
		os.Remove(dbName)
	}
	db = OpenDatabase()
	if err := migrate(db); err != nil {
		log.Fatal(err)
	}
	exit := m.Run()
	defer os.Exit(exit)
}
