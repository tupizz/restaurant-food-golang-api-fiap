package domain

import (
	"context"
	"github.com/tupizz/restaurant-food-golang-api-fiap/internal/domain/entity"
)

type UserRepository interface {
	GetAll(ctx context.Context) ([]entity.User, error)
	Create(ctx context.Context, user entity.User) (entity.User, error)
}

type ClientRepository interface {
	Create(ctx context.Context, client entity.Client) (entity.Client, error)
	GetByCpf(ctx context.Context, cpf string) (entity.Client, error)
}
