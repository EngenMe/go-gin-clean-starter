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
	suite.Suite
	emailContainer *container.TestDatabaseContainer
}

// SetupSuite initializes the test environment by starting a MailHog test container and setting necessary environment variables.
func (suite *EmailConfigTestSuite) SetupSuite() {
	testContainer, err := container.StartTestContainer()
	require.NoError(suite.T(), err)
	suite.emailContainer = testContainer

	err = os.Setenv("SMTP_HOST", testContainer.Host)
	if err != nil {
		panic(err)
	}
	err = os.Setenv("SMTP_PORT", testContainer.Port)
	if err != nil {
		panic(err)
	}
	err = os.Setenv("SMTP_SENDER_NAME", "Test Sender")
	if err != nil {
		panic(err)
	}
	err = os.Setenv("SMTP_AUTH_EMAIL", "test@example.com")
	if err != nil {
		panic(err)
	}
	err = os.Setenv("SMTP_AUTH_PASSWORD", "password123")
	if err != nil {
		panic(err)
	}
}

// TearDownSuite cleans up test resources by unsetting environment variables and stopping the email test container.
func (suite *EmailConfigTestSuite) TearDownSuite() {
	err := os.Unsetenv("SMTP_HOST")
	if err != nil {
		panic(err)
	}
	err = os.Unsetenv("SMTP_PORT")
	if err != nil {
		panic(err)
	}
	err = os.Unsetenv("SMTP_SENDER_NAME")
	if err != nil {
		panic(err)
	}
	err = os.Unsetenv("SMTP_AUTH_EMAIL")
	if err != nil {
		panic(err)
	}
	err = os.Unsetenv("SMTP_AUTH_PASSWORD")
	if err != nil {
		panic(err)
	}

	if suite.emailContainer != nil {
		err := suite.emailContainer.Stop()
		require.NoError(suite.T(), err)
	}
}

// TestNewEmailConfig_Integration validates that a new EmailConfig is created correctly using environment variables.
// Ensures no errors occur and verifies that configuration values match the expected environment variable settings.
func (suite *EmailConfigTestSuite) TestNewEmailConfig_Integration() {
	emailConfig, err := config.NewEmailConfig()
	require.NoError(suite.T(), err)
	assert.NotNil(suite.T(), emailConfig)

	assert.Equal(suite.T(), os.Getenv("SMTP_HOST"), emailConfig.Host)
	port, _ := strconv.Atoi(os.Getenv("SMTP_PORT"))
	assert.Equal(suite.T(), port, emailConfig.Port)
	assert.Equal(suite.T(), os.Getenv("SMTP_SENDER_NAME"), emailConfig.SenderName)
	assert.Equal(suite.T(), os.Getenv("SMTP_AUTH_EMAIL"), emailConfig.AuthEmail)
	assert.Equal(suite.T(), os.Getenv("SMTP_AUTH_PASSWORD"), emailConfig.AuthPassword)
}

// TestEmailConfigTestSuite executes the test suite for email configuration using the EmailConfigTestSuite.
func TestEmailConfigTestSuite(t *testing.T) {
	suite.Run(t, new(EmailConfigTestSuite))
}
