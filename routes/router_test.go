package routes

import (
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/samber/do"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockUserRouteRegistrar is a mock implementation of a user route registrar for testing route registration functionality.
type MockUserRouteRegistrar struct {
	mock.Mock
}

// User registers user-related routes to the provided Gin engine using the supplied dependency injector.
func (m *MockUserRouteRegistrar) User(server *gin.Engine, injector *do.Injector) {
	m.Called(server, injector)
}

// TestRegisterRoutes tests the RegisterRoutes function to ensure routing logic and dependency injection work as expected.
func TestRegisterRoutes(t *testing.T) {
	mockEngine := gin.Default()
	mockInjector := do.New()

	t.Run(
		"Successfully registers routes", func(t *testing.T) {
			mockUserRegistrar := new(MockUserRouteRegistrar)

			mockUserRegistrar.On("User", mockEngine, mockInjector).Once()

			originalUser := User
			User = mockUserRegistrar.User
			defer func() { User = originalUser }()

			RegisterRoutes(mockEngine, mockInjector)

			mockUserRegistrar.AssertExpectations(t)
		},
	)

	t.Run(
		"Passes correct parameters", func(t *testing.T) {
			mockUserRegistrar := new(MockUserRouteRegistrar)

			mockUserRegistrar.On("User", mock.AnythingOfType("*gin.Engine"), mock.AnythingOfType("*do.Injector")).Once()

			originalUser := User
			User = mockUserRegistrar.User
			defer func() { User = originalUser }()

			RegisterRoutes(mockEngine, mockInjector)

			mockUserRegistrar.AssertExpectations(t)

			args := mockUserRegistrar.Calls[0].Arguments
			serverArg := args.Get(0).(*gin.Engine)
			injectorArg := args.Get(1).(*do.Injector)

			assert.Equal(t, mockEngine, serverArg)
			assert.Equal(t, mockInjector, injectorArg)
		},
	)
}
