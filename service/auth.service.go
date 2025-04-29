// service/auth_service.go
package service

import (
	"context"
	"time"

	"github.com/Caknoooo/go-gin-clean-starter/config"
	"github.com/Caknoooo/go-gin-clean-starter/dto"
	"github.com/Caknoooo/go-gin-clean-starter/entity"
	"github.com/Caknoooo/go-gin-clean-starter/helpers"
	"github.com/Caknoooo/go-gin-clean-starter/repository"
	"github.com/google/uuid"
)

type (
	AuthService interface {
		Login(ctx context.Context, req dto.UserLoginRequest) (dto.TokenResponse, error)
		RefreshToken(ctx context.Context, req dto.RefreshTokenRequest) (dto.TokenResponse, error)
		Logout(ctx context.Context, token string) error
	}

	authService struct {
		userRepo   repository.UserRepository
		tokenRepo  repository.TokenRepository
		jwtService JWTService
		config     *config.Config
	}
)

func NewAuthService(
	userRepo repository.UserRepository,
	tokenRepo repository.TokenRepository,
	jwtService JWTService,
	config *config.Config,
) AuthService {
	return &authService{
		userRepo:   userRepo,
		tokenRepo:  tokenRepo,
		jwtService: jwtService,
		config:     config,
	}
}

func (s *authService) Login(ctx context.Context, req dto.UserLoginRequest) (dto.TokenResponse, error) {
	// Verify user credentials
	user, flag, err := s.userRepo.CheckEmail(ctx, nil, req.Email)
	if err != nil || !flag {
		return dto.TokenResponse{}, dto.ErrEmailNotFound
	}

	if !user.IsVerified {
		return dto.TokenResponse{}, dto.ErrAccountNotVerified
	}

	// Verify password
	checkPassword, err := helpers.CheckPassword(user.Password, []byte(req.Password))
	if err != nil || !checkPassword {
		return dto.TokenResponse{}, dto.ErrPasswordNotMatch
	}

	// Generate access token
	accessToken, expiresIn, err := s.jwtService.GenerateAccessToken(user.ID.String(), user.Role)
	if err != nil {
		return dto.TokenResponse{}, err
	}

	// Generate refresh token
	refreshToken := uuid.New().String()
	refreshTokenExpiry := time.Now().Add(time.Hour * time.Duration(s.config.JWT.RefreshTokenExpiry))

	// Store refresh token in database
	refreshTokenEntity := entity.RefreshToken{
		UserID:    user.ID,
		Token:     refreshToken,
		ExpiresAt: refreshTokenExpiry,
	}

	// Delete any existing refresh tokens for this user (optional)
	_ = s.tokenRepo.DeleteRefreshTokensByUserID(ctx, nil, user.ID)

	// Save new refresh token
	if _, err := s.tokenRepo.CreateRefreshToken(ctx, nil, refreshTokenEntity); err != nil {
		return dto.TokenResponse{}, err
	}

	return dto.TokenResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		ExpiresIn:    expiresIn,
		TokenType:    "Bearer",
	}, nil
}

func (s *authService) RefreshToken(ctx context.Context, req dto.RefreshTokenRequest) (dto.TokenResponse, error) {
	// Verify the refresh token exists and is valid
	storedToken, err := s.tokenRepo.GetRefreshToken(ctx, nil, req.RefreshToken)
	if err != nil {
		return dto.TokenResponse{}, dto.ErrInvalidRefreshToken
	}

	// Check if the token has expired
	if time.Now().After(storedToken.ExpiresAt) {
		// Delete the expired token
		_ = s.tokenRepo.DeleteRefreshToken(ctx, nil, req.RefreshToken)
		return dto.TokenResponse{}, dto.ErrExpiredRefreshToken
	}

	// Get the user associated with the token
	user, err := s.userRepo.GetUserById(ctx, nil, storedToken.UserID.String())
	if err != nil {
		return dto.TokenResponse{}, dto.ErrUserNotFound
	}

	// Generate a new access token
	accessToken, expiresIn, err := s.jwtService.GenerateAccessToken(user.ID.String(), user.Role)
	if err != nil {
		return dto.TokenResponse{}, err
	}

	// For security, we can optionally generate a new refresh token each time
	if s.config.JWT.RotateRefreshToken {
		// Delete the old refresh token
		_ = s.tokenRepo.DeleteRefreshToken(ctx, nil, req.RefreshToken)

		// Generate a new refresh token
		newRefreshToken := uuid.New().String()
		refreshTokenExpiry := time.Now().Add(time.Hour * time.Duration(s.config.JWT.RefreshTokenExpiry))

		// Store the new refresh token
		newRefreshTokenEntity := entity.RefreshToken{
			UserID:    user.ID,
			Token:     newRefreshToken,
			ExpiresAt: refreshTokenExpiry,
		}

		if _, err := s.tokenRepo.CreateRefreshToken(ctx, nil, newRefreshTokenEntity); err != nil {
			return dto.TokenResponse{}, err
		}

		return dto.TokenResponse{
			AccessToken:  accessToken,
			RefreshToken: newRefreshToken,
			ExpiresIn:    expiresIn,
			TokenType:    "Bearer",
		}, nil
	}

	// If we don't rotate refresh tokens, return the same refresh token
	return dto.TokenResponse{
		AccessToken:  accessToken,
		RefreshToken: req.RefreshToken,
		ExpiresIn:    expiresIn,
		TokenType:    "Bearer",
	}, nil
}

func (s *authService) Logout(ctx context.Context, token string) error {
	// Delete the refresh token to invalidate it
	return s.tokenRepo.DeleteRefreshToken(ctx, nil, token)
}
