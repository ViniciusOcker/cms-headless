package utils

import (
	"log"
	"os"

	"gorm.io/driver/postgres"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func InitDB() *gorm.DB {
	var dialector gorm.Dialector
	env := os.Getenv("APP_ENV")

	if env == "prod" {
		dsn := os.Getenv("DATABASE_URL")
		dialector = postgres.Open(dsn)
	} else {
		dbPath := "cms_dev.db"
		if env == "test" {
			dbPath = "file::memory:?cache=shared" // SQLite em memória para testes ultra-rápidos
		}
		dialector = sqlite.Open(dbPath)
	}

	db, err := gorm.Open(dialector, &gorm.Config{
		DisableForeignKeyConstraintWhenMigrating: false,
	})

	if err != nil {
		log.Fatalf("Falha ao conectar no banco (%s): %v", env, err)
	}

	return db
}
