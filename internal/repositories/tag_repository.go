package repositories

import (
	"cms-headless/internal/models"
	"cms-headless/internal/utils"

	"gorm.io/gorm"
)

type TagRepository interface {
	FindAll(page, pageSize int) ([]models.Tag, int64, error)
	Create(tag *models.Tag) error
	UpdateName(id uint, newTitle string) error
}

type tagRepository struct {
	db *gorm.DB
}

func NewTagRepository(db *gorm.DB) TagRepository {
	return &tagRepository{db: db}
}

// Listagem de tags
func (r *tagRepository) FindAll(page, pageSize int) ([]models.Tag, int64, error) {
	var tags []models.Tag
	var total int64

	r.db.Model(&models.Tag{}).Count(&total)
	err := r.db.Scopes(utils.PaginateRepository(page, pageSize)).Order("title asc").Find(&tags).Error

	return tags, total, err
}

// Criar uma nova tag
func (r *tagRepository) Create(tag *models.Tag) error {
	return r.db.Create(tag).Error
}

// Renomear uma tag
func (r *tagRepository) UpdateName(id uint, newTitle string) error {
	// Usamos o Model para garantir que o GORM saiba qual tabela e ID usar
	return r.db.Model(&models.Tag{}).Where("id = ?", id).Update("title", newTitle).Error
}
