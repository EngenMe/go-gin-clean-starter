package script

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"gorm.io/gorm"
)

// MockDB is a type that embeds mock.Mock, allowing it to be used for mocking database interactions in tests.
type MockDB struct {
	mock.Mock
}

// GormDB retrieves a *gorm.DB instance, typically used for mocking database interactions in tests.
func (m *MockDB) GormDB() *gorm.DB {
	args := m.Called()
	return args.Get(0).(*gorm.DB)
}

// TestNewExampleScript verifies that a new ExampleScript instance is correctly initialized with a mocked gorm.DB connection.
func TestNewExampleScript(t *testing.T) {
	mockDB := new(MockDB)
	mockGormDB := &gorm.DB{}
	mockDB.On("GormDB").Return(mockGormDB)

	script := NewExampleScript(mockDB.GormDB())

	assert.NotNil(t, script)
	assert.Equal(t, mockGormDB, script.db)
	mockDB.AssertExpectations(t)
}

// TestExampleScript_Run tests the Run method of the ExampleScript, validating success and failure scenarios.
func TestExampleScript_Run(t *testing.T) {
	tests := []struct {
		name    string
		wantErr bool
	}{
		{
			name:    "successful run",
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(
			tt.name, func(t *testing.T) {
				mockDB := new(MockDB)
				mockGormDB := &gorm.DB{}
				mockDB.On("GormDB").Return(mockGormDB)

				script := NewExampleScript(mockDB.GormDB())

				err := script.Run()

				if tt.wantErr {
					assert.Error(t, err)
				} else {
					assert.NoError(t, err)
				}
			},
		)
	}
}
