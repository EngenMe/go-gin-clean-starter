package config

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/Caknoooo/go-gin-clean-starter/constants"
)

// TestGetEnv verifies the behavior of getEnv function under different conditions using test cases.
func TestGetEnv(t *testing.T) {
	t.Run(
		"Should return value when env exists", func(t *testing.T) {
			key := "TEST_KEY"
			value := "test_value"
			os.Setenv(key, value)
			defer os.Unsetenv(key)

			result := getEnv(key, "default")
			assert.Equal(t, value, result)
		},
	)

	t.Run(
		"Should return default when env doesn't exist", func(t *testing.T) {
			result := getEnv("NON_EXISTENT_KEY", "default_value")
			assert.Equal(t, "default_value", result)
		},
	)
}

// TestLoadEnv verifies the behavior of the loadEnv function under different scenarios in a test environment.
// It checks if environment variables are correctly loaded from a .env.test file or if the function panics when the file is missing.
func TestLoadEnv(t *testing.T) {
	t.Run(
		"Should load .env.test file in testing mode", func(t *testing.T) {
			originalEnv := os.Getenv("APP_ENV")
			defer func() {
				os.Setenv("APP_ENV", originalEnv)
			}()

			content := "DB_USER=testuser\nDB_PASS=testpass\nDB_NAME=testdb\n"
			err := os.WriteFile(".env.test", []byte(content), 0644)
			require.NoError(t, err)
			defer os.Remove(".env.test")

			loadEnv()

			assert.Equal(t, "testuser", os.Getenv("DB_USER"))
			assert.Equal(t, "testpass", os.Getenv("DB_PASS"))
			assert.Equal(t, "testdb", os.Getenv("DB_NAME"))
		},
	)

	t.Run(
		"Should panic when .env.test file is missing in testing mode", func(t *testing.T) {
			originalEnv := os.Getenv("APP_ENV")
			os.Setenv("APP_ENV", constants.ENUM_RUN_TESTING)
			defer func() {
				os.Setenv("APP_ENV", originalEnv)
				if r := recover(); r == nil {
					t.Errorf("The code did not panic")
				}
			}()

			os.Remove(".env.test")

			loadEnv()
		},
	)
}
