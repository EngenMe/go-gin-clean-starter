package dto

const (
	// MESSAGE_SUCCESS_REFRESH_TOKEN represents a success message for a successful token refresh operation.
	MESSAGE_SUCCESS_REFRESH_TOKEN = "Successfully refreshed token"

	// MESSAGE_FAILED_REFRESH_TOKEN indicates a failure message when attempting to refresh an authentication token unsuccessfully.
	MESSAGE_FAILED_REFRESH_TOKEN = "Failed to refresh token"

	// MESSAGE_FAILED_INVALID_REFRESH_TOKEN indicates an error message for an invalid refresh token during token refresh operations.
	MESSAGE_FAILED_INVALID_REFRESH_TOKEN = "Invalid refresh token"

	// MESSAGE_FAILED_EXPIRED_REFRESH_TOKEN indicates a failure message when the refresh token has expired and is no longer valid.
	MESSAGE_FAILED_EXPIRED_REFRESH_TOKEN = "Refresh token has expired"
)

// TokenResponse represents a response containing access and refresh tokens along with the associated user role.
type TokenResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	Role         string `json:"role"`
}

// RefreshTokenRequest represents a data structure for handling refresh token requests.
// It contains the required refresh token field, validated as non-empty.
type RefreshTokenRequest struct {
	RefreshToken string `json:"refresh_token" binding:"required"`
}
