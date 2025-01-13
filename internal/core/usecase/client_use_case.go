package usecase

import (
	"context"
	"errors"

	"github.com/tupizz/restaurant-food-golang-api-fiap/internal/core/domain/entities"
	"github.com/tupizz/restaurant-food-golang-api-fiap/internal/core/domain/validator"
	"github.com/tupizz/restaurant-food-golang-api-fiap/internal/core/usecase/dto"
	"github.com/tupizz/restaurant-food-golang-api-fiap/internal/core/usecase/ports"
)

type ClientUseCase interface {
	CreateClient(ctx context.Context, input dto.ClientInput) (dto.ClientOutput, error)
	GetClientByCpf(ctx context.Context, cpf string) (dto.ClientOutput, error)
}

type clientUseCase struct {
	clientRepository ports.ClientRepository
}

func NewClientUseCase(clientRepository ports.ClientRepository) ClientUseCase {
	return &clientUseCase{clientRepository: clientRepository}
}

func (s *clientUseCase) CreateClient(ctx context.Context, input dto.ClientInput) (dto.ClientOutput, error) {
	if err := validator.IsValidCPF(input.CPF); err != nil {
		return dto.ClientOutput{}, errors.New("CPF inválido")
	}

	_, err := s.clientRepository.GetByCpf(ctx, input.CPF)
	if err == nil {
		return dto.ClientOutput{}, errors.New("CPF já cadastrado")
	}

	client := entities.Client{
		Name: input.Name,
		CPF:  input.CPF,
	}

	createdClient, err := s.clientRepository.Create(ctx, client)
	if err != nil {
		return dto.ClientOutput{}, err
	}

	output := dto.ClientOutput{
		ID:        createdClient.ID,
		Name:      createdClient.Name,
		CPF:       createdClient.CPF,
		CreatedAt: createdClient.CreatedAt,
		UpdatedAt: createdClient.UpdatedAt,
	}

	return output, nil
}

func (s *clientUseCase) GetClientByCpf(ctx context.Context, cpf string) (dto.ClientOutput, error) {
	client, err := s.clientRepository.GetByCpf(ctx, cpf)
	if err != nil {
		return dto.ClientOutput{}, err
	}

	output := dto.ClientOutput{
		ID:        client.ID,
		Name:      client.Name,
		CPF:       client.CPF,
		CreatedAt: client.CreatedAt,
		UpdatedAt: client.UpdatedAt,
	}

	return output, nil
}
