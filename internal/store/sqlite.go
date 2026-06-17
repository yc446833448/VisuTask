package store

import (
	"log"

	"github.com/yc446833448/VisuTask/internal/model"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

type DB struct {
	*gorm.DB
}

func New(dbPath string) (*DB, error) {
	db, err := gorm.Open(sqlite.Open(dbPath), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Warn),
	})
	if err != nil {
		return nil, err
	}

	// Auto-migrate all models
	if err := db.AutoMigrate(
		&model.Script{},
		&model.Task{},
		&model.Execution{},
		&model.User{},
	); err != nil {
		return nil, err
	}

	log.Printf("database initialized at %s", dbPath)
	return &DB{db}, nil
}
