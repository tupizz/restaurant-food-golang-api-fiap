package mappers

import (
	"github.com/tupizz/restaurant-food-golang-api-fiap/internal/core/domain/entities"
	"github.com/tupizz/restaurant-food-golang-api-fiap/internal/core/usecase/dto"
)

func ToClientDTO(client entities.Client) dto.ClientOutput {
	return dto.ClientOutput{
		ID:        client.ID,
		Name:      client.Name,
		CPF:       client.CPF,
		CreatedAt: client.CreatedAt,
		UpdatedAt: client.UpdatedAt,
	}
}
