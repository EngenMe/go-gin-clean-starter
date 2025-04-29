package dto

import (
	"errors"
)

const (
	// Auth Messages
	MESSAGE_FAILED_GENERATE_TOKEN = "failed to generate token"
	MESSAGE_FAILED_REFRESH_TOKEN  = "failed to refresh token"
	MESSAGE_INVALID_REFRESH_TOKEN = "invalid refresh token"
	MESSAGE_EXPIRED_REFRESH_TOKEN = "refresh token has expired"
	MESSAGE_SUCCESS_REFRESH_TOKEN = "token refreshed successfully"
)

var (
	ErrInvalidRefreshToken = errors.New("invalid refresh token")
	ErrExpiredRefreshToken = errors.New("refresh token has expired")
)

type (
	TokenResponse struct {
		AccessToken  string `json:"access_token"`
		RefreshToken string `json:"refresh_token"`
		ExpiresIn    int64  `json:"expires_in"` // seconds until access token expires
		TokenType    string `json:"token_type"` // typically "Bearer"
	}

	RefreshTokenRequest struct {
		RefreshToken string `json:"refresh_token" form:"refresh_token" binding:"required"`
	}
)
