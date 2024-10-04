package dto

type UserInput struct {
	Name  string `json:"name" binding:"required" validate:"max=100"`
	Email string `json:"email" binding:"required" validate:"max=100,email"`
	Age   int    `json:"age" binding:"required,gte=0,lte=130" validate:"gte=0,lte=130"`
}
