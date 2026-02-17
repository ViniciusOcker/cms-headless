package models

import (
	"time"
)

type Tag struct {
	ID        uint   `gorm:"primaryKey;autoIncrement"`
	Title     string `gorm:"uniqueIndex;not null"`
	CreatedAt time.Time
	UpdatedAt time.Time
}
