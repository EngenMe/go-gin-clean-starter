package migrations

import (
	"fmt"
	"testing"

	"github.com/Caknoooo/go-gin-clean-starter/entity"
	"github.com/Caknoooo/go-gin-clean-starter/migrations"
	"github.com/Caknoooo/go-gin-clean-starter/tests/integration/container"
)

// TestMigrate_Integration is an integration test that verifies the migration process for database schema and initial setup.
func TestMigrate_Integration(t *testing.T) {
	container.LoadTestEnv()

	dbContainer, err := container.StartTestContainer()
	if err != nil {
		t.Fatalf("Failed to start test container: %v", err)
	}
	defer func() {
		if err := dbContainer.Stop(); err != nil {
			t.Errorf("Failed to stop test container: %v", err)
		}
	}()

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
	defer func() {
		if err := container.CloseDatabaseConnection(db); err != nil {
			t.Errorf("Failed to close database connection: %v", err)
		}
	}()

	err = migrations.Migrate(db)
	if err != nil {
		t.Errorf("Migrate() returned an error: %v", err)
	}

	var tableCount int64
	db.Raw("SELECT COUNT(*) FROM information_schema.tables WHERE table_schema = 'public' AND table_name IN ('users', 'refresh_tokens')").Scan(&tableCount)
	if tableCount != 2 {
		t.Errorf("Expected 2 tables to be created, but found %d", tableCount)
	}

	var user entity.User
	var refreshToken entity.RefreshToken

	if err := db.Model(&user).Limit(1).Find(&user).Error; err != nil {
		t.Errorf("Failed to query User table: %v", err)
	}

	if err := db.Model(&refreshToken).Limit(1).Find(&refreshToken).Error; err != nil {
		t.Errorf("Failed to query RefreshToken table: %v", err)
	}
}
