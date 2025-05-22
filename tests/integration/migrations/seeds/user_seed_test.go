package seeds

import (
	"encoding/json"
	"fmt"
	"github.com/stretchr/testify/require"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"gorm.io/gorm"

	"github.com/Caknoooo/go-gin-clean-starter/entity"
	"github.com/Caknoooo/go-gin-clean-starter/helpers"
	"github.com/Caknoooo/go-gin-clean-starter/migrations/seeds"
	"github.com/Caknoooo/go-gin-clean-starter/tests/integration/container"
)

// SeedsTestSuite represents a test suite for database seed operations.
// It includes a database connection, test data, temporary file paths, and project configuration details.
type SeedsTestSuite struct {
	suite.Suite
	db           *gorm.DB
	testData     []SeedUserRequest
	tempJSONPath string
	projectRoot  string
	usedTestJSON bool
}

// SeedUserRequest represents the structure for seeding a user in the database, containing user details and verification status.
type SeedUserRequest struct {
	Name       string `json:"name"`
	TelpNumber string `json:"telp_number"`
	Email      string `json:"email"`
	Password   string `json:"password"`
	Role       string `json:"role"`
	IsVerified bool   `json:"is_verified"`
}

// SetupSuite initializes the test suite, setting up the database container, environment variables, and test data.
func (suite *SeedsTestSuite) SetupSuite() {
	container.LoadTestEnv()

	dbContainer, err := container.StartTestContainer()
	if err != nil {
		suite.T().Fatalf("Failed to start test container: %v", err)
	}

	envVars := map[string]string{
		"DB_HOST": dbContainer.Host,
		"DB_PORT": dbContainer.Port,
		"DB_USER": container.GetEnvWithDefault("DB_USER", "testuser"),
		"DB_PASS": container.GetEnvWithDefault("DB_PASS", "testpassword"),
		"DB_NAME": container.GetEnvWithDefault("DB_NAME", "testdb"),
	}
	if err := container.SetEnv(envVars); err != nil {
		panic(fmt.Sprintf("Failed to set env vars: %v", err))
	}

	db := container.SetUpDatabaseConnection()
	suite.db = db

	projectRoot, err := helpers.GetProjectRoot()
	if err != nil {
		suite.T().Fatalf("Failed to get project root: %v", err)
	}
	suite.projectRoot = projectRoot

	suite.tempJSONPath = filepath.Join(os.TempDir(), "test_users.json")
	suite.testData = []SeedUserRequest{
		{
			Name:       "Test User 1",
			TelpNumber: "08123456789",
			Email:      "test1@example.com",
			Password:   "password123",
			Role:       "user",
			IsVerified: true,
		},
		{
			Name:       "Test Admin",
			TelpNumber: "08123456788",
			Email:      "admin@example.com",
			Password:   "admin123",
			Role:       "admin",
			IsVerified: true,
		},
	}

	err = createTestJSONFile(suite.tempJSONPath, suite.testData)
	if err != nil {
		suite.T().Fatalf("Failed to create test JSON file: %v", err)
	}
}

// TearDownSuite performs cleanup after all tests in the suite, closing the database connection and removing temporary files.
func (suite *SeedsTestSuite) TearDownSuite() {
	if suite.db != nil {
		if err := container.CloseDatabaseConnection(suite.db); err != nil {
			suite.T().Logf("Failed to close database connection: %v", err)
		}
	}

	err := os.Remove(suite.tempJSONPath)
	require.NoError(suite.T(), err, "Failed to remove temporary JSON file")
}

// BeforeTest cleans up the database by dropping the User table and resets the `usedTestJSON` flag before running a test.
func (suite *SeedsTestSuite) BeforeTest(suiteName, testName string) {
	err := suite.db.Migrator().DropTable(&entity.User{})
	require.NoError(suite.T(), err, "Failed to drop table")
	suite.usedTestJSON = false
}

