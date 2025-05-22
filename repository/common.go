package repository

import (
	"math"

	"gorm.io/gorm"

	"github.com/Caknoooo/go-gin-clean-starter/dto"
)

// Paginate applies pagination to a Gorm database query based on the provided pagination request.
func Paginate(req dto.PaginationRequest) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		offset := (req.Page - 1) * req.PerPage
		return db.Offset(offset).Limit(req.PerPage)
	}
}

// TotalPage calculates the total number of pages required based on the total items (count) and items per page (perPage).
func TotalPage(count, perPage int64) int64 {
	totalPage := int64(math.Ceil(float64(count) / float64(perPage)))

	return totalPage
}
