package config_test

import (
	"os"
	"strconv"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"

	"github.com/Caknoooo/go-gin-clean-starter/config"
	"github.com/Caknoooo/go-gin-clean-starter/tests/integration/container"
)

// EmailConfigTestSuite is a test suite for verifying email configuration using environment variables and test containers.
type EmailConfigTestSuite struct {
	container.BaseSuite
	emailContainer *container.TestDatabaseContainer
}

// SetupSuite initializes the email configuration test suite by starting a test container and setting environment variables.
func (s *EmailConfigTestSuite) SetupSuite() {
	s.BaseSuite.SetupSuite()

	dbContainer, err := container.StartTestContainer()
	require.NoError(s.T(), err)
	s.emailContainer = dbContainer

	envVars := map[string]string{
		"SMTP_HOST":          dbContainer.Host,
		"SMTP_PORT":          dbContainer.Port,
		"SMTP_SENDER_NAME":   container.GetEnvWithDefault("SMTP_SENDER_NAME", "Test Sender"),
		"SMTP_AUTH_EMAIL":    container.GetEnvWithDefault("SMTP_AUTH_EMAIL", "test@example.com"),
		"SMTP_AUTH_PASSWORD": container.GetEnvWithDefault("SMTP_AUTH_PASSWORD", "password123"),
	}
	s.SetupEnv(envVars)
}

// TearDownSuite cleans up the test suite by unsetting environment variables and stopping the email test container.
func (s *EmailConfigTestSuite) TearDownSuite() {
	s.CleanupEnv(
		[]string{
			"SMTP_HOST",
			"SMTP_PORT",
			"SMTP_SENDER_NAME",
			"SMTP_AUTH_EMAIL",
			"SMTP_AUTH_PASSWORD",
		},
	)

	if s.emailContainer != nil {
		require.NoError(s.T(), s.emailContainer.Stop())
	}
}

// TestNewEmailConfig_Integration validates that a new EmailConfig is created correctly using environment variables.
// Ensures no errors occur and verifies that configuration values match the expected environment variable settings.
func (s *EmailConfigTestSuite) TestNewEmailConfig_Integration() {
	emailConfig, err := config.NewEmailConfig()
	require.NoError(s.T(), err)
	assert.NotNil(s.T(), emailConfig)

	assert.Equal(s.T(), os.Getenv("SMTP_HOST"), emailConfig.Host)
	port, _ := strconv.Atoi(os.Getenv("SMTP_PORT"))
	assert.Equal(s.T(), port, emailConfig.Port)
	assert.Equal(s.T(), os.Getenv("SMTP_SENDER_NAME"), emailConfig.SenderName)
	assert.Equal(s.T(), os.Getenv("SMTP_AUTH_EMAIL"), emailConfig.AuthEmail)
	assert.Equal(s.T(), os.Getenv("SMTP_AUTH_PASSWORD"), emailConfig.AuthPassword)
}

// TestEmailConfigTestSuite executes the test suite for email configuration using the EmailConfigTestSuite.
func TestEmailConfigTestSuite(t *testing.T) {
	suite.Run(t, new(EmailConfigTestSuite))
}
