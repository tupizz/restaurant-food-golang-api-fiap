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

type ProductRepository interface {
	Create(ctx context.Context, product entity.Product) (entity.Product, error)
	GetById(ctx context.Context, id int) (entity.Product, error)
	GetAll(ctx context.Context, page int, pageSize int) ([]entity.Product, int, error)
	Delete(ctx context.Context, id int) error
	Update(ctx context.Context, product entity.Product) (entity.Product, error)
}
