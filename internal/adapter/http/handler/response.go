package handler

// ErrorResponse represents a standard error response.
// It includes an error message and an optional error code.
type ErrorResponse struct {
	Error   string `json:"error" example:"Invalid request"`
	Code    int    `json:"code,omitempty" example:"400"`
	Details string `json:"details,omitempty" example:"Additional information"`
}
