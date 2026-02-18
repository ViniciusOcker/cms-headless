package repositories_test

import (
	"cms-headless/internal/models"
	"cms-headless/internal/repositories"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestPostRepository(t *testing.T) {
	db := SetupTestDB()
	repo := repositories.NewPostRepository(db)

	// Preparando ambiente com Tag e Categoria
	tag := models.Tag{Title: "Tutoriais"}
	cat := models.Category{Title: "Backend"}
	db.Create(&tag)
	db.Create(&cat)

	t.Run("Deve criar um post com associações", func(t *testing.T) {
		post := &models.Post{
			Title:      "Aprendendo Go",
			Slug:       "aprendendo-go",
			Tags:       []models.Tag{tag},
			Categories: []models.Category{cat},
		}

		err := repo.Create(post)
		assert.NoError(t, err)
		assert.NotZero(t, post.ID)

		// Validar se as associações existem na tabela de junção específica de Posts
		var countTags int64
		db.Table("post_tags").Where("post_id = ?", post.ID).Count(&countTags)
		assert.Equal(t, int64(1), countTags)
	})

	t.Run("Deve buscar post por ID com Preload", func(t *testing.T) {
		// Criar um post direto para teste
		p := models.Post{Title: "Post ID", Slug: "post-id", Tags: []models.Tag{tag}}
		db.Create(&p)

		found, err := repo.FindByID(p.ID)
		assert.NoError(t, err)
		assert.Equal(t, "Post ID", found.Title)
		assert.NotEmpty(t, found.Tags)
		assert.Equal(t, "Tutoriais", found.Tags[0].Title)
	})

	t.Run("Deve respeitar visibilidade (PostedAt) na busca por Slug", func(t *testing.T) {
		slugRascunho := "post-rascunho"
		db.Create(&models.Post{Title: "Rascunho", Slug: slugRascunho, PostedAt: nil})

		// 1. Não deve achar como publicado
		found, err := repo.FindBySlug(slugRascunho, true)
		assert.Error(t, err) // gorm.ErrRecordNotFound
		assert.Nil(t, found)

		// 2. Deve achar na busca administrativa (onlyPosted = false)
		found, err = repo.FindBySlug(slugRascunho, false)
		assert.NoError(t, err)
		assert.NotNil(t, found)
	})

	t.Run("Deve atualizar data de postagem (SetPostedAt)", func(t *testing.T) {
		p := models.Post{Title: "Agendado", Slug: "agendado"}
		db.Create(&p)

		now := time.Now()
		err := repo.SetPostedAt(p.ID, &now)
		assert.NoError(t, err)

		var updated models.Post
		db.First(&updated, p.ID)
		assert.NotNil(t, updated.PostedAt)
		assert.True(t, updated.PostedAt.After(time.Now().Add(-1*time.Minute)))
	})

	t.Run("Deve remover postagem (SetPostedAt = nil)", func(t *testing.T) {
		now := time.Now()
		p := models.Post{Title: "Publicado", Slug: "publicado", PostedAt: &now}
		db.Create(&p)

		err := repo.SetPostedAt(p.ID, nil)
		assert.NoError(t, err)

		var updated models.Post
		db.First(&updated, p.ID)
		assert.Nil(t, updated.PostedAt)
	})

	t.Run("Deve substituir categorias via Association Replace", func(t *testing.T) {
		newCat := models.Category{Title: "DevOps"}
		db.Create(&newCat)

		p, _ := repo.FindBySlug("aprendendo-go", false)

		err := repo.ReplaceCategories(p, []models.Category{newCat})
		assert.NoError(t, err)

		// Recarregar e validar
		pUpdated, _ := repo.FindByID(p.ID)
		assert.Len(t, pUpdated.Categories, 1)
		assert.Equal(t, "DevOps", pUpdated.Categories[0].Title)
	})
}
