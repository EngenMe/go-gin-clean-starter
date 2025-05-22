package repository

import (
	"context"
	"time"

	"gorm.io/gorm"

	"github.com/Caknoooo/go-gin-clean-starter/entity"
)

// RefreshTokenRepository is an interface for managing refresh tokens in a database.
// Create adds a new refresh token to the database.
// FindByToken retrieves a refresh token record by its token value.
// DeleteByUserID removes all refresh tokens associated with a specific user ID.
// DeleteByToken deletes a refresh token by its token value.
// DeleteExpired removes all expired refresh tokens from the database.
type RefreshTokenRepository interface {
	Create(ctx context.Context, tx *gorm.DB, token entity.RefreshToken) (entity.RefreshToken, error)
	FindByToken(ctx context.Context, tx *gorm.DB, token string) (entity.RefreshToken, error)
	DeleteByUserID(ctx context.Context, tx *gorm.DB, userID string) error
	DeleteByToken(ctx context.Context, tx *gorm.DB, token string) error
	DeleteExpired(ctx context.Context, tx *gorm.DB) error
}

// refreshTokenRepository is a struct that implements the RefreshTokenRepository interface for managing refresh tokens.
// It provides methods to create, retrieve, and delete refresh tokens in a database using GORM.
// The struct contains a GORM DB pointer for database operations.
type refreshTokenRepository struct {
	db *gorm.DB
}

// NewRefreshTokenRepository creates a new instance of RefreshTokenRepository for managing refresh tokens using GORM.
func NewRefreshTokenRepository(db *gorm.DB) RefreshTokenRepository {
	return &refreshTokenRepository{
		db: db,
	}
}

// Create inserts a new refresh token into the database and returns the created token or an error if the operation fails.
func (r *refreshTokenRepository) Create(
	ctx context.Context,
	tx *gorm.DB,
	token entity.RefreshToken,
) (entity.RefreshToken, error) {
	if tx == nil {
		tx = r.db
	}

	if err := tx.WithContext(ctx).Create(&token).Error; err != nil {
		return entity.RefreshToken{}, err
	}

	return token, nil
}

// FindByToken retrieves a refresh token record by its token string, including the associated user, using the provided DB transaction.
func (r *refreshTokenRepository) FindByToken(ctx context.Context, tx *gorm.DB, token string) (
	entity.RefreshToken,
	error,
) {
	if tx == nil {
		tx = r.db
	}

	var refreshToken entity.RefreshToken
	if err := tx.WithContext(ctx).Where("token = ?", token).Preload("User").Take(&refreshToken).Error; err != nil {
		return entity.RefreshToken{}, err
	}

	return refreshToken, nil
}

// DeleteByUserID removes all refresh tokens associated with the given user ID from the database. Returns an error if it fails.
func (r *refreshTokenRepository) DeleteByUserID(ctx context.Context, tx *gorm.DB, userID string) error {
	if tx == nil {
		tx = r.db
	}

	if err := tx.WithContext(ctx).Where("user_id = ?", userID).Delete(&entity.RefreshToken{}).Error; err != nil {
		return err
	}

	return nil
}

// DeleteByToken deletes a refresh token record from the database based on the provided token string. Returns an error if it fails.
func (r *refreshTokenRepository) DeleteByToken(ctx context.Context, tx *gorm.DB, token string) error {
	if tx == nil {
		tx = r.db
	}

	if err := tx.WithContext(ctx).Where("token = ?", token).Delete(&entity.RefreshToken{}).Error; err != nil {
		return err
	}

	return nil
}

// DeleteExpired removes all refresh tokens from the database that have expired based on the `expires_at` timestamp.
func (r *refreshTokenRepository) DeleteExpired(ctx context.Context, tx *gorm.DB) error {
	if tx == nil {
		tx = r.db
	}

	if err := tx.WithContext(ctx).Where("expires_at < ?", time.Now()).Delete(&entity.RefreshToken{}).Error; err != nil {
		return err
	}

	return nil
}
