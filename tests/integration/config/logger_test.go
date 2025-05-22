package config_test

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	"github.com/Caknoooo/go-gin-clean-starter/config"
	"github.com/Caknoooo/go-gin-clean-starter/tests/integration/container"
)

// LoggerIntegrationTestSuite is a test suite for verifying database logging and integration functionality.
type LoggerIntegrationTestSuite struct {
	container.BaseSuite
	dbContainer *container.TestDatabaseContainer
	db          *gorm.DB
	testLogDir  string
}

// SetupSuite initializes the test suite by setting up a test database, environment variables, and logging configuration.
func (s *LoggerIntegrationTestSuite) SetupSuite() {
	s.BaseSuite.SetupSuite()

	s.testLogDir = "./test_logs_integration"
	config.LogDir = s.testLogDir
	require.NoError(s.T(), os.MkdirAll(s.testLogDir, 0755))

	dbContainer, err := container.StartTestContainer()
	require.NoError(s.T(), err)
	s.dbContainer = dbContainer

	envVars := map[string]string{
		"DB_HOST": dbContainer.Host,
		"DB_PORT": dbContainer.Port,
		"DB_USER": container.GetEnvWithDefault("DB_USER", "testuser"),
		"DB_PASS": container.GetEnvWithDefault("DB_PASS", "testpassword"),
		"DB_NAME": container.GetEnvWithDefault("DB_NAME", "testdb"),
	}
	s.SetupEnv(envVars)

	dsn := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%s sslmode=disable",
		os.Getenv("DB_HOST"),
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASS"),
		os.Getenv("DB_NAME"),
		os.Getenv("DB_PORT"),
	)

	s.db, err = gorm.Open(
		postgres.Open(dsn), &gorm.Config{
			Logger: config.SetupLogger(),
		},
	)
	require.NoError(s.T(), err)

	err = s.db.Exec("CREATE EXTENSION IF NOT EXISTS \"uuid-ossp\"").Error
	require.NoError(s.T(), err)
}

// TearDownSuite cleans up resources used during the test suite, including closing the database, stopping containers, and resetting environment variables.
func (s *LoggerIntegrationTestSuite) TearDownSuite() {
	if s.db != nil {
		sqlDB, err := s.db.DB()
		if err == nil {
			require.NoError(s.T(), sqlDB.Close())
		}
	}

	if s.dbContainer != nil {
		require.NoError(s.T(), s.dbContainer.Stop())
	}

	s.CleanupEnv(
		[]string{
			"DB_HOST",
			"DB_PORT",
			"DB_USER",
			"DB_PASS",
			"DB_NAME",
		},
	)

	require.NoError(s.T(), os.RemoveAll(s.testLogDir))
}

// TestLoggerWithDatabaseOperations verifies that database operations are logged correctly and log files are properly generated.
func (s *LoggerIntegrationTestSuite) TestLoggerWithDatabaseOperations() {
	type TestModel struct {
		ID   uint `gorm:"primaryKey"`
		Name string
	}

	err := s.db.AutoMigrate(&TestModel{})
	require.NoError(s.T(), err)

	tests := []struct {
		name string
		op   func() error
	}{
		{
			name: "Create record",
			op: func() error {
				return s.db.Create(&TestModel{Name: "test"}).Error
			},
		},
		{
			name: "Find existing record",
			op: func() error {
				var result TestModel
				return s.db.First(&result, 1).Error
			},
		},
		{
			name: "Find non-existent record",
			op: func() error {
				var result TestModel
				return s.db.First(&result, 999).Error
			},
		},
	}

	for _, tt := range tests {
		s.Run(
			tt.name, func() {
				err := tt.op()
				if tt.name == "Find non-existent record" {
					assert.Error(s.T(), err)
				} else {
					assert.NoError(s.T(), err)
				}
			},
		)
	}

	currentMonth := strings.ToLower(time.Now().Format("January"))
	logFileName := fmt.Sprintf("%s_query.log", currentMonth)
	logPath := filepath.Join(s.testLogDir, logFileName)

	fileInfo, err := os.Stat(logPath)
	require.NoError(s.T(), err)
	assert.Greater(s.T(), fileInfo.Size(), int64(0))
}

// TestLoggerIntegrationTestSuite runs the LoggerIntegrationTestSuite to verify database logging and integration functionality.
func TestLoggerIntegrationTestSuite(t *testing.T) {
	suite.Run(t, new(LoggerIntegrationTestSuite))
}
