package repos

import (
	"context"
	"database/sql"
	"fmt"
)

type PostgresRepo struct {
	db *sql.DB
}

func NewPostgresRepo(db *sql.DB) *PostgresRepo {
	return &PostgresRepo{
		db: db,
	}
}

func (postgres *PostgresRepo) GetRate(ctx context.Context, fromCurrency string, toCurrency string) (float64, error) {
	if fromCurrency == toCurrency {
		return 1.00, nil
	}

	query := `
	SELECT rate FROM rates WHERE currency_code = $1 
	AND base_currency = 'USD' AND updated_at > NOW() - INTERVAL '1 day'
  `

	/* --- --- --- */

	var rate float64

	if fromCurrency == "USD" {
		err := postgres.db.QueryRowContext(ctx, query, toCurrency).Scan(&rate)
		if err != nil {
			return 0, err
		}

		return rate, nil
	}

	if toCurrency == "USD" {
		err := postgres.db.QueryRowContext(ctx, query, fromCurrency).Scan(&rate)
		if err != nil {
			return 0, err
		}

		return 1 / rate, nil
	}

	/* --- --- --- */

	var fromRate, toRate float64

	err := postgres.db.QueryRowContext(ctx, query, fromCurrency).Scan(&fromRate)
	if err != nil {
		return 0, err
	}

	err = postgres.db.QueryRowContext(ctx, query, toCurrency).Scan(&toRate)
	if err != nil {
		return 0, err
	}

	return toRate / fromRate, nil
}

func (postgres *PostgresRepo) GetAllRates(ctx context.Context, baseCurrency string) (map[string]float64, error) {
	query := `
	SELECT currency_code, rate FROM rates WHERE base_currency = $1
  AND updated_at > NOW() - INTERVAL '1 day'
	`

	rows, err := postgres.db.QueryContext(ctx, query, baseCurrency)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	/* --- --- --- */

	rates := make(map[string]float64)

	for rows.Next() {
		var code string
		var rate float64

		if err := rows.Scan(&code, &rate); err != nil {
			return nil, err
		}

		rates[code] = rate
	}

	return rates, nil
}

func (postgres *PostgresRepo) UpdateRates(ctx context.Context, baseCurrency string, rates map[string]float64) error {
	if len(rates) == 0 {
		return fmt.Errorf("rates map is empty")
	}

	/* --- --- --- */

	tx, err := postgres.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	/* --- --- --- */

	query := "DELETE FROM rates WHERE base_currency = $1"
	if _, err := tx.ExecContext(ctx, query, baseCurrency); err != nil {
		return err
	}

	/* --- --- --- */

	query = `
	INSERT INTO rates (base_currency, currency_code, rate, updated_at)
  VALUES ($1, $2, $3, NOW())
	`

	for code, rate := range rates {
		if _, err := tx.ExecContext(ctx, query, baseCurrency, code, rate); err != nil {
			return err
		}
	}

	return tx.Commit()
}
