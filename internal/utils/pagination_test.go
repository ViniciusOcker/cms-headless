package utils_test

import (
	"cms-headless/internal/models"
	"cms-headless/internal/utils"
	"testing"

	"github.com/stretchr/testify/assert"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

// Setup minimalista apenas para análise de SQL
func setupMinimalDB() *gorm.DB {
	db, _ := gorm.Open(sqlite.Open(""), &gorm.Config{
		DryRun: true, // Garante que não tentará executar nada
	})
	return db
}

func TestPaginateRepository(t *testing.T) {
	db := setupMinimalDB()

	t.Run("Deve aplicar valores padrão (LIMIT 50) quando entradas forem zero", func(t *testing.T) {
		sql := db.ToSQL(func(tx *gorm.DB) *gorm.DB {
			return tx.Model(&models.Category{}).Scopes(utils.PaginateRepository(0, 0)).Find(&[]models.Category{})
		})

		// Esperado: Sem OFFSET (0) e LIMIT 50
		assert.Contains(t, sql, "LIMIT 50")
		assert.NotContains(t, sql, "OFFSET")
	})

	t.Run("Deve calcular OFFSET corretamente para página 2", func(t *testing.T) {
		sql := db.ToSQL(func(tx *gorm.DB) *gorm.DB {
			return tx.Model(&models.Project{}).Scopes(utils.PaginateRepository(2, 10)).Find(&[]models.Project{})
		})

		// Page 2, Size 10 -> Offset 10, Limit 10
		assert.Contains(t, sql, "LIMIT 10")
		assert.Contains(t, sql, "OFFSET 10")
	})

	t.Run("Deve limitar o teto em 500 itens", func(t *testing.T) {
		sql := db.ToSQL(func(tx *gorm.DB) *gorm.DB {
			return tx.Model(&models.Post{}).Scopes(utils.PaginateRepository(1, 999)).Find(&[]models.Post{})
		})

		assert.Contains(t, sql, "LIMIT 500")
	})

	t.Run("Página negativa deve resultar em OFFSET 0", func(t *testing.T) {
		sql := db.ToSQL(func(tx *gorm.DB) *gorm.DB {
			return tx.Model(&models.Tag{}).Scopes(utils.PaginateRepository(-5, 20)).Find(&[]models.Tag{})
		})

		assert.Contains(t, sql, "LIMIT 20")
		assert.NotContains(t, sql, "OFFSET")
	})
}
