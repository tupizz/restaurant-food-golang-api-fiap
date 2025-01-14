package usecase

import (
	"context"

	"github.com/tupizz/restaurant-food-golang-api-fiap/internal/core/domain/entities"
	"github.com/tupizz/restaurant-food-golang-api-fiap/internal/core/usecase/ports"
)

type GetClientByCPFUseCase interface {
	Run(ctx context.Context, cpf string) (*entities.Client, error)
}

type getClientByCPFUseCase struct {
	clientRepository ports.ClientRepository
}

func NewGetClientByCPFUseCase(clientRepository ports.ClientRepository) GetClientByCPFUseCase {
	return &getClientByCPFUseCase{clientRepository}
}

func (g *getClientByCPFUseCase) Run(ctx context.Context, cpf string) (*entities.Client, error) {
	client, err := g.clientRepository.GetByCpf(ctx, cpf)
	if err != nil {
		return nil, err
	}

	return &client, nil
}