// setupTestJSON creates a test JSON file for seeding users if not already present and returns its path or an error.
func (suite *SeedsTestSuite) setupTestJSON() (string, error) {
	testSeedDir := filepath.Join(suite.projectRoot, "migrations", "json")
	err := os.MkdirAll(testSeedDir, 0755)
	if err != nil {
		return "", err
	}

	usersJSONPath := filepath.Join(testSeedDir, "users.json")
	if _, err := os.Stat(usersJSONPath); err == nil {
		return usersJSONPath, nil
	}

	testJSONPath := filepath.Join(testSeedDir, "users_test.json")
	err = copyFile(suite.tempJSONPath, testJSONPath)
	if err != nil {
		return "", err
	}
	suite.usedTestJSON = true
	return testJSONPath, nil
}

// cleanupTestJSON removes the specified JSON file if it was created during the test run, as indicated by the `usedTestJSON` flag.
func (suite *SeedsTestSuite) cleanupTestJSON(jsonPath string) {
	if suite.usedTestJSON {
		err := os.Remove(jsonPath)
		require.NoError(suite.T(), err, "Failed to remove test JSON file")
	}
}

// TestListUserSeeder_Success verifies that the seeder successfully inserts user data into the database without errors.
// It checks if all test users are seeded correctly with matching name, role, and verification status.
func (suite *SeedsTestSuite) TestListUserSeeder_Success() {
	jsonPath, err := suite.setupTestJSON()
	if err != nil {
		suite.T().Fatalf("Failed to setup test JSON: %v", err)
	}
	defer suite.cleanupTestJSON(jsonPath)

	oldGetProjectRoot := helpers.GetProjectRoot
	helpers.GetProjectRoot = func() (string, error) {
		return suite.projectRoot, nil
	}
	defer func() { helpers.GetProjectRoot = oldGetProjectRoot }()

	err = seeds.ListUserSeeder(suite.db)
	assert.NoError(suite.T(), err, "Seeder should not return error")

	var seededData []SeedUserRequest
	data, err := os.ReadFile(jsonPath)
	assert.NoError(suite.T(), err, "Should read JSON file")
	err = json.Unmarshal(data, &seededData)
	assert.NoError(suite.T(), err, "Should parse JSON file")

	var users []entity.User
	result := suite.db.Find(&users)
	assert.NoError(suite.T(), result.Error, "Should be able to query users")
	assert.Equal(suite.T(), len(seededData), int(result.RowsAffected), "Should insert all test users")

	for _, testUser := range seededData {
		var user entity.User
		err := suite.db.Where("email = ?", testUser.Email).First(&user).Error
		assert.NoError(suite.T(), err, "Should find seeded user")
		assert.Equal(suite.T(), testUser.Name, user.Name, "User name should match")
		assert.Equal(suite.T(), testUser.Role, user.Role, "User role should match")
		assert.Equal(suite.T(), testUser.IsVerified, user.IsVerified, "User verification status should match")
	}
}

// TestListUserSeeder_TableCreation verifies that the seeder correctly creates the User table if it doesn't exist in the database.
func (suite *SeedsTestSuite) TestListUserSeeder_TableCreation() {
	jsonPath, err := suite.setupTestJSON()
	if err != nil {
		suite.T().Fatalf("Failed to setup test JSON: %v", err)
	}
	defer suite.cleanupTestJSON(jsonPath)

	oldGetProjectRoot := helpers.GetProjectRoot
	helpers.GetProjectRoot = func() (string, error) {
		return suite.projectRoot, nil
	}
	defer func() { helpers.GetProjectRoot = oldGetProjectRoot }()

	err = suite.db.Migrator().DropTable(&entity.User{})
	require.NoError(suite.T(), err, "Failed to drop table")

	err = seeds.ListUserSeeder(suite.db)
	assert.NoError(suite.T(), err, "Seeder should not return error")

	hasTable := suite.db.Migrator().HasTable(&entity.User{})
	assert.True(suite.T(), hasTable, "Seeder should create table if it doesn't exist")
}

