package repositories_test

import (
	"cms-headless/internal/models"
	"cms-headless/internal/repositories"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestProjectRepository(t *testing.T) {
	db := SetupTestDB()
	repo := repositories.NewProjectRepository(db)
	now := time.Now().UTC().Truncate(time.Second)

	// Seed de dependências
	tag := models.Tag{Title: "Go"}
	cat := models.Category{Title: "Backend"}
	db.Create(&tag) // ID 1
	db.Create(&cat) // ID 1

	t.Run("Deve criar um projeto com associações", func(t *testing.T) {
		project := &models.Project{
			Title:      "Projeto Alpha",
			Slug:       "projeto-alpha",
			Tags:       []models.Tag{tag},
			Categories: []models.Category{cat},
			CreatedAt:  now,
			UpdatedAt:  now,
		}

		err := repo.Create(project)
		assert.NoError(t, err)
		assert.Equal(t, uint(1), project.ID)

		// Validar junção manual no banco para garantir o "match"
		var countTags int64
		db.Table("project_tags").Where("project_id = ?", project.ID).Count(&countTags)
		assert.Equal(t, int64(1), countTags)
	})

	t.Run("Deve retornar NIL e erro ao não encontrar slug", func(t *testing.T) {
		// Teste crucial para validar seu ajuste de retorno
		project, err := repo.FindBySlug("slug-inexistente", false)

		assert.Error(t, err)
		assert.EqualError(t, err, "record not found")
		assert.Nil(t, project) // Aqui validamos se o ponteiro é realmente nulo
	})

	t.Run("Deve ocultar rascunho quando onlyPosted for true", func(t *testing.T) {
		// Projeto rascunho (PostedAt = nil)
		db.Create(&models.Project{Title: "Rascunho", Slug: "rascunho", PostedAt: nil, CreatedAt: now, UpdatedAt: now})

		project, err := repo.FindBySlug("rascunho", true)

		assert.Error(t, err) // Deve dar RecordNotFound
		assert.EqualError(t, err, "record not found")
		assert.Nil(t, project)
	})

	t.Run("Deve encontrar projeto postado", func(t *testing.T) {
		// Use time.Second para garantir uma margem real
		postDate := now.Add(-1 * time.Hour)

		db.Create(&models.Project{
			Title:     "Postado",
			Slug:      "postado",
			PostedAt:  &postDate,
			CreatedAt: now,
			UpdatedAt: now,
		})

		project, err := repo.FindBySlug("postado", true)

		assert.NoError(t, err)
		assert.NotNil(t, project)
		assert.Equal(t, "Postado", project.Title)
	})

	t.Run("Deve realizar Soft Delete corretamente", func(t *testing.T) {
		p := &models.Project{Title: "Deletar", Slug: "deletar"}
		db.Create(p)

		err := repo.Delete(p.ID)
		assert.NoError(t, err)

		// FindByID não deve encontrar
		found, err := repo.FindByID(p.ID)
		assert.Error(t, err)
		assert.EqualError(t, err, "record not found")
		assert.Nil(t, found)

		// Verificar se continua no banco via Unscoped (Prova do Soft Delete)
		var check models.Project
		db.Unscoped().First(&check, p.ID)
		assert.NotZero(t, check.DeletedAt)
	})

	t.Run("Deve substituir tags usando Replace", func(t *testing.T) {
		// Tag nova
		newTag := models.Tag{Title: "NextJS"}
		db.Create(&newTag) // ID 2

		p, _ := repo.FindByID(1) // Projeto Alpha
		err := repo.ReplaceTags(p, []models.Tag{newTag})
		assert.NoError(t, err)

		// Recarregar para garantir o Preload
		updated, _ := repo.FindByID(1)
		assert.Len(t, updated.Tags, 1)
		assert.Equal(t, "NextJS", updated.Tags[0].Title)
	})

	t.Run("Deve substituir categorias usando Replace", func(t *testing.T) {
		// Tag nova
		newCategory := models.Category{Title: "Frontend"}
		db.Create(&newCategory) // ID 2

		p, _ := repo.FindByID(1) // Projeto Alpha
		err := repo.ReplaceCategories(p, []models.Category{newCategory})
		assert.NoError(t, err)

		// Recarregar para garantir o Preload
		updated, _ := repo.FindByID(1)
		assert.Len(t, updated.Categories, 1)
		assert.Equal(t, "Frontend", updated.Categories[0].Title)
	})

	t.Run("Deve atualizar a data da postagem para o ano de 2015", func(t *testing.T) {
		newDate := time.Date(2015, 6, 6, 0, 0, 0, 0, time.UTC)

		// Passamos o endereço da variável diretamente
		err := repo.SetPostedAt(1, &newDate)
		assert.NoError(t, err)

		updated, _ := repo.FindByID(1)

		// Validações
		assert.NotNil(t, updated.PostedAt)
		assert.Equal(t, 2015, updated.PostedAt.Year())
		// Compara os valores das datas (Time.Equal) em vez de apenas o ponteiro
		assert.True(t, newDate.Equal(*updated.PostedAt))
	})

	t.Run("Deve remover a data de postagem ao passar NIL", func(t *testing.T) {
		// Tenta despostar o projeto ID 1
		err := repo.SetPostedAt(1, nil)
		assert.NoError(t, err)

		updated, _ := repo.FindByID(1)

		// O ponteiro deve ser nil agora (NULL no banco)
		assert.Nil(t, updated.PostedAt)
	})

	t.Run("Deve mudar o titulo e slug para MESH e manter associações", func(t *testing.T) {
		// 1. Busca o item (ID 1 já deve ter tags do teste anterior)
		item, err := repo.FindByID(1)
		assert.NoError(t, err)

		// 2. Modifica
		item.Title = "MESH"
		item.Slug = "mesh"

		// 3. Update
		err = repo.Update(item)
		assert.NoError(t, err)

		// 4. Verifica persistência
		updated, _ := repo.FindByID(1)
		assert.Equal(t, "MESH", updated.Title)
		assert.Equal(t, "mesh", updated.Slug)

		// Importante: verificar se as Tags que o Alpha tinha ainda estão lá
		// Isso garante que o r.db.Save(item) não sobrescreveu a junção com vazio
		assert.NotEmpty(t, updated.Tags)
	})
}
