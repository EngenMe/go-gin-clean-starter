package utils_test

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
	"gorm.io/gorm"

	"github.com/Caknoooo/go-gin-clean-starter/tests/integration/container"
	"github.com/Caknoooo/go-gin-clean-starter/utils"
)

// EmailIntegrationTestSuite is a test suite for verifying email integration using SMTP and a test database environment.
// It sets up and tears down the required containers for database and SMTP server during the test lifecycle.
type EmailIntegrationTestSuite struct {
	suite.Suite
	smtpContainer testcontainers.Container
	dbContainer   *container.TestDatabaseContainer
	db            *gorm.DB
}

// SetupSuite initializes and configures the test environment, including database and SMTP containers with necessary settings.
func (suite *EmailIntegrationTestSuite) SetupSuite() {
	container.LoadTestEnv()

	ctx := context.Background()

	dbContainer, err := container.StartTestContainer()
	require.NoError(suite.T(), err)
	suite.dbContainer = dbContainer

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

	suite.db = container.SetUpDatabaseConnection()

	smtpReq := testcontainers.ContainerRequest{
		Image:        "mailhog/mailhog",
		ExposedPorts: []string{"1025/tcp", "8025/tcp"},
		WaitingFor:   wait.ForListeningPort("1025/tcp"),
	}
	smtpContainer, err := testcontainers.GenericContainer(
		ctx, testcontainers.GenericContainerRequest{
			ContainerRequest: smtpReq,
			Started:          true,
		},
	)
	require.NoError(suite.T(), err)
	suite.smtpContainer = smtpContainer

	smtpHost, err := smtpContainer.Host(ctx)
	require.NoError(suite.T(), err)

	smtpPort, err := smtpContainer.MappedPort(ctx, "1025")
	require.NoError(suite.T(), err)

	envVars = map[string]string{
		"SMTP_HOST":          smtpHost,
		"SMTP_PORT":          smtpPort.Port(),
		"SMTP_SENDER_NAME":   container.GetEnvWithDefault("SMTP_SENDER_NAME", "Test Sender"),
		"SMTP_AUTH_EMAIL":    container.GetEnvWithDefault("SMTP_AUTH_EMAIL", "test@example.com"),
		"SMTP_AUTH_PASSWORD": container.GetEnvWithDefault("SMTP_AUTH_PASSWORD", "password123"),
	}
	if err := container.SetEnv(envVars); err != nil {
		panic(fmt.Sprintf("Failed to set env vars: %v", err))
	}
}

// TearDownSuite cleans up the test environment by unsetting environment variables, closing database connections, and stopping containers.
func (suite *EmailIntegrationTestSuite) TearDownSuite() {
	ctx := context.Background()
	timeout := 10 * time.Second

	for _, env := range []string{
		"DB_HOST", "DB_PORT", "DB_USER", "DB_PASS", "DB_NAME",
		"SMTP_HOST", "SMTP_PORT", "SMTP_AUTH_EMAIL", "SMTP_AUTH_PASSWORD", "SMTP_SENDER_NAME",
	} {
		err := os.Unsetenv(env)
		if err != nil {
			panic(err)
		}
	}

	if suite.db != nil {
		err := container.CloseDatabaseConnection(suite.db)
		assert.NoError(suite.T(), err)
	}

	if suite.smtpContainer != nil {
		_ = suite.smtpContainer.Stop(ctx, &timeout)
	}
	if suite.dbContainer != nil {
		_ = suite.dbContainer.Stop()
	}
}

// TestSendMail_Integration validates the behavior of the SendMail function using integration tests with different scenarios.
func (suite *EmailIntegrationTestSuite) TestSendMail_Integration() {
	tests := []struct {
		name      string
		toEmail   string
		subject   string
		body      string
		wantError bool
	}{
		{
			name:      "Successfully send email",
			toEmail:   "recipient@example.com",
			subject:   "Test Subject",
			body:      "<p>Test Body</p>",
			wantError: false,
		},
		{
			name:      "Invalid recipient email",
			toEmail:   "",
			subject:   "Test Subject",
			body:      "<p>Test Body</p>",
			wantError: true,
		},
	}

	for _, tt := range tests {
		suite.Run(
			tt.name, func() {
				err := utils.SendMail(tt.toEmail, tt.subject, tt.body)
				if tt.wantError {
					assert.Error(suite.T(), err)
				} else {
					assert.NoError(suite.T(), err)
				}
			},
		)
	}
}
