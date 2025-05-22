package repository

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"

	"github.com/Caknoooo/go-gin-clean-starter/dto"
	"github.com/Caknoooo/go-gin-clean-starter/entity"
	"github.com/Caknoooo/go-gin-clean-starter/repository"
	"github.com/Caknoooo/go-gin-clean-starter/tests/integration/container"
)

// PaginationRequest represents parameters for paginated requests, including the page number, items per page, and a search term.
type PaginationRequest struct {
	Page    int
	PerPage int
	Search  string
}

// Default sets the default values for Page and PerPage if they are not initialized. Page defaults to 1, and PerPage defaults to 10.
func (p *PaginationRequest) Default() {
	if p.Page == 0 {
		p.Page = 1
	}
	if p.PerPage == 0 {
		p.PerPage = 10
	}
}

// PaginationResponse represents pagination details for paginated API responses.
// It includes current page, items per page, total count of items, and the maximum number of pages.
type PaginationResponse struct {
	Page    int
	PerPage int
	Count   int64
	MaxPage int64
}

// GetAllUserRepositoryResponse represents a response containing a list of users and pagination details.
type GetAllUserRepositoryResponse struct {
	Users              []entity.User
	PaginationResponse PaginationResponse
}

// Paginate applies pagination logic to a Gorm database query based on the provided PaginationRequest parameters.
// It calculates the offset and limits the number of records returned per page.
// Designed to be used as a query function passed to Gorm's chainable methods.
func Paginate(req PaginationRequest) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		offset := (req.Page - 1) * req.PerPage
		return db.Offset(offset).Limit(req.PerPage)
	}
}

// TotalPage calculates the total number of pages based on the given item count and items per page.
func TotalPage(count, perPage int64) int64 {
	if perPage == 0 {
		return 0
	}
	return (count + perPage - 1) / perPage
}

