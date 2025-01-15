package usecase

import (
	"context"
	"errors"

	"github.com/tupizz/restaurant-food-golang-api-fiap/internal/core/domain/entities"
	"github.com/tupizz/restaurant-food-golang-api-fiap/internal/core/domain/validator"
	"github.com/tupizz/restaurant-food-golang-api-fiap/internal/core/usecase/dto"
	"github.com/tupizz/restaurant-food-golang-api-fiap/internal/core/usecase/ports"
)

type CreateClientUseCase interface {
	Run(ctx context.Context, input dto.ClientInput) (*entities.Client, error)
}

type createClientUseCase struct {
	clientRepository ports.ClientRepository
}

func NewCreateClientUseCase(clientRepository ports.ClientRepository) CreateClientUseCase {
	return &createClientUseCase{clientRepository}
}

func (s *createClientUseCase) Run(ctx context.Context, input dto.ClientInput) (*entities.Client, error) {
	if err := validator.IsValidCPF(input.CPF); err != nil {
		return nil, errors.New("CPF inválido")
	}

	_, err := s.clientRepository.GetByCpf(ctx, input.CPF)
	if err == nil {
		return nil, errors.New("CPF já cadastrado")
	}

	client := entities.Client{
		Name: input.Name,
		CPF:  input.CPF,
	}

	createdClient, err := s.clientRepository.Create(ctx, client)
	if err != nil {
		return nil, err
	}

	return &createdClient, nil
}
