package repository

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"

	"github.com/Caknoooo/go-gin-clean-starter/dto"
)

// TestPaginate tests the Paginate function to ensure it correctly calculates and applies offset and limit to a Gorm query.
func TestPaginate(t *testing.T) {
	var tests []struct {
		name     string
		req      dto.PaginationRequest
		expected struct {
			offset int
			limit  int
		}
	}

	for _, tt := range tests {
		t.Run(
			tt.name, func(t *testing.T) {
				db := &gorm.DB{
					Statement: &gorm.Statement{
						Vars: make([]interface{}, 0),
					},
				}

				result := Paginate(tt.req)(db)

				assert.Equal(t, tt.expected.offset, result.Statement.Offset)
				assert.Equal(t, tt.expected.limit, result.Statement.Limit)
			},
		)
	}
}

// TestTotalPage validates the TotalPage function by testing various scenarios for item counts and items per page.
func TestTotalPage(t *testing.T) {
	tests := []struct {
		name     string
		count    int64
		perPage  int64
		expected int64
	}{
		{
			name:     "20 items with 10 per page",
			count:    20,
			perPage:  10,
			expected: 2,
		},
		{
			name:     "21 items with 10 per page",
			count:    21,
			perPage:  10,
			expected: 3,
		},
		{
			name:     "0 items with 10 per page",
			count:    0,
			perPage:  10,
			expected: 0,
		},
		{
			name:     "5 items with 10 per page",
			count:    5,
			perPage:  10,
			expected: 1,
		},
	}

	for _, tt := range tests {
		t.Run(
			tt.name, func(t *testing.T) {
				result := TotalPage(tt.count, tt.perPage)
				assert.Equal(t, tt.expected, result)
			},
		)
	}
}
