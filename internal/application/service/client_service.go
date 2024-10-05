package service

import (
	"context"
	"errors"

	"github.com/tupizz/restaurant-food-golang-api-fiap/internal/application/dto"
	"github.com/tupizz/restaurant-food-golang-api-fiap/internal/domain"
	"github.com/tupizz/restaurant-food-golang-api-fiap/internal/domain/entity"
	"github.com/tupizz/restaurant-food-golang-api-fiap/internal/domain/validation"
)

type ClientService interface {
	CreateClient(ctx context.Context, input dto.ClientInput) (dto.ClientOutput, error)
	GetClientByCpf(ctx context.Context, cpf string) (dto.ClientOutput, error)
}

type clientService struct {
	clientRepository domain.ClientRepository
}

func NewClientService(clientRepo domain.ClientRepository) ClientService {
	return &clientService{clientRepository: clientRepo}
}

func (s *clientService) CreateClient(ctx context.Context, input dto.ClientInput) (dto.ClientOutput, error) {
	// Validação do CPF
	if !validation.IsValidCPF(input.CPF) {
		return dto.ClientOutput{}, errors.New("CPF inválido")
	}

	// Verificar se o CPF já está cadastrado
	_, err := s.clientRepository.GetByCpf(ctx, input.CPF)
	if err == nil {
		return dto.ClientOutput{}, errors.New("CPF já cadastrado")
	}

	client := entity.Client{
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

func (s *clientService) GetClientByCpf(ctx context.Context, cpf string) (dto.ClientOutput, error) {
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
