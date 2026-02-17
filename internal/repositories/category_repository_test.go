package repositories_test

import (
	"cms-headless/internal/models"
	"cms-headless/internal/repositories"
	"cms-headless/internal/utils"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCategoryRepository(t *testing.T) {
	db := utils.SetupTestDB()
	repo := repositories.NewCategoryRepository(db)

	t.Run("Deve criar uma category com sucesso", func(t *testing.T) {
		category := &models.Category{Title: "Golang"}
		err := repo.Create(category)
		assert.NoError(t, err)
		assert.NotZero(t, category.ID)
	})

	t.Run("Deve listar todas as categorias", func(t *testing.T) {
		categories, total, err := repo.FindAll(1, 50)
		assert.NoError(t, err)
		assert.Equal(t, int64(1), total)
		assert.GreaterOrEqual(t, len(categories), 1)
	})

	t.Run("Deve renomear uma category", func(t *testing.T) {
		// Criar tag inicial
		category := &models.Category{Title: "React"}
		db.Create(category)

		err := repo.UpdateName(category.ID, "React Next.js")
		assert.NoError(t, err)

		var updatedcategory models.Category
		db.First(&updatedcategory, category.ID)
		assert.Equal(t, "React Next.js", updatedcategory.Title)
	})
}
