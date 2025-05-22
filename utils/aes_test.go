package utils

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestAESEncryptDecrypt tests the AES encryption and decryption functions to ensure encrypted data can be correctly decrypted.
func TestAESEncryptDecrypt(t *testing.T) {
	tests := []struct {
		name      string
		input     string
		wantError bool
	}{
		{
			name:      "empty string",
			input:     "",
			wantError: false,
		},
		{
			name:      "simple text",
			input:     "hello world",
			wantError: false,
		},
		{
			name:      "special characters",
			input:     "p@ssw0rd!$%^&*()",
			wantError: false,
		},
		{
			name:      "long text",
			input:     "Lorem ipsum dolor sit amet, consectetur adipiscing elit. Nullam auctor, nisl eget ultricies tincidunt, nisl nisl aliquam nisl, eget ultricies nisl nisl eget nisl.",
			wantError: false,
		},
	}

	for _, tt := range tests {
		t.Run(
			tt.name, func(t *testing.T) {
				encrypted, err := AESEncrypt(tt.input)
				if tt.wantError {
					assert.Error(t, err)
					return
				}
				require.NoError(t, err)
				assert.NotEmpty(t, encrypted)

				decrypted, err := AESDecrypt(encrypted)
				require.NoError(t, err)
				assert.Equal(t, tt.input, decrypted)
			},
		)
	}
}

func TestAESDecryptInvalidInput(t *testing.T) {
	tests := []struct {
		name        string
		input       string
		wantError   bool
		errorString string
	}{
		{
			name:        "invalid hex string",
			input:       "not a hex string",
			wantError:   true,
			errorString: "error in decoding encrypted string",
		},
		{
			name:        "short ciphertext",
			input:       "aabbcc",
			wantError:   true,
			errorString: "ciphertext too short",
		},
		{
			name:        "malformed ciphertext",
			input:       "aabbccddeeff00112233445566778899aabbccddeeff",
			wantError:   true,
			errorString: "error decrypting",
		},
	}

	for _, tt := range tests {
		t.Run(
			tt.name, func(t *testing.T) {
				decrypted, err := AESDecrypt(tt.input)
				if tt.wantError {
					assert.Error(t, err)
					assert.Contains(t, err.Error(), tt.errorString)
					assert.Empty(t, decrypted)
				} else {
					assert.NoError(t, err)
				}
			},
		)
	}
}

// TestAESEncryptConsistency ensures the consistency of AES encryption and decryption by verifying input-output correctness.
// Verifies that encryption produces different ciphertexts for the same input due to nonce randomness.
// Ensures decrypting the ciphertexts returns the original plaintext input without error.
func TestAESEncryptConsistency(t *testing.T) {
	input := "test input"
	encrypted1, err1 := AESEncrypt(input)
	require.NoError(t, err1)
	assert.NotEmpty(t, encrypted1)

	encrypted2, err2 := AESEncrypt(input)
	require.NoError(t, err2)
	assert.NotEmpty(t, encrypted2)

	assert.NotEqual(t, encrypted1, encrypted2)

	decrypted1, err := AESDecrypt(encrypted1)
	require.NoError(t, err)
	assert.Equal(t, input, decrypted1)

	decrypted2, err := AESDecrypt(encrypted2)
	require.NoError(t, err)
	assert.Equal(t, input, decrypted2)
}

// TestAESDecryptPanicRecovery validates the recovery mechanism of AESDecrypt when a panic occurs during decryption.
func TestAESDecryptPanicRecovery(t *testing.T) {
	decrypted, err := AESDecrypt("00112233445566778899aabb00112233445566778899aabb")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "error decrypting")
	assert.Empty(t, decrypted)
}
