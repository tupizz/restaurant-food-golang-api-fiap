package dto

type ClientInput struct {
	Name string `json:"name" binding:"required,max=100"`
	CPF  string `json:"cpf" binding:"required"`
}
