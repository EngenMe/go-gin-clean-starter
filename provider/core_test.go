package provider

import (
	"testing"

	"github.com/samber/do"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"gorm.io/gorm"

	"github.com/Caknoooo/go-gin-clean-starter/config"
	"github.com/Caknoooo/go-gin-clean-starter/constants"
	"github.com/Caknoooo/go-gin-clean-starter/service"
)

// mockConfig is a mock implementation that embeds mock.Mock to simulate and test behaviors in database configuration.
type mockConfig struct {
	mock.Mock
}

// SetUpDatabaseConnection mocks the initialization of a Gorm database connection and simulates error handling.
func (m *mockConfig) SetUpDatabaseConnection() *gorm.DB {
	args := m.Called()
	db := args.Get(0).(*gorm.DB)
	if err := args.Error(1); err != nil {
		panic(err)
	}
	return db
}

// mockUserProvider is a mock implementation used for testing user-related dependency provision in a dependency injection container.
type mockUserProvider struct {
	mock.Mock
}

// ProvideUserDependencies registers user-related dependencies into the provided dependency injection container.
func (m *mockUserProvider) ProvideUserDependencies(injector *do.Injector) {
	m.Called(injector)
}

// TestInitDatabase validates the integration of InitDatabase with a dependency injector and ensures correct database setup.
func TestInitDatabase(t *testing.T) {
	injector := do.New()

	mockCfg := &mockConfig{}
	mockDB := &gorm.DB{}
	mockCfg.On("SetUpDatabaseConnection").Return(mockDB, nil)
	originalSetUp := config.SetUpDatabaseConnection
	config.SetUpDatabaseConnection = mockCfg.SetUpDatabaseConnection
	defer func() { config.SetUpDatabaseConnection = originalSetUp }()

	InitDatabase(injector)

	db, err := do.InvokeNamed[*gorm.DB](injector, constants.DB)
	assert.NoError(t, err, "should provide DB without error")
	assert.Equal(t, mockDB, db, "should provide the mock DB")
	mockCfg.AssertExpectations(t)
}

// TestRegisterDependencies validates the registration process of application dependencies in the dependency injection container.
func TestRegisterDependencies(t *testing.T) {
	injector := do.New()

	mockCfg := &mockConfig{}
	mockDB := &gorm.DB{}
	mockCfg.On("SetUpDatabaseConnection").Return(mockDB, nil)
	originalSetUp := config.SetUpDatabaseConnection
	config.SetUpDatabaseConnection = mockCfg.SetUpDatabaseConnection
	defer func() { config.SetUpDatabaseConnection = originalSetUp }()

	mockUserProv := &mockUserProvider{}
	mockUserProv.On("ProvideUserDependencies", injector).Return()
	originalProvide := ProvideUserDependencies
	ProvideUserDependencies = mockUserProv.ProvideUserDependencies
	defer func() { ProvideUserDependencies = originalProvide }()

	RegisterDependencies(injector)

	db, err := do.InvokeNamed[*gorm.DB](injector, constants.DB)
	assert.NoError(t, err, "should provide DB without error")
	assert.Equal(t, mockDB, db, "should provide the mock DB")

	jwtService, err := do.InvokeNamed[service.JWTService](injector, constants.JWTService)
	assert.NoError(t, err, "should provide JWTService without error")
	assert.NotNil(t, jwtService, "JWTService should not be nil")

	mockUserProv.AssertExpectations(t)
	mockCfg.AssertExpectations(t)
}
