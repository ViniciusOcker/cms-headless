package utils

import (
	"cms-headless/internal/models"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func SetupTestDB() *gorm.DB {
	db, _ := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	db.AutoMigrate(&models.Category{})
	db.AutoMigrate(&models.Post{})
	db.AutoMigrate(&models.Project{})
	db.AutoMigrate(&models.Tag{})
	return db
}
