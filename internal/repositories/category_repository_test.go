package repositories_test

import (
	"cms-headless/internal/models"
	"cms-headless/internal/repositories"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestCategoryRepository(t *testing.T) {
	// Cada execução de TestCategoryRepository ganha seu próprio DB isolado
	db := SetupTestDB()
	repo := repositories.NewCategoryRepository(db)
	now := time.Now().UTC().Truncate(time.Second)

	t.Run("Deve criar uma categoria com sucesso", func(t *testing.T) {
		cat := &models.Category{Title: "Design", CreatedAt: now, UpdatedAt: now}
		err := repo.Create(cat)

		assert.NoError(t, err)
		assert.NotZero(t, cat.ID)
		assert.Equal(t, uint(1), cat.ID)
		assert.Equal(t, "Design", cat.Title)
		assert.Equal(t, now, cat.CreatedAt)
		assert.Equal(t, now, cat.UpdatedAt)
	})

	t.Run("Deve listar categorias com paginação e total", func(t *testing.T) {
		// Seed específico para este run
		db.Create(&models.Category{Title: "A", CreatedAt: now, UpdatedAt: now})
		db.Create(&models.Category{Title: "B", CreatedAt: now, UpdatedAt: now})

		cats, total, err := repo.FindAll(1, 10)

		assert.NoError(t, err)
		assert.Equal(t, total, int64(3))
		assert.Equal(t, "A", cats[0].Title)
		assert.Equal(t, "B", cats[1].Title)
		assert.Equal(t, "Design", cats[2].Title)
		assert.Equal(t, uint(2), cats[0].ID)
		assert.Equal(t, uint(3), cats[1].ID)
		assert.Equal(t, uint(1), cats[2].ID)
		assert.Equal(t, now, cats[2].CreatedAt)
		assert.Equal(t, now, cats[1].UpdatedAt)
		assert.True(t, cats[0].Title < cats[1].Title, "Deveria estar ordenado por título")
	})

	t.Run("Deve renomear uma categoria sem afetar outras", func(t *testing.T) {
		cat := &models.Category{Title: "Original"}
		db.Create(cat)

		err := repo.UpdateName(cat.ID, "Renomeado")
		assert.NoError(t, err)

		var check models.Category
		db.First(&check, cat.ID)
		assert.Equal(t, "Renomeado", check.Title)
	})

	t.Run("Não deve permitir categorias com títulos duplicados", func(t *testing.T) {
		db.Create(&models.Category{Title: "Repetido"})
		err := repo.Create(&models.Category{Title: "Repetido"})

		assert.Error(t, err, "O banco deveria barrar títulos duplicados via UniqueIndex")
		assert.EqualError(t, err, "UNIQUE constraint failed: categories.title")
	})
}
