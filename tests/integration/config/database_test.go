package config_test

import (
	"os"
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
	suite.Suite
	dbContainer *container.TestDatabaseContainer
	db          *gorm.DB
}

// SetupSuite initializes the test suite by starting a database container and setting required environment variables.
func (suite *DatabaseConfigTestSuite) SetupSuite() {
	container, err := container.StartTestContainer()
	require.NoError(suite.T(), err)
	suite.dbContainer = container

	err = os.Setenv("APP_ENV", constants.ENUM_RUN_TESTING)
	if err != nil {
		panic(err)
	}
	err = os.Setenv("DB_USER", "testuser")
	if err != nil {
		panic(err)
	}
	err = os.Setenv("DB_PASS", "testpassword")
	if err != nil {
		panic(err)
	}
	err = os.Setenv("DB_NAME", "testdb")
	if err != nil {
		panic(err)
	}
	err = os.Setenv("DB_HOST", container.Host)
	if err != nil {
		panic(err)
	}
	err = os.Setenv("DB_PORT", container.Port)
	if err != nil {
		panic(err)
	}
}

// TearDownSuite cleans up the test suite by closing database connections, stopping containers, and unsetting environment variables.
func (suite *DatabaseConfigTestSuite) TearDownSuite() {
	if suite.db != nil {
		err := container.CloseDatabaseConnection(suite.db)
		if err != nil {
			panic(err)
		}
	}

	if suite.dbContainer != nil {
		err := suite.dbContainer.Stop()
		require.NoError(suite.T(), err)
	}

	err := os.Unsetenv("APP_ENV")
	if err != nil {
		panic(err)
	}
	err = os.Unsetenv("DB_USER")
	if err != nil {
		panic(err)
	}
	err = os.Unsetenv("DB_PASS")
	if err != nil {
		panic(err)
	}
	err = os.Unsetenv("DB_NAME")
	if err != nil {
		panic(err)
	}
	err = os.Unsetenv("DB_HOST")
	if err != nil {
		panic(err)
	}
	err = os.Unsetenv("DB_PORT")
	if err != nil {
		panic(err)
	}
}

// TestSetUpDatabaseConnection tests the setup of a database connection and validates its functionality and extensions.
func (suite *DatabaseConfigTestSuite) TestSetUpDatabaseConnection() {
	db := container.SetUpDatabaseConnection()
	suite.db = db

	var result int
	err := db.Raw("SELECT 1").Scan(&result).Error
	require.NoError(suite.T(), err)
	assert.Equal(suite.T(), 1, result)

	var extensions []string
	err = db.Raw("SELECT extname FROM pg_extension WHERE extname = 'uuid-ossp'").Scan(&extensions).Error
	require.NoError(suite.T(), err)
	assert.NotEmpty(suite.T(), extensions)
}

// TestCloseDatabaseConnection verifies the database connection can be closed properly and the closure is successfully enforced.
func (suite *DatabaseConfigTestSuite) TestCloseDatabaseConnection() {
	db := container.SetUpDatabaseConnection()
	suite.db = db

	var result int
	err := db.Raw("SELECT 1").Scan(&result).Error
	require.NoError(suite.T(), err)

	err = container.CloseDatabaseConnection(db)
	if err != nil {
		panic(err)
	}

	dbSQL, err := db.DB()
	require.NoError(suite.T(), err)
	err = dbSQL.Ping()
	require.Error(suite.T(), err)
}

// TestDatabaseConfigTestSuite runs the DatabaseConfigTestSuite to validate database configuration and lifecycle management.
func TestDatabaseConfigTestSuite(t *testing.T) {
	suite.Run(t, new(DatabaseConfigTestSuite))
}
