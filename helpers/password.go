package helpers

import (
	"golang.org/x/crypto/bcrypt"
)

// HashPassword hashes a given plaintext password using bcrypt with a predefined cost and returns the hashed password.
func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 4)
	return string(bytes), err
}

// CheckPassword compares a hashed password with a plaintext password and returns whether they match or an error.
func CheckPassword(hashPassword string, plainPassword []byte) (bool, error) {
	hashPW := []byte(hashPassword)
	if err := bcrypt.CompareHashAndPassword(hashPW, plainPassword); err != nil {
		return false, err
	}
	return true, nil
}
