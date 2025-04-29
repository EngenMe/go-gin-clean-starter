// repository/token_repository.go
package repository

import (
	"context"
	"time"

	"github.com/Caknoooo/go-gin-clean-starter/entity"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type (
	TokenRepository interface {
		CreateRefreshToken(ctx context.Context, tx *gorm.DB, token entity.RefreshToken) (entity.RefreshToken, error)
		GetRefreshToken(ctx context.Context, tx *gorm.DB, token string) (entity.RefreshToken, error)
		DeleteRefreshToken(ctx context.Context, tx *gorm.DB, token string) error
		DeleteRefreshTokensByUserID(ctx context.Context, tx *gorm.DB, userID uuid.UUID) error
		DeleteExpiredTokens(ctx context.Context, tx *gorm.DB) error
	}

	tokenRepository struct {
		db *gorm.DB
	}
)

func NewTokenRepository(db *gorm.DB) TokenRepository {
	return &tokenRepository{
		db: db,
	}
}

func (r *tokenRepository) CreateRefreshToken(ctx context.Context, tx *gorm.DB, token entity.RefreshToken) (entity.RefreshToken, error) {
	if tx == nil {
		tx = r.db
	}
	if err := tx.WithContext(ctx).Create(&token).Error; err != nil {
		return entity.RefreshToken{}, err
	}
	return token, nil
}

func (r *tokenRepository) GetRefreshToken(ctx context.Context, tx *gorm.DB, token string) (entity.RefreshToken, error) {
	if tx == nil {
		tx = r.db
	}
	var refreshToken entity.RefreshToken
	if err := tx.WithContext(ctx).Where("token = ?", token).First(&refreshToken).Error; err != nil {
		return entity.RefreshToken{}, err
	}
	return refreshToken, nil
}

func (r *tokenRepository) DeleteRefreshToken(ctx context.Context, tx *gorm.DB, token string) error {
	if tx == nil {
		tx = r.db
	}
	if err := tx.WithContext(ctx).Where("token = ?", token).Delete(&entity.RefreshToken{}).Error; err != nil {
		return err
	}
	return nil
}

func (r *tokenRepository) DeleteRefreshTokensByUserID(ctx context.Context, tx *gorm.DB, userID uuid.UUID) error {
	if tx == nil {
		tx = r.db
	}
	if err := tx.WithContext(ctx).Where("user_id = ?", userID).Delete(&entity.RefreshToken{}).Error; err != nil {
		return err
	}
	return nil
}

func (r *tokenRepository) DeleteExpiredTokens(ctx context.Context, tx *gorm.DB) error {
	if tx == nil {
		tx = r.db
	}
	if err := tx.WithContext(ctx).Where("expires_at < ?", time.Now()).Delete(&entity.RefreshToken{}).Error; err != nil {
		return err
	}
	return nil
}
