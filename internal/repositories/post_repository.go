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
	FindByID(id uint) (*models.Post, error)
	Create(post *models.Post) error
	Update(post *models.Post) error
	Delete(id uint) error
	SetPostedAt(id uint, t *time.Time) error
	ReplaceTags(post *models.Post, tags []models.Tag) error
	ReplaceCategories(post *models.Post, categories []models.Category) error
	Search(page, pageSize int, categoryID, tagID uint, queryText string, onlyPosted bool) ([]models.Post, int64, error)
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
		now := time.Now().UTC()
		query = query.Where("posted_at IS NOT NULL AND posted_at <= ?", now)
	}

	query.Count(&total)
	err := query.Scopes(utils.PaginateRepository(page, pageSize)).
		Preload("Tags").Preload("Categories").
		Order("posted_at desc").Find(&posts).Error

	return posts, total, err
}

func (r *postRepository) FindBySlug(slug string, onlyPosted bool) (*models.Post, error) {
	var post *models.Post // Começa como nil
	query := r.db.Model(&models.Post{}).Where("slug = ?", slug)

	if onlyPosted {
		query = query.Where("posted_at IS NOT NULL AND posted_at <= ?", time.Now())
	}

	// O GORM preencherá o ponteiro se encontrar o registro
	err := query.Preload("Tags").Preload("Categories").First(&post).Error
	if err != nil {
		return nil, err
	}

	return post, nil
}

func (r *postRepository) FindByID(id uint) (*models.Post, error) {
	var post *models.Post
	err := r.db.Preload("Tags").Preload("Categories").First(&post, id).Error

	if err != nil {
		return nil, err
	}

	return post, nil
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

func (r *postRepository) Search(page, pageSize int, categoryID, tagID uint, queryText string, onlyPosted bool) ([]models.Post, int64, error) {
	var posts []models.Post
	var total int64

	query := r.db.Model(&models.Post{})

	if categoryID > 0 {
		query = query.Joins("JOIN post_categories ON post_categories.post_id = posts.id").
			Where("post_categories.category_id = ?", categoryID)
	}
	if tagID > 0 {
		query = query.Joins("JOIN post_tags ON post_tags.post_id = posts.id").
			Where("post_tags.tag_id = ?", tagID)
	}

	if queryText != "" {
		// Importante: especificar posts.title para evitar ambiguidade com tags.title
		query = query.Where("posts.title LIKE ? OR posts.short_description LIKE ?", "%"+queryText+"%", "%"+queryText+"%")
	}

	if onlyPosted {
		query = query.Where("posts.posted_at IS NOT NULL AND posts.posted_at <= ?", time.Now().UTC())
	}

	// Contagem distinta para não contar o mesmo post múltiplas vezes devido aos joins
	query.Distinct("posts.id").Count(&total)

	// No Find, selecionamos apenas as colunas de posts para o Scan correto
	err := query.Select("posts.*").
		Scopes(utils.PaginateRepository(page, pageSize)).
		Preload("Tags").
		Preload("Categories").
		Order("posts.posted_at desc").
		Find(&posts).Error

	return posts, total, err
}
