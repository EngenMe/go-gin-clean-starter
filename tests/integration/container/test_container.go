package container

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// TestDatabaseContainer wraps a testcontainers.Container with additional fields for Host and Port configuration.
type TestDatabaseContainer struct {
	testcontainers.Container
	Host string
	Port string
}

// StartTestContainer creates and starts a PostgreSQL test container and returns a TestDatabaseContainer with connection info.
// Returns an error if the container fails to start or retrieve connection details.
func StartTestContainer() (*TestDatabaseContainer, error) {
	ctx := context.Background()

	req := testcontainers.ContainerRequest{
		Image:        "postgres:13-alpine",
		ExposedPorts: []string{"5432/tcp"},
		Env: map[string]string{
			"POSTGRES_USER":     GetEnvWithDefault("DB_USER", "testuser"),
			"POSTGRES_PASSWORD": GetEnvWithDefault("DB_PASS", "testpassword"),
			"POSTGRES_DB":       GetEnvWithDefault("DB_NAME", "testdb"),
		},
		WaitingFor: wait.ForLog("database system is ready to accept connections").
			WithOccurrence(2).
			WithStartupTimeout(60 * time.Second),
	}

	container, err := testcontainers.GenericContainer(
		ctx, testcontainers.GenericContainerRequest{
			ContainerRequest: req,
			Started:          true,
		},
	)
	if err != nil {
		return nil, fmt.Errorf("failed to start container: %w", err)
	}

	mappedPort, err := container.MappedPort(ctx, "5432")
	if err != nil {
		return nil, fmt.Errorf("failed to get container port: %w", err)
	}

	host, err := container.Host(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get container host: %w", err)
	}

	return &TestDatabaseContainer{
		Container: container,
		Host:      host,
		Port:      mappedPort.Port(),
	}, nil
}

// Stop stops the running container with a specified timeout and releases its resources. It returns an error if unable to stop.
func (c *TestDatabaseContainer) Stop() error {
	ctx := context.Background()
	timeout := 10 * time.Second
	return c.Container.Stop(ctx, &timeout)
}

// CloseDatabaseConnection safely closes the underlying database connection from the provided *gorm.DB instance.
func CloseDatabaseConnection(db *gorm.DB) error {
	dbSQL, err := db.DB()
	if err != nil {
		return fmt.Errorf("failed to get underlying DB: %w", err)
	}
	return dbSQL.Close()
}

// SetUpDatabaseConnection initializes and returns a GORM database connection using environment variables for configuration.
// It enables the `uuid-ossp` PostgreSQL extension if not already enabled. Panics on connection or setup failure.
func SetUpDatabaseConnection() *gorm.DB {
	dsn := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%s sslmode=disable",
		os.Getenv("DB_HOST"),
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASS"),
		os.Getenv("DB_NAME"),
		os.Getenv("DB_PORT"),
	)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		panic(fmt.Errorf("failed to connect to database: %w", err))
	}

	err = db.Exec("CREATE EXTENSION IF NOT EXISTS \"uuid-ossp\"").Error
	if err != nil {
		panic(fmt.Errorf("failed to enable uuid-ossp extension: %w", err))
	}

	return db
}
