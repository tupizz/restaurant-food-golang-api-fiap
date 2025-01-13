package ports

import (
	"context"

	"github.com/tupizz/restaurant-food-golang-api-fiap/internal/core/domain/entities"
)

type ClientRepository interface {
	Create(ctx context.Context, client entities.Client) (entities.Client, error)
	GetByCpf(ctx context.Context, cpf string) (entities.Client, error)
}
