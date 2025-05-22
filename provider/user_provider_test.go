package provider

import (
	"gorm.io/gorm"
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v4"
	"github.com/samber/do"
	"github.com/stretchr/testify/assert"

	"github.com/Caknoooo/go-gin-clean-starter/constants"
	"github.com/Caknoooo/go-gin-clean-starter/controller"
	"github.com/Caknoooo/go-gin-clean-starter/service"
)

// mockJWTService is a mock implementation of the JWTService interface, used for testing purposes.
type mockJWTService struct{}

// GenerateAccessToken generates a mock access token for the given userID and role, primarily used for testing purposes.
func (m *mockJWTService) GenerateAccessToken(_ string, _ string) string {
	return "mock-access-token"
}

// GenerateRefreshToken generates a mock refresh token and its expiration time for testing purposes.
func (m *mockJWTService) GenerateRefreshToken() (string, time.Time) {
	return "mock-refresh-token", time.Now().Add(24 * time.Hour)
}

// GetUserIDByToken extracts the user ID from the given token and returns it along with any potential error encountered.
func (m *mockJWTService) GetUserIDByToken(_ string) (string, error) {
	return "mock-user-id", nil
}

// ValidateToken validates the provided JWT token and returns a mock jwt.Token object and an error if applicable.
func (m *mockJWTService) ValidateToken(token string) (*jwt.Token, error) {
	claims := jwt.MapClaims{
		"user_id": "mock-user-id",
		"role":    "user",
	}
	return &jwt.Token{
		Raw:    token,
		Claims: claims,
		Valid:  true,
	}, nil
}

// TestProvideUserDependencies verifies that user-related dependencies such as controllers are correctly provided by the injector.
func TestProvideUserDependencies(t *testing.T) {
	injector := do.New()

	mockDB := &gorm.DB{}
	do.ProvideNamedValue[*gorm.DB](injector, constants.DB, mockDB)

	mockJWT := &mockJWTService{}
	do.ProvideNamedValue[service.JWTService](injector, constants.JWTService, mockJWT)

	ProvideUserDependencies(injector)

	userController, err := do.Invoke[controller.UserController](injector)
	assert.NoError(t, err, "should provide UserController without error")
	assert.NotNil(t, userController, "UserController should not be nil")
}

// TestProvideUserDependencies_MissingDB verifies that ProvideUserDependencies panics if the database dependency is missing.
func TestProvideUserDependencies_MissingDB(t *testing.T) {
	injector := do.New()

	mockJWT := &mockJWTService{}
	do.ProvideNamedValue[service.JWTService](injector, constants.JWTService, mockJWT)

	assert.Panics(
		t,
		func() {
			ProvideUserDependencies(injector)
		},
		"should panic when DB is missing",
	)
}

// TestProvideUserDependencies_MissingJWTService checks if ProvideUserDependencies panics when JWTService dependency is missing.
func TestProvideUserDependencies_MissingJWTService(t *testing.T) {
	injector := do.New()

	mockDB := &gorm.DB{}
	do.ProvideNamedValue[*gorm.DB](injector, constants.DB, mockDB)

	assert.Panics(
		t,
		func() {
			ProvideUserDependencies(injector)
		},
		"should panic when JWTService is missing",
	)
}
