package models

import (
	"time"

	"gorm.io/gorm"
)

type Post struct {
	ID               uint   `gorm:"primaryKey;autoIncrement"`
	Title            string `gorm:"uniqueIndex;not null"`
	Slug             string `gorm:"uniqueIndex;not null"`
	ShortDescription string
	Body             string `gorm:"type:text"`
	CreatedAt        time.Time
	UpdatedAt        time.Time
	PostedAt         *time.Time     `gorm:"index"`
	DeletedAt        gorm.DeletedAt `gorm:"index"`

	// Relacionamentos
	Tags       []Tag      `gorm:"many2many:post_tags;"`
	Categories []Category `gorm:"many2many:post_categories;"`
}
