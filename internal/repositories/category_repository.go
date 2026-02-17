package repositories

import (
	"cms-headless/internal/models"
	"cms-headless/internal/utils"

	"gorm.io/gorm"
)

type CategoryRepository interface {
	FindAll(page, pageSize int) ([]models.Category, int64, error)
	Create(category *models.Category) error
	UpdateName(id uint, newTitle string) error
}

type categoryRepository struct {
	db *gorm.DB
}

func NewCategoryRepository(db *gorm.DB) CategoryRepository {
	return &categoryRepository{db: db}
}

func (r *categoryRepository) FindAll(page, pageSize int) ([]models.Category, int64, error) {
	var categories []models.Category
	var total int64

	r.db.Model(&models.Category{}).Count(&total)
	err := r.db.Scopes(utils.PaginateRepository(page, pageSize)).Order("title asc").Find(&categories).Error

	return categories, total, err
}

func (r *categoryRepository) Create(category *models.Category) error {
	return r.db.Create(category).Error
}

func (r *categoryRepository) UpdateName(id uint, newTitle string) error {
	return r.db.Model(&models.Category{}).Where("id = ?", id).Update("title", newTitle).Error
}
