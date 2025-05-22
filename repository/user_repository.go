package repository

import (
	"context"

	"gorm.io/gorm"

	"github.com/Caknoooo/go-gin-clean-starter/dto"
	"github.com/Caknoooo/go-gin-clean-starter/entity"
)

type (
	// UserRepository defines the contract for database operations related to user management, including CRUD and search functionalities.
	UserRepository interface {
		Register(ctx context.Context, tx *gorm.DB, user entity.User) (entity.User, error)
		GetAllUserWithPagination(
			ctx context.Context,
			tx *gorm.DB,
			req dto.PaginationRequest,
		) (dto.GetAllUserRepositoryResponse, error)
		GetUserById(ctx context.Context, tx *gorm.DB, userId string) (entity.User, error)
		GetUserByEmail(ctx context.Context, tx *gorm.DB, email string) (entity.User, error)
		CheckEmail(ctx context.Context, tx *gorm.DB, email string) (entity.User, bool, error)
		Update(ctx context.Context, tx *gorm.DB, user entity.User) (entity.User, error)
		Delete(ctx context.Context, tx *gorm.DB, userId string) error
	}

	// userRepository struct represents the implementation of UserRepository interface using GORM for database operations.
	userRepository struct {
		db *gorm.DB
	}
)

// NewUserRepository initializes and returns a new instance of UserRepository with the provided GORM database connection.
func NewUserRepository(db *gorm.DB) UserRepository {
	return &userRepository{
		db: db,
	}
}

// Register inserts a new user record into the database and returns the created user or an error if the operation fails.
func (r *userRepository) Register(ctx context.Context, tx *gorm.DB, user entity.User) (entity.User, error) {
	if tx == nil {
		tx = r.db
	}

	if err := tx.WithContext(ctx).Create(&user).Error; err != nil {
		return entity.User{}, err
	}

	return user, nil
}

// GetAllUserWithPagination retrieves a paginated list of users and total count based on the provided pagination request.
func (r *userRepository) GetAllUserWithPagination(
	ctx context.Context,
	tx *gorm.DB,
	req dto.PaginationRequest,
) (dto.GetAllUserRepositoryResponse, error) {
	if tx == nil {
		tx = r.db
	}

	var users []entity.User
	var err error
	var count int64

	req.Default()

	query := tx.WithContext(ctx).Model(&entity.User{})
	if req.Search != "" {
		query = query.Where("name LIKE ?", "%"+req.Search+"%")
	}

	if err := query.Count(&count).Error; err != nil {
		return dto.GetAllUserRepositoryResponse{}, err
	}

	if err := query.Scopes(Paginate(req)).Find(&users).Error; err != nil {
		return dto.GetAllUserRepositoryResponse{}, err
	}

	totalPage := TotalPage(count, int64(req.PerPage))
	return dto.GetAllUserRepositoryResponse{
		Users: users,
		PaginationResponse: dto.PaginationResponse{
			Page:    req.Page,
			PerPage: req.PerPage,
			Count:   count,
			MaxPage: totalPage,
		},
	}, err
}

// GetUserById retrieves a user by their unique ID using the provided context and database transaction.
func (r *userRepository) GetUserById(ctx context.Context, tx *gorm.DB, userId string) (entity.User, error) {
	if tx == nil {
		tx = r.db
	}

	var user entity.User
	if err := tx.WithContext(ctx).Where("id = ?", userId).Take(&user).Error; err != nil {
		return entity.User{}, err
	}

	return user, nil
}

// GetUserByEmail retrieves a user record from the database by email using the provided context and transaction.
func (r *userRepository) GetUserByEmail(ctx context.Context, tx *gorm.DB, email string) (entity.User, error) {
	if tx == nil {
		tx = r.db
	}

	var user entity.User
	if err := tx.WithContext(ctx).Where("email = ?", email).Take(&user).Error; err != nil {
		return entity.User{}, err
	}

	return user, nil
}

// CheckEmail verifies the existence of a user by email and returns the user entity and a boolean indicating existence.
func (r *userRepository) CheckEmail(ctx context.Context, tx *gorm.DB, email string) (entity.User, bool, error) {
	if tx == nil {
		tx = r.db
	}

	var user entity.User
	if err := tx.WithContext(ctx).Where("email = ?", email).Take(&user).Error; err != nil {
		return entity.User{}, false, err
	}

	return user, true, nil
}

// Update updates an existing user record in the database and returns the updated user or an error if the operation fails.
func (r *userRepository) Update(ctx context.Context, tx *gorm.DB, user entity.User) (entity.User, error) {
	if tx == nil {
		tx = r.db
	}

	if err := tx.WithContext(ctx).Updates(&user).Error; err != nil {
		return entity.User{}, err
	}

	return user, nil
}

// Delete removes a user identified by userId from the database, using the provided context and optional transaction.
func (r *userRepository) Delete(ctx context.Context, tx *gorm.DB, userId string) error {
	if tx == nil {
		tx = r.db
	}

	if err := tx.WithContext(ctx).Delete(&entity.User{}, "id = ?", userId).Error; err != nil {
		return err
	}

	return nil
}
