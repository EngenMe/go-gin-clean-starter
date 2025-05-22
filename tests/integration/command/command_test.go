package command_test

import (
	"os"
	"testing"

	"github.com/samber/do"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"gorm.io/gorm"

	"github.com/Caknoooo/go-gin-clean-starter/command"
	"github.com/Caknoooo/go-gin-clean-starter/constants"
	"github.com/Caknoooo/go-gin-clean-starter/entity"
	"github.com/Caknoooo/go-gin-clean-starter/tests/integration/container"
)

// CommandTestSuite is a test suite structure for testing commands functionality with a dependency injection container.
// It includes a GORM database connection and arguments management in the test setup and teardown process.
type CommandTestSuite struct {
	suite.Suite
	injector *do.Injector
	db       *gorm.DB
	oldArgs  []string
}

// SetupSuite initializes the test suite by configuring dependencies, starting a test container, and setting up a database.
func (suite *CommandTestSuite) SetupSuite() {
	suite.injector = do.New()

	dbContainer, err := container.StartTestContainer()
	if err != nil {
		suite.T().Fatalf("Failed to start test container: %v", err)
	}

	os.Setenv("DB_HOST", dbContainer.Host)
	os.Setenv("DB_PORT", dbContainer.Port)
	os.Setenv("DB_USER", "testuser")
	os.Setenv("DB_PASS", "testpassword")
	os.Setenv("DB_NAME", "testdb")

	db := container.SetUpDatabaseConnection()
	suite.db = db

	do.ProvideNamed[*gorm.DB](
		suite.injector, constants.DB, func(i *do.Injector) (*gorm.DB, error) {
			return db, nil
		},
	)

	suite.oldArgs = os.Args
}

// TearDownSuite resets the command-line arguments and closes the database connection if it exists. Logs any errors.
func (suite *CommandTestSuite) TearDownSuite() {
	os.Args = suite.oldArgs

	if suite.db != nil {
		if err := container.CloseDatabaseConnection(suite.db); err != nil {
			suite.T().Logf("Failed to close database connection: %v", err)
		}
	}
}

// TestCommands_Migrate verifies that the migration process is executed correctly and the necessary tables are created in the database.
func (suite *CommandTestSuite) TestCommands_Migrate() {
	os.Args = []string{"cmd", "--migrate"}

	result := command.Commands(suite.injector)

	assert.False(suite.T(), result, "Expected run to be false when migrate flag is set")

	assert.True(suite.T(), suite.db.Migrator().HasTable("users"), "Users table should exist after migration")
	assert.True(
		suite.T(),
		suite.db.Migrator().HasTable("refresh_tokens"),
		"Refresh tokens table should exist after migration",
	)
}

// TestCommands_Seed validates that the seeding process populates the database and correctly sets the seed flag behavior.
func (suite *CommandTestSuite) TestCommands_Seed() {
	suite.db.AutoMigrate(&entity.User{})

	os.Args = []string{"cmd", "--seed"}

	result := command.Commands(suite.injector)

	assert.False(suite.T(), result, "Expected run to be false when seed flag is set")

	var count int64
	suite.db.Model(&entity.User{}).Count(&count)
	assert.Greater(suite.T(), count, int64(0), "Expected users to be seeded")
}

// TestCommands_Script verifies the behavior of the Commands function when the --script flag is provided with a script name.
func (suite *CommandTestSuite) TestCommands_Script() {
	os.Args = []string{"cmd", "--script:example_script"}

	result := command.Commands(suite.injector)

	assert.False(suite.T(), result, "Expected run to be false when script flag is set")
}

// TestCommands_Run verifies that the Commands function returns true when the --run flag is set in the command-line arguments.
func (suite *CommandTestSuite) TestCommands_Run() {
	os.Args = []string{"cmd", "--run"}

	result := command.Commands(suite.injector)

	assert.True(suite.T(), result, "Expected run to be true when run flag is set")
}

// TestCommands_NoFlags verifies the behavior of the Commands function when no command-line flags are provided.
func (suite *CommandTestSuite) TestCommands_NoFlags() {
	os.Args = []string{"cmd"}

	result := command.Commands(suite.injector)

	assert.False(suite.T(), result, "Expected run to be false when no flags are set")
}

// TestCommandTestSuite runs the test suite for the CommandTestSuite structure by invoking the testify suite's Run method.
func TestCommandTestSuite(t *testing.T) {
	suite.Run(t, new(CommandTestSuite))
}
