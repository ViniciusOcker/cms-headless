package repositories_test

import (
	"cms-headless/internal/models"
	"cms-headless/internal/repositories"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestTagRepository(t *testing.T) {
	db := SetupTestDB()
	repo := repositories.NewTagRepository(db)
	now := time.Now().UTC().Truncate(time.Second)

	t.Run("Deve criar uma tag com sucesso", func(t *testing.T) {
		tag := &models.Tag{Title: "React", CreatedAt: now, UpdatedAt: now}
		err := repo.Create(tag)

		assert.NoError(t, err)
		assert.Equal(t, uint(1), tag.ID)
		assert.Equal(t, "React", tag.Title)
		assert.WithinDuration(t, now, tag.CreatedAt, time.Second)
	})

	t.Run("Deve listar tags com paginação e ordenação", func(t *testing.T) {
		// Seed
		db.Create(&models.Tag{Title: "Golang", CreatedAt: now})
		db.Create(&models.Tag{Title: "Vue", CreatedAt: now})

		tags, total, err := repo.FindAll(1, 10)

		assert.NoError(t, err)
		assert.Equal(t, int64(3), total) // Golang (1) + React (2) + Vue (3)

		// Verificando ordenação alfabética: Golang, React, Vue
		assert.Equal(t, "Golang", tags[0].Title)
		assert.Equal(t, "React", tags[1].Title)
		assert.Equal(t, "Vue", tags[2].Title)

		assert.Equal(t, uint(2), tags[0].ID)
		assert.Equal(t, uint(1), tags[1].ID)
		assert.Equal(t, uint(3), tags[2].ID)
	})

	t.Run("Deve renomear uma tag mantendo o ID", func(t *testing.T) {
		tag := &models.Tag{Title: "OldName"}
		db.Create(tag)

		err := repo.UpdateName(tag.ID, "NewName")
		assert.NoError(t, err)

		var check models.Tag
		db.First(&check, tag.ID)
		assert.Equal(t, "NewName", check.Title)
		assert.Equal(t, tag.ID, check.ID)
	})

	t.Run("Não deve permitir tags duplicadas", func(t *testing.T) {
		err := repo.Create(&models.Tag{Title: "Duplicate"})
		assert.NoError(t, err)

		err = repo.Create(&models.Tag{Title: "Duplicate"})
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "UNIQUE constraint failed")
	})
}
