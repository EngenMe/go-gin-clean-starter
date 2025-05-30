package entity_test

import (
	"fmt"
	"os"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"

	"github.com/Caknoooo/go-gin-clean-starter/entity"
	"github.com/Caknoooo/go-gin-clean-starter/tests/integration/container"
)

// TestMain sets up the environment for integration tests by initializing a test container and configuring database variables.
// It ensures proper cleanup by stopping the test container after tests have run.
func TestMain(m *testing.M) {
	container.LoadTestEnv()

	dbContainer, err := container.StartTestContainer()
	if err != nil {
		panic(err)
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

	code := m.Run()

	if err := dbContainer.Stop(); err != nil {
		panic(err)
	}

	os.Exit(code)
}

// setupTestDB initializes a test database connection, migrates the User schema, and returns the database instance.
// It fails the test if the database migration encounters an error.
func setupTestDB(t *testing.T) *gorm.DB {
	db := container.SetUpDatabaseConnection()

	err := db.AutoMigrate(&entity.User{})
	if err != nil {
		t.Fatalf("Failed to migrate database: %v", err)
	}

	return db
}

// cleanupTestDB cleans up the test database by dropping the User table and closing the database connection.
func cleanupTestDB(t *testing.T, db *gorm.DB) {
	err := db.Migrator().DropTable(&entity.User{})
	if err != nil {
		t.Fatalf("Failed to drop table: %v", err)
	}

	if err := container.CloseDatabaseConnection(db); err != nil {
		t.Fatalf("Failed to close database connection: %v", err)
	}
}

// TestUser_Integration_Create tests the integration for creating a user, validating correct behavior and handling edge cases.
func TestUser_Integration_Create(t *testing.T) {
	db := setupTestDB(t)
	defer cleanupTestDB(t, db)

	tests := []struct {
		name        string
		user        *entity.User
		expectError bool
		validate    func(t *testing.T, user *entity.User, db *gorm.DB)
	}{
		{
			name: "Valid user creation",
			user: &entity.User{
				Name:        "John Doe",
				Email:       "john-doe@example.com",
				PhoneNumber: "1234567890",
				Password:    "password123",
				Role:        "user",
				ImageUrl:    "https://example.com/image.jpg",
			},
			expectError: false,
			validate: func(t *testing.T, user *entity.User, db *gorm.DB) {
				var savedUser entity.User
				err := db.Where("email = ?", user.Email).First(&savedUser).Error
				assert.NoError(t, err, "User should exist in the database")
				assert.NotEqual(t, uuid.Nil, savedUser.ID, "ID should be generated")
				assert.Equal(t, user.Name, savedUser.Name, "Name should match")
				assert.Equal(t, user.Email, savedUser.Email, "Email should match")
				assert.Equal(t, user.PhoneNumber, savedUser.PhoneNumber, "PhoneNumber should match")
				assert.NotEqual(t, "password123", savedUser.Password, "Password should be hashed")
				assert.Equal(t, "user", savedUser.Role, "Role should be user")
				assert.False(t, savedUser.IsVerified, "IsVerified should be false")
			},
		},
		{
			name: "Duplicate email",
			user: &entity.User{
				Name:        "Jane Doe",
				Email:       "john@example.com",
				PhoneNumber: "0987654321",
				Password:    "password123",
				Role:        "user",
			},
			expectError: true,
			validate: func(t *testing.T, user *entity.User, db *gorm.DB) {
				var count int64
				db.Model(&entity.User{}).Where("email = ?", user.Email).Count(&count)
				assert.Equal(t, int64(1), count, "Only one user with this email should exist")
			},
		},
		{
			name: "Invalid role",
			user: &entity.User{
				Name:     "Invalid User",
				Email:    "invalid@example.com",
				Password: "password123",
				Role:     "invalid_role",
			},
			expectError: true,
		},
	}

	db.Create(
		&entity.User{
			Name:        "Existing User",
			Email:       "john@example.com",
			PhoneNumber: "1234567890",
			Password:    "password123",
			Role:        "user",
		},
	)

	for _, tt := range tests {
		t.Run(
			tt.name, func(t *testing.T) {
				err := db.Create(tt.user).Error
				if tt.expectError {
					assert.Error(t, err, "Expected an error")
				} else {
					assert.NoError(t, err, "Expected no error")
					tt.validate(t, tt.user, db)
				}
			},
		)
	}
}

// TestUser_Integration_Update tests the integration of updating a User entity in the database, validating correct functionality.
func TestUser_Integration_Update(t *testing.T) {
	db := setupTestDB(t)
	defer cleanupTestDB(t, db)

	user := &entity.User{
		Name:        "John Doe",
		Email:       "john@example.com",
		PhoneNumber: "1234567890",
		Password:    "password123",
		Role:        "user",
	}
	err := db.Create(user).Error
	assert.NoError(t, err, "Failed to create test user")

	tests := []struct {
		name        string
		update      func(user *entity.User)
		expectError bool
		validate    func(t *testing.T, user *entity.User, db *gorm.DB)
	}{
		{
			name: "Update password and name",
			update: func(user *entity.User) {
				user.Name = "John Updated"
				user.Password = "newpassword123"
			},
			expectError: false,
			validate: func(t *testing.T, user *entity.User, db *gorm.DB) {
				var updatedUser entity.User
				err := db.Where("email = ?", user.Email).First(&updatedUser).Error
				assert.NoError(t, err, "User should exist in the database")
				assert.Equal(t, "John Updated", updatedUser.Name, "Name should be updated")
				assert.NotEqual(t, "newpassword123", updatedUser.Password, "Password should be hashed")
			},
		},
		{
			name: "Update without password change",
			update: func(user *entity.User) {
				user.PhoneNumber = "0987654321"
				user.Role = "admin"
			},
			expectError: false,
			validate: func(t *testing.T, user *entity.User, db *gorm.DB) {
				var updatedUser entity.User
				err := db.Where("email = ?", user.Email).First(&updatedUser).Error
				assert.NoError(t, err, "User should exist in the database")
				assert.Equal(t, "0987654321", updatedUser.PhoneNumber, "PhoneNumber should be updated")
				assert.Equal(t, "admin", updatedUser.Role, "Role should be updated")
			},
		},
	}

	for _, tt := range tests {
		t.Run(
			tt.name, func(t *testing.T) {
				tt.update(user)
				err := db.Save(user).Error
				if tt.expectError {
					assert.Error(t, err, "Expected an error")
				} else {
					assert.NoError(t, err, "Expected no error")
					tt.validate(t, user, db)
				}
			},
		)
	}
}
