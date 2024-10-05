package repository

import (
	"context"

	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/tupizz/restaurant-food-golang-api-fiap/internal/domain"
	"github.com/tupizz/restaurant-food-golang-api-fiap/internal/domain/entity"
)

type clientRepository struct {
	db *pgxpool.Pool
}

func NewClientRepository(db *pgxpool.Pool) domain.ClientRepository {
	return &clientRepository{db: db}
}

func (r *clientRepository) Create(ctx context.Context, client entity.Client) (entity.Client, error) {
	query := `INSERT INTO clients (name, cpf) VALUES ($1, $2) RETURNING id`
	err := r.db.QueryRow(ctx, query, client.Name, client.CPF).Scan(&client.ID)
	if err != nil {
		return entity.Client{}, err
	}
	return client, nil
}

func (r *clientRepository) GetByCpf(ctx context.Context, cpf string) (entity.Client, error) {
	query := `SELECT id, name, cpf FROM clients WHERE cpf = $1 AND deleted_at IS NULL`
	var client entity.Client
	err := r.db.QueryRow(ctx, query, cpf).Scan(&client.ID, &client.Name, &client.CPF)
	if err != nil {
		if err == pgx.ErrNoRows {
			return entity.Client{}, domain.ErrNotFound
		}
		return entity.Client{}, err
	}
	return client, nil
}
