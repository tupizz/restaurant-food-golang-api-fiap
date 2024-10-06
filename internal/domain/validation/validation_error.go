package validation

import (
	"github.com/go-playground/validator/v10"
)

type ErrorResponse struct {
	Field   string `json:"field"`
	Message string `json:"message"`
}

// getErrorMessage provides custom error messages based on the validation tag
func getErrorMessage(fieldError validator.FieldError) string {
	switch fieldError.Tag() {
	case "required":
		return "This field is required"
	case "min":
		return "The value is too short"
	case "gte":
		return "The value must be greater than or equal to 0"
	case "url":
		return "This must be a valid URL"
	default:
		return "Invalid value"
	}
}

// HandleValidationError processes validation errors and returns a list of field-specific error messages
func HandleValidationError(err error) []ErrorResponse {
	var errors []ErrorResponse

	if validationErrors, ok := err.(validator.ValidationErrors); ok {
		for _, fieldError := range validationErrors {
			errorResponse := ErrorResponse{
				Field:   fieldError.Field(),
				Message: getErrorMessage(fieldError), // Customize message based on the error tag
			}
			errors = append(errors, errorResponse)
		}
	}

	return errors
}
