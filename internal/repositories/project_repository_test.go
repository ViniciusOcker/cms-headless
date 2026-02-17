package repositories_test

import (
	"cms-headless/internal/models"
	"cms-headless/internal/repositories"
	"cms-headless/internal/utils"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestProjectRepository(t *testing.T) {
	db := utils.SetupTestDB()
	repo := repositories.NewProjectRepository(db)

	// Seed de dados básicos para relacionamento
	tag := models.Tag{Title: "React"}
	cat := models.Category{Title: "Frontend"}
	db.Create(&tag)
	db.Create(&cat)

	t.Run("Deve criar um projeto com tags e categorias", func(t *testing.T) {
		project := &models.Project{
			Title:      "Meu Portfólio",
			Slug:       "meu-portfolio",
			Tags:       []models.Tag{tag},
			Categories: []models.Category{cat},
		}

		err := repo.Create(project)
		assert.NoError(t, err)
		assert.NotZero(t, project.ID)

		// Verificar se as associações foram criadas no banco
		var countTags int64
		db.Table("project_tags").Where("project_id = ?", project.ID).Count(&countTags)
		assert.Equal(t, int64(1), countTags)
	})

	t.Run("Deve filtrar apenas projetos postados", func(t *testing.T) {
		now := time.Now()
		db.Create(&models.Project{Title: "Postado", Slug: "postado", PostedAt: &now})
		db.Create(&models.Project{Title: "Rascunho", Slug: "rascunho", PostedAt: nil}) // Zero value

		// Testando FindAll com onlyPosted = true
		projects, total, err := repo.FindAll(1, 10, true)
		assert.NoError(t, err)
		assert.Equal(t, int64(1), total)
		assert.Equal(t, "Postado", projects[0].Title)
	})

	t.Run("Deve buscar por Slug respeitando visibilidade", func(t *testing.T) {
		slug := "projeto-secreto"
		db.Create(&models.Project{Title: "Secreto", Slug: slug}) // Sem PostedAt

		// Não deve achar se onlyPosted for true
		p, err := repo.FindBySlug(slug, true)
		assert.Error(t, err) // GORM retorna RecordNotFound
		assert.Nil(t, p)

		// Deve achar se onlyPosted for false
		p, err = repo.FindBySlug(slug, false)
		assert.NoError(t, err)
		assert.NotNil(t, p)
	})

	t.Run("Deve realizar Soft Delete", func(t *testing.T) {
		p := models.Project{Title: "Para Deletar", Slug: "deletar"}
		db.Create(&p)

		err := repo.Delete(p.ID)
		assert.NoError(t, err)

		// O registro não deve ser encontrado em uma busca comum
		var check models.Project
		err = db.First(&check, p.ID).Error
		assert.Error(t, err)

		// Mas deve existir no banco (Soft Delete check)
		err = db.Unscoped().First(&check, p.ID).Error
		assert.NoError(t, err)
		assert.NotNil(t, check.DeletedAt)
	})

	t.Run("Deve substituir tags (ReplaceTags)", func(t *testing.T) {
		newTag := models.Tag{Title: "Go"}
		db.Create(&newTag)

		p, _ := repo.FindBySlug("meu-portfolio", false)

		err := repo.ReplaceTags(p, []models.Tag{newTag})
		assert.NoError(t, err)

		// Recarregar e verificar
		pUpdated, _ := repo.FindByID(p.ID)
		assert.Len(t, pUpdated.Tags, 1)
		assert.Equal(t, "Go", pUpdated.Tags[0].Title)
	})
}
