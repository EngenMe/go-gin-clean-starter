package config_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
	"gorm.io/gorm"

	"github.com/Caknoooo/go-gin-clean-starter/constants"
	"github.com/Caknoooo/go-gin-clean-starter/tests/integration/container"
)

// DatabaseConfigTestSuite is a test suite for managing database configuration and container-based test setups.
// It includes utility functions for setup, teardown, and validation of database connections in a controlled environment.
// This suite uses a test container to mock the database and manage its lifecycle during testing sessions.
type DatabaseConfigTestSuite struct {
	container.BaseSuite
	dbContainer *container.TestDatabaseContainer
	db          *gorm.DB
}

// SetupSuite initializes the test suite by starting a database container and configuring required environment variables.
func (s *DatabaseConfigTestSuite) SetupSuite() {
	s.BaseSuite.SetupSuite()

	dbContainer, err := container.StartTestContainer()
	require.NoError(s.T(), err)
	s.dbContainer = dbContainer

	envVars := map[string]string{
		"APP_ENV": constants.ENUM_RUN_TESTING,
		"DB_HOST": dbContainer.Host,
		"DB_PORT": dbContainer.Port,
		"DB_USER": container.GetEnvWithDefault("DB_USER", "testuser"),
		"DB_PASS": container.GetEnvWithDefault("DB_PASS", "testpassword"),
		"DB_NAME": container.GetEnvWithDefault("DB_NAME", "testdb"),
	}
	s.SetupEnv(envVars)
}

// TearDownSuite cleans up resources after tests, including closing the database, stopping the container, and resetting environment variables.
func (s *DatabaseConfigTestSuite) TearDownSuite() {
	if s.db != nil {
		require.NoError(s.T(), container.CloseDatabaseConnection(s.db))
	}

	if s.dbContainer != nil {
		require.NoError(s.T(), s.dbContainer.Stop())
	}

	s.CleanupEnv(
		[]string{
			"APP_ENV",
			"DB_USER",
			"DB_PASS",
			"DB_NAME",
			"DB_HOST",
			"DB_PORT",
		},
	)
}

// TestSetUpDatabaseConnection tests the setup of a database connection and validates its functionality and extensions.
func (s *DatabaseConfigTestSuite) TestSetUpDatabaseConnection() {
	db := container.SetUpDatabaseConnection()
	s.db = db

	var result int
	err := db.Raw("SELECT 1").Scan(&result).Error
	require.NoError(s.T(), err)
	assert.Equal(s.T(), 1, result)

	var extensions []string
	err = db.Raw("SELECT extname FROM pg_extension WHERE extname = 'uuid-ossp'").Scan(&extensions).Error
	require.NoError(s.T(), err)
	assert.NotEmpty(s.T(), extensions)
}

// TestCloseDatabaseConnection verifies the database connection can be closed properly and the closure is successfully enforced.
func (s *DatabaseConfigTestSuite) TestCloseDatabaseConnection() {
	db := container.SetUpDatabaseConnection()
	s.db = db

	var result int
	err := db.Raw("SELECT 1").Scan(&result).Error
	require.NoError(s.T(), err)

	err = container.CloseDatabaseConnection(db)
	if err != nil {
		panic(err)
	}

	dbSQL, err := db.DB()
	require.NoError(s.T(), err)
	err = dbSQL.Ping()
	require.Error(s.T(), err)
}

// TestDatabaseConfigTestSuite runs the DatabaseConfigTestSuite to validate database configuration and lifecycle management.
func TestDatabaseConfigTestSuite(t *testing.T) {
	suite.Run(t, new(DatabaseConfigTestSuite))
}
