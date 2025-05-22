package utils

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"gopkg.in/gomail.v2"

	"github.com/Caknoooo/go-gin-clean-starter/config"
)

// MockDialer is a mock implementation of the Dialer interface, used for testing email sending functionality.
type MockDialer struct {
	mock.Mock
}

// DialAndSend simulates sending email messages and returns any error captured during the mock execution.
func (m *MockDialer) DialAndSend(messages ...*gomail.Message) error {
	args := m.Called(messages)
	return args.Error(0)
}

// TestSendMail tests the SendMail function by verifying email sending scenarios with mocked configurations and errors.
func TestSendMail(t *testing.T) {
	originalNewEmailConfig := config.NewEmailConfig
	defer func() { config.NewEmailConfig = originalNewEmailConfig }()

	tests := []struct {
		name           string
		emailConfig    *config.EmailConfig
		emailConfigErr error
		dialerErr      error
		wantErr        bool
	}{
		{
			name: "Successfully send email",
			emailConfig: &config.EmailConfig{
				Host:         "smtp.example.com",
				Port:         587,
				AuthEmail:    "test@example.com",
				AuthPassword: "password",
			},
			wantErr: false,
		},
		{
			name:           "Failed to get email config",
			emailConfig:    nil,
			emailConfigErr: errors.New("config error"),
			wantErr:        true,
		},
		{
			name: "Failed to send email",
			emailConfig: &config.EmailConfig{
				Host:         "smtp.example.com",
				Port:         587,
				AuthEmail:    "test@example.com",
				AuthPassword: "password",
			},
			dialerErr: errors.New("smtp error"),
			wantErr:   true,
		},
	}

	for _, tt := range tests {
		t.Run(
			tt.name, func(t *testing.T) {
				config.NewEmailConfig = func() (*config.EmailConfig, error) {
					return tt.emailConfig, tt.emailConfigErr
				}

				mockDialer := new(MockDialer)
				if tt.emailConfig != nil {
					mockDialer.On("DialAndSend", mock.Anything).Return(tt.dialerErr)
				}

				originalNewDialer := NewDialer
				defer func() { NewDialer = originalNewDialer }()
				NewDialer = func(host string, port int, username, password string) Dialer {
					return mockDialer
				}

				err := SendMail("recipient@example.com", "Test Subject", "<p>Test Body</p>")

				if tt.wantErr {
					assert.Error(t, err)
				} else {
					assert.NoError(t, err)
				}

				if tt.emailConfig != nil {
					mockDialer.AssertExpectations(t)
				}
			},
		)
	}
}
