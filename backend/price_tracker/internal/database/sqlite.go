package database

import (
	"price-tracker/internal/models"

	"github.com/glebarez/sqlite"
	"gorm.io/gorm"
)

var DB *gorm.DB

func InitDB() {
	var err error
	DB, err = gorm.Open(sqlite.Open("tracker.db"), &gorm.Config{})
	if err != nil {
		panic(err)
	}
	db, err := DB.DB()
	if err == nil {
		db.Exec("PRAGMA journal_mode=WAL;")
		db.Exec("PRAGMA synchronous=NORMAL;")
	}
	DB.AutoMigrate(&models.User{}, &models.Product{})
}
