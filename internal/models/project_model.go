package models

import (
	"time"

	"gorm.io/gorm"
)

type Project struct {
	ID               uint   `gorm:"primaryKey;autoIncrement"`
	Title            string `gorm:"uniqueIndex;not null"`
	Slug             string `gorm:"uniqueIndex;not null"`
	ShortDescription string
	Body             string `gorm:"type:text"`
	DemoURL          string
	RepoURL          string
	CreatedAt        time.Time      // Padronizado para CreatedAt
	UpdatedAt        time.Time      // Padronizado para UpdatedAt
	PostedAt         *time.Time     `gorm:"index"`
	DeletedAt        gorm.DeletedAt `gorm:"index"` // Alterado para Soft Delete do GORM

	// Relacionamentos
	Tags       []Tag      `gorm:"many2many:project_tags;"`
	Categories []Category `gorm:"many2many:project_categories;"`
}
