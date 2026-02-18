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

	t.Run("Deve substituir tags via Association Replace", func(t *testing.T) {
		newCat := models.Tag{Title: "DevOps"}
		db.Create(&newCat)

		p, _ := repo.FindBySlug("aprendendo-go", false)

		err := repo.ReplaceTags(p, []models.Tag{newCat})
		assert.NoError(t, err)

		// Recarregar e validar
		pUpdated, _ := repo.FindByID(p.ID)
		assert.Len(t, pUpdated.Tags, 1)
		assert.Equal(t, "DevOps", pUpdated.Tags[0].Title)
	})

	t.Run("Deve mudar o titulo e slug para MESH e manter associações", func(t *testing.T) {
		item, err := repo.FindBySlug("aprendendo-go", false)
		assert.NoError(t, err)
		item.Title = "MESH"
		item.Slug = "mesh"
		err = repo.Update(item)
		assert.NoError(t, err)
		updated, _ := repo.FindByID(1)
		assert.Equal(t, "MESH", updated.Title)
		assert.Equal(t, "mesh", updated.Slug)
		assert.NotEmpty(t, updated.Tags)
		assert.NotEmpty(t, updated.Categories)
	})
}

func TestPostRepository_Search(t *testing.T) {
	db := SetupTestDB()
	repo := repositories.NewPostRepository(db)
	now := time.Now().UTC()

	// Seed: Criar Tags e Categorias
	t1 := models.Tag{Title: "Go"}
	t2 := models.Tag{Title: "JS"}
	c1 := models.Category{Title: "Backend"}
	db.Create(&t1)
	db.Create(&t2)
	db.Create(&c1)

	// Seed: Criar Posts
	p1 := models.Post{Title: "Mastering Go", Tags: []models.Tag{t1}, Categories: []models.Category{c1}, PostedAt: &now}
	p2 := models.Post{Title: "Learning JS", Tags: []models.Tag{t2}, PostedAt: &now}
	db.Create(&p1)
	db.Create(&p2)

	t.Run("Deve filtrar por Categoria e Tag simultaneamente", func(t *testing.T) {
		// Slugs únicos são obrigatórios
		p1 := models.Post{Title: "Mastering Go", Slug: "mastering-go", Tags: []models.Tag{t1}, Categories: []models.Category{c1}, PostedAt: &now}
		p2 := models.Post{Title: "Learning JS", Slug: "learning-js", Tags: []models.Tag{t2}, PostedAt: &now}
		db.Create(&p1)
		db.Create(&p2)

		res, total, err := repo.Search(1, 10, c1.ID, t1.ID, "", true)

		assert.NoError(t, err)
		assert.Equal(t, int64(1), total)
		// Se o total for > 0, aí sim acessamos o índice 0 para evitar panic
		if assert.NotEmpty(t, res) {
			assert.Equal(t, "Mastering Go", res[0].Title)
		}
	})

	t.Run("Deve filtrar por texto no título", func(t *testing.T) {
		res, total, err := repo.Search(1, 10, 0, 0, "Learning", true)

		assert.NoError(t, err)
		assert.Equal(t, int64(1), total)
		assert.Equal(t, "Learning JS", res[0].Title)
	})

	t.Run("Não deve retornar nada para busca sem resultados", func(t *testing.T) {
		res, total, err := repo.Search(1, 10, 999, 0, "", true)

		assert.NoError(t, err)
		assert.Equal(t, int64(0), total)
		assert.Empty(t, res)
	})
}
