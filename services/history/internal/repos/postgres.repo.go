package repos

import (
	"context"
	"database/sql"

	pb "history/internal/protos/history"
)

type PostgresRepo struct {
	db *sql.DB
}

func NewPostgresRepo(db *sql.DB) *PostgresRepo {
	return &PostgresRepo{
		db: db,
	}
}

/* --- --- --- */

func (p *PostgresRepo) SaveConversion(ctx context.Context, record *pb.ConversionRecord) error {
	query := "INSERT INTO conversion_history (id, from_currency, to_currency, amount, result, rate, created_at) VALUES ($1, $2, $3, $4, $5, $6, $7)"
	_, err := p.db.ExecContext(ctx, query,
		record.Id,
		record.FromCurrency,
		record.ToCurrency, record.Amount,
		record.Result,
		record.Rate,
		record.CreatedAt.AsTime(),
	)

	return err
}