// TestUserRepository is a test function for the user repository, ensuring proper functionality of CRUD operations.
func TestUserRepository(t *testing.T) {
	container.LoadTestEnv()

	dbContainer, err := container.StartTestContainer()
	if err != nil {
		t.Fatalf("failed to start test container: %v", err)
	}
	defer func(dbContainer *container.TestDatabaseContainer) {
		err := dbContainer.Stop()
		if err != nil {
			panic(err)
		}
	}(dbContainer)

	envVars := map[string]string{
		"DB_HOST": dbContainer.Host,
		"DB_PORT": dbContainer.Port,
		"DB_USER": container.GetEnvWithDefault("DB_USER", "testuser"),
		"DB_PASS": container.GetEnvWithDefault("DB_PASS", "testpassword"),
		"DB_NAME": container.GetEnvWithDefault("DB_NAME", "testdb"),
	}
	if err := container.SetEnv(envVars); err != nil {
		panic(fmt.Sprintf("Failed to set env vars: %v", err))
	}

	db := container.SetUpDatabaseConnection()
	defer func(db *gorm.DB) {
		err := container.CloseDatabaseConnection(db)
		if err != nil {
			panic(err)
		}
	}(db)

	err = db.AutoMigrate(&entity.User{})
	if err != nil {
		t.Fatalf("failed to migrate schema: %v", err)
	}

	repo := repository.NewUserRepository(db)

	cleanDB := func() {
		err := db.Exec("TRUNCATE TABLE users RESTART IDENTITY CASCADE").Error
		if err != nil {
			t.Fatalf("failed to clean database: %v", err)
		}
	}

	ctx := context.Background()

	t.Run(
		"Register", func(t *testing.T) {
			t.Cleanup(cleanDB)
			user := entity.User{
				ID:         uuid.New(),
				Name:       "Test User",
				Email:      "test@example.com",
				TelpNumber: "1234567890",
				Password:   "password123",
				Role:       "user",
				Timestamp: entity.Timestamp{
					CreatedAt: time.Now(),
					UpdatedAt: time.Now(),
				},
			}

			createdUser, err := repo.Register(ctx, nil, user)
			assert.NoError(t, err)
			assert.Equal(t, user.Name, createdUser.Name)
			assert.Equal(t, user.Email, createdUser.Email)
			assert.Equal(t, user.ID, createdUser.ID)
		},
	)

	t.Run(
		"GetAllUserWithPagination", func(t *testing.T) {
			t.Cleanup(cleanDB)
			users := []entity.User{
				{
					ID:         uuid.New(),
					Name:       "User One",
					Email:      "user1@example.com",
					TelpNumber: "1111111111",
					Password:   "password123",
					Role:       "user",
					Timestamp: entity.Timestamp{
						CreatedAt: time.Now(),
						UpdatedAt: time.Now(),
					},
				},
				{
					ID:         uuid.New(),
					Name:       "User Two",
					Email:      "user2@example.com",
					TelpNumber: "2222222222",
					Password:   "password123",
					Role:       "user",
					Timestamp: entity.Timestamp{
						CreatedAt: time.Now(),
						UpdatedAt: time.Now(),
					},
				},
			}

			for _, u := range users {
				_, err := repo.Register(ctx, nil, u)
				assert.NoError(t, err)
			}

			req := dto.PaginationRequest{
				Page:    1,
				PerPage: 1,
				Search:  "User",
			}

			response, err := repo.GetAllUserWithPagination(ctx, nil, req)
			assert.NoError(t, err)
			assert.Len(t, response.Users, 1)
			assert.Equal(t, int64(2), response.PaginationResponse.Count)
			assert.Equal(t, int64(2), response.PaginationResponse.MaxPage)
			assert.Equal(t, 1, response.PaginationResponse.Page)
			assert.Equal(t, 1, response.PaginationResponse.PerPage)
		},
	)

	t.Run(
		"GetUserById", func(t *testing.T) {
			t.Cleanup(cleanDB)
			user := entity.User{
				ID:         uuid.New(),
				Name:       "ID Test User",
				Email:      "idtest@example.com",
				TelpNumber: "3333333333",
				Password:   "password123",
				Role:       "user",
				Timestamp: entity.Timestamp{
					CreatedAt: time.Now(),
					UpdatedAt: time.Now(),
				},
			}
			_, err := repo.Register(ctx, nil, user)
			assert.NoError(t, err)

			foundUser, err := repo.GetUserById(ctx, nil, user.ID.String())
			assert.NoError(t, err)
			assert.Equal(t, user.ID, foundUser.ID)
			assert.Equal(t, user.Name, foundUser.Name)
			assert.Equal(t, user.Email, foundUser.Email)
		},
	)

	t.Run(
		"GetUserByEmail", func(t *testing.T) {
			t.Cleanup(cleanDB)
			user := entity.User{
				ID:         uuid.New(),
				Name:       "Email Test User",
				Email:      "emailtest@example.com",
				TelpNumber: "4444444444",
				Password:   "password123",
				Role:       "user",
				Timestamp: entity.Timestamp{
					CreatedAt: time.Now(),
					UpdatedAt: time.Now(),
				},
			}
			_, err := repo.Register(ctx, nil, user)
			assert.NoError(t, err)

			foundUser, err := repo.GetUserByEmail(ctx, nil, "emailtest@example.com")
			assert.NoError(t, err)
			assert.Equal(t, user.ID, foundUser.ID)
			assert.Equal(t, user.Name, foundUser.Name)
			assert.Equal(t, user.Email, foundUser.Email)
		},
	)

	t.Run(
		"CheckEmail", func(t *testing.T) {
			t.Cleanup(cleanDB)
			user := entity.User{
				ID:         uuid.New(),
				Name:       "Check Email User",
				Email:      "checkemail@example.com",
				TelpNumber: "5555555555",
				Password:   "password123",
				Role:       "user",
				Timestamp: entity.Timestamp{
					CreatedAt: time.Now(),
					UpdatedAt: time.Now(),
				},
			}
			_, err := repo.Register(ctx, nil, user)
			assert.NoError(t, err)

			foundUser, exists, err := repo.CheckEmail(ctx, nil, "checkemail@example.com")
			assert.NoError(t, err)
			assert.True(t, exists)
			assert.Equal(t, user.ID, foundUser.ID)
			assert.Equal(t, user.Email, foundUser.Email)

			_, exists, err = repo.CheckEmail(ctx, nil, "nonexistent@example.com")
			assert.ErrorIs(t, err, gorm.ErrRecordNotFound)
			assert.False(t, exists)
		},
	)

	t.Run(
		"Update", func(t *testing.T) {
			t.Cleanup(cleanDB)
			user := entity.User{
				ID:         uuid.New(),
				Name:       "Update Test User",
				Email:      "updatetest@example.com",
				TelpNumber: "6666666666",
				Password:   "password123",
				Role:       "user",
				Timestamp: entity.Timestamp{
					CreatedAt: time.Now(),
					UpdatedAt: time.Now(),
				},
			}
			_, err := repo.Register(ctx, nil, user)
			assert.NoError(t, err)

			updatedUser := entity.User{
				ID:         user.ID,
				Name:       "Updated User",
				Email:      "updated@example.com",
				TelpNumber: "7777777777",
				Role:       "admin",
				Timestamp: entity.Timestamp{
					UpdatedAt: time.Now(),
				},
			}

			result, err := repo.Update(ctx, nil, updatedUser)
			assert.NoError(t, err)
			assert.Equal(t, updatedUser.Name, result.Name)
			assert.Equal(t, updatedUser.Email, result.Email)
			assert.Equal(t, updatedUser.TelpNumber, result.TelpNumber)
			assert.Equal(t, updatedUser.Role, result.Role)
		},
	)

	t.Run(
		"Delete", func(t *testing.T) {
			t.Cleanup(cleanDB)
			user := entity.User{
				ID:         uuid.New(),
				Name:       "Delete Test User",
				Email:      "deletetest@example.com",
				TelpNumber: "8888888888",
				Password:   "password123",
				Role:       "user",
				Timestamp: entity.Timestamp{
					CreatedAt: time.Now(),
					UpdatedAt: time.Now(),
				},
			}
			_, err := repo.Register(ctx, nil, user)
			assert.NoError(t, err)

			err = repo.Delete(ctx, nil, user.ID.String())
			assert.NoError(t, err)

			_, err = repo.GetUserById(ctx, nil, user.ID.String())
			assert.ErrorIs(t, err, gorm.ErrRecordNotFound)
		},
	)
}
