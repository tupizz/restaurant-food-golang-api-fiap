package dto

import "github.com/go-playground/validator/v10"

type ProductInputUpdate struct {
	ID          int      `json:"id" validate:"required"`
	Name        string   `json:"name" validate:"omitempty,min=2"`
	Price       float64  `json:"price" validate:"omitempty,gte=0"`
	Description string   `json:"description" validate:"omitempty,min=10"`
	Category    string   `json:"category" validate:"omitempty,min=3"`
	Images      []string `json:"images" validate:"omitempty,dive,url"`
}

type ProductInputCreate struct {
	Name        string   `json:"name" validate:"required,min=2"`
	Price       float64  `json:"price" validate:"required,gte=0"`
	Description string   `json:"description" validate:"required,min=10"`
	Category    string   `json:"category" validate:"required,min=3"`
	Images      []string `json:"images" validate:"required,dive,url"`
}

var validate *validator.Validate

func init() {
	// Initialize the validator
	validate = validator.New()
}

// ValidateProductUpdate validates the ProductInputUpdate struct
func ValidateProductUpdate(input ProductInputUpdate) error {
	return validate.Struct(input)
}

// ValidateProductCreate validates the ProductInputCreate struct
func ValidateProductCreate(input ProductInputCreate) error {
	return validate.Struct(input)
}
