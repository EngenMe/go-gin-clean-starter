package helpers

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

// TestHashPassword tests the functionality of the HashPassword function for various input scenarios and potential edge cases.
func TestHashPassword(t *testing.T) {
	tests := []struct {
		name        string
		password    string
		expectError bool
	}{
		{
			name:        "successful hash",
			password:    "securePassword123!",
			expectError: false,
		},
		{
			name:        "empty password",
			password:    "",
			expectError: false,
		},
		{
			name:        "password at bcrypt limit (72 bytes)",
			password:    "this-is-72-bytes-xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx",
			expectError: false,
		},
		{
			name:        "password exceeds bcrypt limit",
			password:    "this-is-a-very-long-password-with-more-than-72-bytes-which-is-bcrypts-maximum-password-length-xxxx",
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(
			tt.name, func(t *testing.T) {
				hash, err := HashPassword(tt.password)

				if tt.expectError {
					assert.Error(t, err)
					assert.Empty(t, hash)
				} else {
					assert.NoError(t, err)
					assert.NotEmpty(t, hash)
					assert.Greater(t, len(hash), 10)
				}
			},
		)
	}
}

// TestCheckPassword verifies the CheckPassword function logic with a variety of test cases, including edge cases and errors.
func TestCheckPassword(t *testing.T) {
	validHash, err := HashPassword("correctPassword")
	assert.NoError(t, err)

	emptyHash, err := HashPassword("")
	assert.NoError(t, err)

	maxLenHash, err := HashPassword("this-is-72-bytes-xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx")
	assert.NoError(t, err)

	tests := []struct {
		name        string
		hash        string
		password    string
		expected    bool
		expectError bool
	}{
		{
			name:        "correct password",
			hash:        validHash,
			password:    "correctPassword",
			expected:    true,
			expectError: false,
		},
		{
			name:        "incorrect password",
			hash:        validHash,
			password:    "wrongPassword",
			expected:    false,
			expectError: true,
		},
		{
			name:        "empty password with empty hash",
			hash:        emptyHash,
			password:    "",
			expected:    true,
			expectError: false,
		},
		{
			name:        "empty password with non-empty hash",
			hash:        validHash,
			password:    "",
			expected:    false,
			expectError: true,
		},
		{
			name:        "invalid hash format",
			hash:        "not-a-valid-hash",
			password:    "anypassword",
			expected:    false,
			expectError: true,
		},
		{
			name:        "password at bcrypt limit",
			hash:        maxLenHash,
			password:    "this-is-72-bytes-xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx",
			expected:    true,
			expectError: false,
		},
	}

	for _, tt := range tests {
		t.Run(
			tt.name, func(t *testing.T) {
				result, err := CheckPassword(tt.hash, []byte(tt.password))

				if tt.expectError {
					assert.Error(t, err)
				} else {
					assert.NoError(t, err)
				}

				assert.Equal(t, tt.expected, result)
			},
		)
	}
}

// TestHashAndCheckPasswordIntegration tests end-to-end password hashing and validation for various input scenarios.
func TestHashAndCheckPasswordIntegration(t *testing.T) {
	testCases := []struct {
		name     string
		password string
	}{
		{
			name:     "normal password",
			password: "integrationTestPassword123!",
		},
		{
			name:     "empty password",
			password: "",
		},
		{
			name:     "72-byte password",
			password: "this-is-72-bytes-xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx",
		},
	}

	for _, tc := range testCases {
		t.Run(
			tc.name, func(t *testing.T) {
				hash, err := HashPassword(tc.password)
				if len(tc.password) > 72 {
					assert.Error(t, err)
					return
				}

				assert.NoError(t, err)
				assert.NotEmpty(t, hash)

				result, err := CheckPassword(hash, []byte(tc.password))
				assert.NoError(t, err)
				assert.True(t, result)

				result, err = CheckPassword(hash, []byte("wrongPassword"))
				assert.Error(t, err)
				assert.False(t, result)
			},
		)
	}
}