// TestListUserSeeder_DuplicateUsers ensures that the seeder does not insert duplicate user records when run multiple times.
func (suite *SeedsTestSuite) TestListUserSeeder_DuplicateUsers() {
	jsonPath, err := suite.setupTestJSON()
	if err != nil {
		suite.T().Fatalf("Failed to setup test JSON: %v", err)
	}
	defer suite.cleanupTestJSON(jsonPath)

	oldGetProjectRoot := helpers.GetProjectRoot
	helpers.GetProjectRoot = func() (string, error) {
		return suite.projectRoot, nil
	}
	defer func() { helpers.GetProjectRoot = oldGetProjectRoot }()

	err = seeds.ListUserSeeder(suite.db)
	assert.NoError(suite.T(), err, "First seeder run should not return error")

	var initialCount int64
	suite.db.Model(&entity.User{}).Count(&initialCount)

	err = seeds.ListUserSeeder(suite.db)
	assert.NoError(suite.T(), err, "Second seeder run should not return error")

	var newCount int64
	suite.db.Model(&entity.User{}).Count(&newCount)

	assert.Equal(suite.T(), initialCount, newCount, "Should not insert duplicate users")
}

// TestListUserSeeder_InvalidJSONPath verifies that the seeder returns an error when provided with an invalid JSON file path.
func (suite *SeedsTestSuite) TestListUserSeeder_InvalidJSONPath() {
	oldGetProjectRoot := helpers.GetProjectRoot
	defer func() { helpers.GetProjectRoot = oldGetProjectRoot }()

	helpers.GetProjectRoot = func() (string, error) {
		return filepath.Join(os.TempDir(), "nonexistent_project"), nil
	}

	err := seeds.ListUserSeeder(suite.db)
	assert.Error(suite.T(), err, "Should return error for invalid JSON path")
}

// TestListUserSeeder_InvalidJSONContent ensures the seeder returns an error when processing a file with invalid JSON content.
func (suite *SeedsTestSuite) TestListUserSeeder_InvalidJSONContent() {
	tempDir := suite.T().TempDir()
	testSeedDir := filepath.Join(tempDir, "migrations", "json")
	err := os.MkdirAll(testSeedDir, 0755)
	if err != nil {
		suite.T().Fatalf("Failed to create test seed directory: %v", err)
	}

	invalidJSONPath := filepath.Join(testSeedDir, "users.json")
	err = os.WriteFile(invalidJSONPath, []byte("invalid json content"), 0644)
	if err != nil {
		suite.T().Fatalf("Failed to create invalid JSON file: %v", err)
	}

	oldGetProjectRoot := helpers.GetProjectRoot
	helpers.GetProjectRoot = func() (string, error) {
		return tempDir, nil
	}
	defer func() { helpers.GetProjectRoot = oldGetProjectRoot }()

	err = seeds.ListUserSeeder(suite.db)
	assert.Error(suite.T(), err, "Should return error for invalid JSON content")
}

// TestSeedsTestSuite runs the SeedsTestSuite to test seed-related database operations using the testify suite framework.
func TestSeedsTestSuite(t *testing.T) {
	suite.Run(t, new(SeedsTestSuite))
}

// createTestJSONFile creates a JSON file at the specified path and encodes the provided SeedUserRequest data into it.
// It uses pretty formatting with indentation and returns an error if file creation or encoding fails.
func createTestJSONFile(path string, data []SeedUserRequest) error {
	file, err := os.Create(path)
	if err != nil {
		return err
	}
	defer func(file *os.File) {
		err := file.Close()
		if err != nil {
			panic(err)
		}
	}(file)

	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "  ")
	return encoder.Encode(data)
}

// copyFile copies a file from the source path `src` to the destination path `dst`.
// It returns an error if the source file does not exist, cannot be read, or if the write operation fails.
func copyFile(src, dst string) error {
	input, err := os.ReadFile(src)
	if err != nil {
		return err
	}

	return os.WriteFile(dst, input, 0644)
}
