package dto

import "time"

type ClientOutput struct {
	ID        int       `json:"id"`
	Name      string    `json:"name"`
	CPF       string    `json:"cpf"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
