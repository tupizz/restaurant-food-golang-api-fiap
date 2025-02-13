package repository

import (
	"context"

	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/tupizz/restaurant-food-golang-api-fiap/internal/core/domain/entities"
	domainError "github.com/tupizz/restaurant-food-golang-api-fiap/internal/core/domain/error"
	"github.com/tupizz/restaurant-food-golang-api-fiap/internal/core/usecase/ports"
)

type clientRepository struct {
	db *pgxpool.Pool
}

func NewClientRepository(db *pgxpool.Pool) ports.ClientRepository {
	return &clientRepository{db: db}
}

func (r *clientRepository) Create(ctx context.Context, client entities.Client) (entities.Client, error) {
	query := `INSERT INTO clients (name, cpf) VALUES ($1, $2) RETURNING id, created_at, updated_at`
	err := r.db.QueryRow(ctx, query, client.Name, client.CPF).Scan(&client.ID, &client.CreatedAt, &client.UpdatedAt)
	if err != nil {
		return entities.Client{}, err
	}

	return client, nil
}

func (r *clientRepository) GetByCpf(ctx context.Context, cpf string) (entities.Client, error) {
	query := `SELECT id, name, cpf, created_at, updated_at FROM clients WHERE cpf = $1 AND deleted_at IS NULL`

	var client entities.Client
	err := r.db.QueryRow(ctx, query, cpf).Scan(&client.ID, &client.Name, &client.CPF, &client.CreatedAt, &client.UpdatedAt)
	if err != nil {
		if err == pgx.ErrNoRows {
			return entities.Client{}, domainError.ErrNotFound("client")
		}
		return entities.Client{}, err
	}

	return client, nil
}
