package repositories

import (
	"cms-headless/internal/models"
	"cms-headless/internal/utils"
	"time"

	"gorm.io/gorm"
)

type PostRepository interface {
	FindAll(page, pageSize int, onlyPosted bool) ([]models.Post, int64, error)
	FindBySlug(slug string, onlyPosted bool) (*models.Post, error)
	// FindByID(id uint) (*models.Post, error)
	Create(post *models.Post) error
	Update(post *models.Post) error
	Delete(id uint) error
	SetPostedAt(id uint, t *time.Time) error
	ReplaceTags(post *models.Post, tags []models.Tag) error
	ReplaceCategories(post *models.Post, categories []models.Category) error
}

type postRepository struct {
	db *gorm.DB
}

func NewPostRepository(db *gorm.DB) PostRepository {
	return &postRepository{db: db}
}

func (r *postRepository) FindAll(page, pageSize int, onlyPosted bool) ([]models.Post, int64, error) {
	var posts []models.Post
	var total int64

	query := r.db.Model(&models.Post{})
	if onlyPosted {
		query = query.Where("posted_at IS NOT NULL AND posted_at <= ?", time.Now())
	}

	query.Count(&total)
	err := query.Scopes(utils.PaginateRepository(page, pageSize)).
		Preload("Tags").Preload("Categories").
		Order("posted_at desc").Find(&posts).Error

	return posts, total, err
}

func (r *postRepository) FindBySlug(slug string, onlyPosted bool) (*models.Post, error) {
	var post models.Post
	query := r.db.Where("slug = ?", slug)
	if onlyPosted {
		query = query.Where("posted_at IS NOT NULL AND posted_at <= ?", time.Now())
	}
	err := query.Preload("Tags").Preload("Categories").First(&post).Error
	return &post, err
}

func (r *postRepository) Create(post *models.Post) error {
	return r.db.Create(post).Error
}

func (r *postRepository) Update(post *models.Post) error {
	return r.db.Save(post).Error
}

func (r *postRepository) Delete(id uint) error {
	return r.db.Delete(&models.Post{}, id).Error
}

func (r *postRepository) SetPostedAt(id uint, t *time.Time) error {
	return r.db.Model(&models.Post{}).Where("id = ?", id).Update("posted_at", t).Error
}

func (r *postRepository) ReplaceTags(post *models.Post, tags []models.Tag) error {
	return r.db.Model(post).Association("Tags").Replace(tags)
}

func (r *postRepository) ReplaceCategories(post *models.Post, categories []models.Category) error {
	return r.db.Model(post).Association("Categories").Replace(categories)
}
