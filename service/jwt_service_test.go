package service

import (
	"os"
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v4"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestNewJWTService validates the creation of a new JWTService instance and checks its configurations and behavior.
func TestNewJWTService(t *testing.T) {
	os.Setenv("JWT_SECRET", "test-secret")
	defer os.Unsetenv("JWT_SECRET")

	service := NewJWTService()
	assert.NotNil(t, service)
	jwtService, ok := service.(*jwtService)
	assert.True(t, ok)
	assert.Equal(t, "test-secret", jwtService.secretKey)
	assert.Equal(t, "Template", jwtService.issuer)
	assert.Equal(t, time.Minute*15, jwtService.accessExpiry)
	assert.Equal(t, time.Hour*24*7, jwtService.refreshExpiry)

	os.Unsetenv("JWT_SECRET")
	service = NewJWTService()
	assert.Equal(t, "test-secret", jwtService.secretKey)
}

// TestGenerateAccessToken tests the functionality of generating a valid JWT access token, parsing, and validating its claims.
func TestGenerateAccessToken(t *testing.T) {
	service := NewJWTService()
	userID := "test-user"
	role := "admin"

	token := service.GenerateAccessToken(userID, role)
	assert.NotEmpty(t, token)

	parsedToken, err := jwt.ParseWithClaims(
		token, &jwtCustomClaim{}, func(token *jwt.Token) (interface{}, error) {
			return []byte("Template"), nil
		},
	)
	require.NoError(t, err)
	assert.True(t, parsedToken.Valid)

	claims, ok := parsedToken.Claims.(*jwtCustomClaim)
	assert.True(t, ok)
	assert.Equal(t, userID, claims.UserID)
	assert.Equal(t, role, claims.Role)
	assert.Equal(t, "Template", claims.Issuer)
	assert.WithinDuration(t, time.Now().Add(time.Minute*15), claims.ExpiresAt.Time, time.Second)
}

// TestGenerateRefreshToken tests the GenerateRefreshToken method to ensure it creates a valid token and expiration timestamp.
func TestGenerateRefreshToken(t *testing.T) {
	service := NewJWTService()

	token, expiresAt := service.GenerateRefreshToken()
	assert.NotEmpty(t, token)
	assert.False(t, expiresAt.IsZero())
	assert.WithinDuration(t, time.Now().Add(time.Hour*24*7), expiresAt, time.Second)

	assert.Len(t, token, 44)
}

// TestValidateToken tests the functionality of token validation, including handling valid, invalid, and incorrectly signed tokens.
func TestValidateToken(t *testing.T) {
	service := NewJWTService()

	validToken := service.GenerateAccessToken("test-user", "admin")
	parsedToken, err := service.ValidateToken(validToken)
	assert.NoError(t, err)
	assert.True(t, parsedToken.Valid)

	_, err = service.ValidateToken("invalid.token.string")
	assert.Error(t, err)

	wrongMethodToken := jwt.NewWithClaims(
		jwt.SigningMethodHS512,
		&jwtCustomClaim{
			UserID: "test-user",
			Role:   "admin",
			RegisteredClaims: jwt.RegisteredClaims{
				ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Minute * 15)),
				Issuer:    "Template",
				IssuedAt:  jwt.NewNumericDate(time.Now()),
			},
		},
	)
	signedWrong, err := wrongMethodToken.SignedString([]byte("Template"))
	require.NoError(t, err, "Failed to sign token with wrong method")

	_, err = service.ValidateToken(signedWrong)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "unexpected signing method")
}

// TestGetUserIDByToken verifies the GetUserIDByToken method by testing valid and invalid JWT token scenarios.
func TestGetUserIDByToken(t *testing.T) {
	service := NewJWTService()
	userID := "test-user"

	validToken := service.GenerateAccessToken(userID, "admin")
	retrievedID, err := service.GetUserIDByToken(validToken)
	assert.NoError(t, err)
	assert.Equal(t, userID, retrievedID)

	_, err = service.GetUserIDByToken("invalid.token.string")
	assert.Error(t, err)
}

// TestGetSecretKey verifies that the getSecretKey function returns the correct value based on the presence of environment variables.
func TestGetSecretKey(t *testing.T) {
	os.Setenv("JWT_SECRET", "custom-secret")
	defer os.Unsetenv("JWT_SECRET")
	assert.Equal(t, "custom-secret", getSecretKey())

	os.Unsetenv("JWT_SECRET")
	assert.Equal(t, "Template", getSecretKey())
}
