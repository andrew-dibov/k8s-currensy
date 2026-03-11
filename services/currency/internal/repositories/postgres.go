package repositories

import (
	"context"
	"database/sql"
	"fmt"
)

type Postgres struct {
	db *sql.DB
}

func NewPostgres(db *sql.DB) *Postgres {
	return &Postgres{
		db: db,
	}
}

func (r *Postgres) GetRate(ctx context.Context, fromCurrency string, toCurrency string) (float64, error) {
	if fromCurrency == toCurrency {
		return 1.00, nil
	}

	querySQL := `
	SELECT rate FROM rates WHERE currency_code = $1 
	AND base_currency = 'USD' AND updated_at > NOW() - INTERVAL '1 day'
  `

	/* --- --- --- */

	var rate float64

	if fromCurrency == "USD" {
		err := r.db.QueryRowContext(ctx, querySQL, toCurrency).Scan(&rate)
		if err != nil {
			return 0, err
		}

		return rate, nil
	}

	/* --- --- --- */

	if toCurrency == "USD" {
		err := r.db.QueryRowContext(ctx, querySQL, fromCurrency).Scan(&rate)
		if err != nil {
			return 0, err
		}

		return 1 / rate, nil
	}

	/* --- --- --- */

	var fromRate, toRate float64

	err := r.db.QueryRowContext(ctx, querySQL, fromCurrency).Scan(&fromRate)
	if err != nil {
		return 0, err
	}

	err = r.db.QueryRowContext(ctx, querySQL, toCurrency).Scan(&toRate)
	if err != nil {
		return 0, err
	}

	return toRate / fromRate, nil
}

func (r *Postgres) GetAllRates(ctx context.Context, baseCurrency string) (map[string]float64, error) {
	querySQL := `
	SELECT currency_code, rate FROM rates WHERE base_currency = $1
  AND updated_at > NOW() - INTERVAL '1 day'
	`

	rows, err := r.db.QueryContext(ctx, querySQL, baseCurrency)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	rates := make(map[string]float64)

	for rows.Next() {
		var currencyCode string
		var rate float64

		if err := rows.Scan(&currencyCode, &rate); err != nil {
			return nil, err
		}

		rates[currencyCode] = rate
	}

	return rates, nil
}

func (r *Postgres) UpdateRates(ctx context.Context, baseCurrency string, rates map[string]float64) error {
	if len(rates) == 0 {
		return fmt.Errorf("empty rates")
	}

	tx, err := r.db.BeginTx(ctx, nil)

	if err != nil {
		return err
	}
	defer tx.Rollback()

	querySQL := "DELETE FROM rates WHERE base_currency = $1"

	if _, err := tx.ExecContext(ctx, querySQL, baseCurrency); err != nil {
		return err
	}

	querySQL = `
	INSERT INTO rates (base_currency, currency_code, rate, updated_at)
  VALUES ($1, $2, $3, NOW())
	`

	for currencyCode, rate := range rates {
		if _, err := tx.ExecContext(ctx, querySQL, baseCurrency, currencyCode, rate); err != nil {
			return err
		}
	}

	return tx.Commit()
}
