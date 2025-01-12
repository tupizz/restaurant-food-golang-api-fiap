package ports

import (
	"context"

	"github.com/tupizz/restaurant-food-golang-api-fiap/internal/core/domain"
)

type ClientRepository interface {
	Create(ctx context.Context, client domain.Client) (domain.Client, error)
	GetByCpf(ctx context.Context, cpf string) (domain.Client, error)
}
