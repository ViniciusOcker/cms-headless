package repositories_test

import (
	"cms-headless/internal/models"
	"cms-headless/internal/repositories"
	"cms-headless/internal/utils"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestTagRepository(t *testing.T) {
	db := utils.SetupTestDB()
	repo := repositories.NewTagRepository(db)

	t.Run("Deve criar uma tag com sucesso", func(t *testing.T) {
		tag := &models.Tag{Title: "Golang"}
		err := repo.Create(tag)
		assert.NoError(t, err)
		assert.NotZero(t, tag.ID)
	})

	t.Run("Deve listar todas as tags", func(t *testing.T) {
		tags, total, err := repo.FindAll(1, 50)
		assert.NoError(t, err)
		assert.Equal(t, int64(1), total)
		assert.GreaterOrEqual(t, len(tags), 1)
	})

	t.Run("Deve renomear uma tag", func(t *testing.T) {
		// Criar tag inicial
		tag := &models.Tag{Title: "React"}
		db.Create(tag)

		err := repo.UpdateName(tag.ID, "React Next.js")
		assert.NoError(t, err)

		var updatedTag models.Tag
		db.First(&updatedTag, tag.ID)
		assert.Equal(t, "React Next.js", updatedTag.Title)
	})
}
