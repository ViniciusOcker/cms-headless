package utils

import "gorm.io/gorm"

// Helper para paginação segura
func PaginateRepository(page, pageSize int) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		if page <= 0 {
			page = 1
		}
		switch {
		case pageSize > 500:
			pageSize = 500
		case pageSize <= 0:
			pageSize = 50
		}
		offset := (page - 1) * pageSize
		return db.Offset(offset).Limit(pageSize)
	}
}
