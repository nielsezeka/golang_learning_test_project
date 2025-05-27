package utils

// APIError represents a standard API error response
// swagger:model
type APIError struct {
	Error string `json:"error"`
}
