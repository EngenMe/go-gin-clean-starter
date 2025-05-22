package utils

// Response represents a standard structure for API responses.
// Status indicates the success or failure of the operation.
// Message provides a human-readable message about the operation.
// Error holds error details when the operation fails (optional).
// Data contains the result of the operation, if any (optional).
// Meta includes additional metadata related to the response (optional).
type Response struct {
	Status  bool   `json:"status"`
	Message string `json:"message"`
	Error   any    `json:"error,omitempty"`
	Data    any    `json:"data,omitempty"`
	Meta    any    `json:"meta,omitempty"`
}

// BuildResponseSuccess constructs a successful Response with the provided message and data.
func BuildResponseSuccess(message string, data any) Response {
	res := Response{
		Status:  true,
		Message: message,
		Data:    data,
	}
	return res
}

// BuildResponseFailed constructs and returns a failed Response with given message, error, and data.
func BuildResponseFailed(message string, err string, data any) Response {
	res := Response{
		Status:  false,
		Message: message,
		Error:   err,
		Data:    data,
	}
	return res
}
