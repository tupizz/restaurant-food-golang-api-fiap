package repository

import (
	"context"

	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/tupizz/restaurant-food-golang-api-fiap/internal/core/domain/entities"
	"github.com/tupizz/restaurant-food-golang-api-fiap/internal/core/usecase/ports"
)

type paymentTaxSettingsRepository struct {
	db *pgxpool.Pool
}

func NewPaymentTaxSettingsRepository(db *pgxpool.Pool) ports.PaymentTaxSettingsRepository {
	return &paymentTaxSettingsRepository{db: db}
}

func (r *paymentTaxSettingsRepository) GetAll(ctx context.Context) ([]entities.PaymentTaxSettings, error) {
	query := `
	SELECT id, name, description, amount_type, amount_value, applicable_to, created_at, updated_at, deleted_at
	FROM payment_tax_settings
	WHERE deleted_at IS NULL
	ORDER BY name
`
	rows, err := r.db.Query(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var taxes []entities.PaymentTaxSettings
	for rows.Next() {
		var tax entities.PaymentTaxSettings
		err := rows.Scan(&tax.ID, &tax.Name, &tax.Description, &tax.AmountType, &tax.AmountValue, &tax.ApplicableTo, &tax.CreatedAt, &tax.UpdatedAt, &tax.DeletedAt)
		if err != nil {
			return nil, err
		}
		taxes = append(taxes, tax)
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}
	return taxes, nil
}
