package main

import (
	"os"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/samber/do"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/Caknoooo/go-gin-clean-starter/command"
	"github.com/Caknoooo/go-gin-clean-starter/provider"
	"github.com/Caknoooo/go-gin-clean-starter/routes"
)

// MockCommand is a mock type that embeds mock.Mock to facilitate testing of command-related behaviors.
type MockCommand struct {
	mock.Mock
}

// Commands invokes the mock implementation for processing command-line arguments and returns the mocked boolean result.
func (m *MockCommand) Commands(injector *do.Injector) bool {
	args := m.Called(injector)
	return args.Bool(0)
}

// MockProvider is a mock implementation for testing purposes, embedding the mock.Mock type for mock functionalities.
type MockProvider struct {
	mock.Mock
}

// RegisterDependencies registers necessary dependencies into the provided dependency injector instance.
func (m *MockProvider) RegisterDependencies(injector *do.Injector) {
	m.Called(injector)
}

// MockRoutes is a mock structure embedding `mock.Mock`, used for testing route registrations in a Gin application.
type MockRoutes struct {
	mock.Mock
}

// RegisterRoutes registers routes into the given Gin engine using the provided dependency injector.
func (m *MockRoutes) RegisterRoutes(server *gin.Engine, injector *do.Injector) {
	m.Called(server, injector)
}

// TestArgs tests the behavior of the args function with different command-line argument scenarios.
func TestArgs(t *testing.T) {
	t.Run(
		"with no arguments", func(t *testing.T) {
			oldArgs := os.Args
			defer func() { os.Args = oldArgs }()
			os.Args = []string{"test"}

			injector := do.New()
			result := args(injector)

			assert.True(t, result)
		},
	)

	t.Run(
		"with arguments", func(t *testing.T) {
			oldArgs := os.Args
			defer func() { os.Args = oldArgs }()
			os.Args = []string{"test", "some-command"}

			mockCmd := new(MockCommand)
			mockCmd.On("Commands", mock.Anything).Return(false)
			command.Commands = mockCmd.Commands

			injector := do.New()
			result := args(injector)

			assert.False(t, result)
			mockCmd.AssertExpectations(t)
		},
	)
}

// TestRun validates the behavior of the run function in different scenarios such as custom ports and various environments.
func TestRun(t *testing.T) {
	t.Run(
		"with custom port", func(t *testing.T) {
			oldPort := os.Getenv("PORT")
			defer func() { os.Setenv("PORT", oldPort) }()
			os.Setenv("PORT", "9999")

			server := gin.Default()
			called := false
			originalRun := run
			defer func() { run = originalRun }()
			run = func(s *gin.Engine) {
				called = true
				assert.Equal(t, server, s)
				assert.Equal(t, "9999", os.Getenv("PORT"))
			}

			run(server)

			assert.True(t, called)
		},
	)

	t.Run(
		"dev environment", func(t *testing.T) {
			oldEnv := os.Getenv("APP_ENV")
			defer func() { os.Setenv("APP_ENV", oldEnv) }()
			os.Setenv("APP_ENV", "dev")

			server := gin.Default()
			called := false
			originalRun := run
			defer func() { run = originalRun }()
			run = func(s *gin.Engine) {
				called = true
				assert.Equal(t, server, s)
			}

			run(server)

			assert.True(t, called)
		},
	)

	t.Run(
		"prod environment", func(t *testing.T) {
			oldEnv := os.Getenv("APP_ENV")
			defer func() { os.Setenv("APP_ENV", oldEnv) }()
			os.Setenv("APP_ENV", "prod")

			server := gin.Default()
			called := false
			originalRun := run
			defer func() { run = originalRun }()
			run = func(s *gin.Engine) {
				called = true
				assert.Equal(t, server, s)
			}

			run(server)

			assert.True(t, called)
		},
	)
}

// TestMainFunc is a test helper function that validates the main function's behavior under different scenarios.
func TestMainFunc(t *testing.T) {
	t.Run(
		"successful execution", func(t *testing.T) {
			mockProvider := new(MockProvider)
			mockProvider.On("RegisterDependencies", mock.Anything).Return()
			provider.RegisterDependencies = mockProvider.RegisterDependencies

			mockRoutes := new(MockRoutes)
			mockRoutes.On("RegisterRoutes", mock.Anything, mock.Anything).Return()
			routes.RegisterRoutes = mockRoutes.RegisterRoutes

			originalArgs := args
			defer func() { args = originalArgs }()
			args = func(injector *do.Injector) bool {
				return true
			}

			runCalled := false
			originalRun := run
			defer func() { run = originalRun }()
			run = func(server *gin.Engine) {
				runCalled = true
			}

			main()

			assert.True(t, runCalled)
			mockProvider.AssertExpectations(t)
			mockRoutes.AssertExpectations(t)
		},
	)

	t.Run(
		"early return from args", func(t *testing.T) {
			mockProvider := new(MockProvider)
			mockProvider.On("RegisterDependencies", mock.Anything).Return()
			provider.RegisterDependencies = mockProvider.RegisterDependencies

			originalArgs := args
			defer func() { args = originalArgs }()
			args = func(injector *do.Injector) bool {
				return false
			}

			runCalled := false
			originalRun := run
			defer func() { run = originalRun }()
			run = func(server *gin.Engine) {
				runCalled = true
			}

			main()

			assert.False(t, runCalled)
			mockProvider.AssertExpectations(t)
		},
	)
}
