package repositories

import (
	"cms-headless/internal/models"
	"cms-headless/internal/utils"
	"time"

	"gorm.io/gorm"
)

type ProjectRepository interface {
	FindAll(page, pageSize int, onlyPosted bool) ([]models.Project, int64, error)
	FindBySlug(slug string, onlyPosted bool) (*models.Project, error)
	FindByID(id uint) (*models.Project, error)
	Create(project *models.Project) error
	Update(project *models.Project) error
	Delete(id uint) error
	SetPostedAt(id uint, t *time.Time) error
	Filter(page, pageSize int, categoryID uint, tagID uint, onlyPosted bool) ([]models.Project, int64, error)
	ReplaceTags(project *models.Project, tags []models.Tag) error
	ReplaceCategories(project *models.Project, categories []models.Category) error
}

type projectRepository struct {
	db *gorm.DB
}

func NewProjectRepository(db *gorm.DB) ProjectRepository {
	return &projectRepository{db: db}
}

func (r *projectRepository) FindAll(page, pageSize int, onlyPosted bool) ([]models.Project, int64, error) {
	var projects []models.Project
	var total int64

	query := r.db.Model(&models.Project{})
	if onlyPosted {
		query = query.Where("posted_at IS NOT NULL AND posted_at <= ?", time.Now())
	}

	query.Count(&total)

	// Usando o utils renomeado
	err := query.Scopes(utils.PaginateRepository(page, pageSize)).
		Preload("Tags").
		Preload("Categories").
		Order("created_at desc").
		Find(&projects).Error

	return projects, total, err
}

func (r *projectRepository) FindBySlug(slug string, onlyPosted bool) (*models.Project, error) {
	var project models.Project // Use a struct, não o ponteiro diretamente aqui
	query := r.db.Where("slug = ?", slug)

	if onlyPosted {
		// Truncamos para segundos para garantir compatibilidade total com SQLite/Postgres
		now := time.Now().UTC()
		query = query.Where("posted_at IS NOT NULL AND posted_at <= ?", now)
	}

	err := query.Preload("Tags").Preload("Categories").First(&project).Error

	if err != nil {
		return nil, err
	}

	return &project, nil
}

func (r *projectRepository) FindByID(id uint) (*models.Project, error) {
	var project models.Project
	err := r.db.Preload("Tags").Preload("Categories").First(&project, id).Error
	if err != nil {
		return nil, err
	}

	return &project, nil
}

func (r *projectRepository) Update(project *models.Project) error {
	return r.db.Save(project).Error
}

func (r *projectRepository) ReplaceTags(project *models.Project, tags []models.Tag) error {
	return r.db.Model(project).Association("Tags").Replace(tags)
}

func (r *projectRepository) ReplaceCategories(project *models.Project, categories []models.Category) error {
	return r.db.Model(project).Association("Categories").Replace(categories)
}

func (r *projectRepository) Delete(id uint) error {
	return r.db.Delete(&models.Project{}, id).Error
}

func (r *projectRepository) SetPostedAt(id uint, t *time.Time) error {
	// Se t for nil, o GORM define como NULL no banco (remove a postagem)
	return r.db.Model(&models.Project{}).Where("id = ?", id).Update("posted_at", t).Error
}

func (r *projectRepository) Filter(page, pageSize int, categoryID uint, tagID uint, onlyPosted bool) ([]models.Project, int64, error) {
	var projects []models.Project
	var total int64

	query := r.db.Model(&models.Project{}).
		Joins("JOIN project_categories on project_categories.project_id = projects.id").
		Joins("JOIN project_tags on project_tags.project_id = projects.id").
		Distinct("projects.*")

	if categoryID > 0 {
		query = query.Where("project_categories.category_id = ?", categoryID)
	}
	if tagID > 0 {
		query = query.Where("project_tags.tag_id = ?", tagID)
	}
	if onlyPosted {
		query = query.Where("posted_at IS NOT NULL")
	}

	query.Count(&total)
	err := query.Scopes(utils.PaginateRepository(page, pageSize)).
		Preload("Tags").Preload("Categories").Find(&projects).Error

	return projects, total, err
}

func (r *projectRepository) Create(project *models.Project) error {
	// O GORM por padrão tentará criar associações se elas não existirem.
	// Se o Service já garantiu que IDs de Tags/Categorias são válidos,
	// o Save/Create fará o insert no projects e nas tabelas de junção.
	return r.db.Create(project).Error
}
