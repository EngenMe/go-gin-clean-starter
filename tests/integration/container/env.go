package container

import (
	"github.com/Caknoooo/go-gin-clean-starter/helpers"
	"github.com/joho/godotenv"
	"log"
	"os"
	"path"
	"sync"
)

var (
	envMutex sync.Mutex
)

// SetEnv sets multiple environment variables provided in the map and returns an error if setting any variable fails.
func SetEnv(vars map[string]string) error {
	envMutex.Lock()
	defer envMutex.Unlock()

	for key, value := range vars {
		if err := os.Setenv(key, value); err != nil {
			return err
		}
	}
	return nil
}

// UnsetEnv removes the specified environment variables by their keys, ensuring thread-safe operation with a mutex.
func UnsetEnv(keys []string) error {
	envMutex.Lock()
	defer envMutex.Unlock()

	for _, key := range keys {
		if err := os.Unsetenv(key); err != nil {
			return err
		}
	}
	return nil
}

// LoadTestEnv initializes the test environment by loading variables from the .env.test file in the project root directory.
// Logs a warning if the .env.test file cannot be loaded, but fails if required environment variables are missing.
func LoadTestEnv() {
	projectRoot, err := helpers.GetProjectRoot()
	if err != nil {
		log.Fatalf("Failed to get project root: %v", err)
	}

	envPath := path.Join(projectRoot, ".env.test")
	if err := godotenv.Overload(envPath); err != nil {
		log.Printf("Warning: Failed to load .env.test file: %v", err)
	}

	validateRequiredVars()
}

// validateRequiredVars ensures that all required environment variables are set, logging a fatal error if any are missing.
func validateRequiredVars() {
	requiredVars := []string{
		"APP_ENV",
		"JWT_SECRET",
		"DB_HOST", "DB_USER", "DB_PASS", "DB_NAME", "DB_PORT",
		"SMTP_HOST", "SMTP_PORT", "SMTP_AUTH_EMAIL",
	}

	for _, envVar := range requiredVars {
		if os.Getenv(envVar) == "" {
			log.Fatalf("Missing required environment variable: %s", envVar)
		}
	}
}

// GetEnvWithDefault retrieves the value of the environment variable by key or returns the defaultValue if not set.
func GetEnvWithDefault(key, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultValue
}
