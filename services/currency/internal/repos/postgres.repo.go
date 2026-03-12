package repos

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

func (p *Postgres) GetRate(ctx context.Context, fromCurrency string, toCurrency string) (float64, error) {
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
		err := p.db.QueryRowContext(ctx, query, toCurrency).Scan(&rate)
		if err != nil {
			return 0, err
		}

		return rate, nil
	}

	if toCurrency == "USD" {
		err := p.db.QueryRowContext(ctx, query, fromCurrency).Scan(&rate)
		if err != nil {
			return 0, err
		}

		return 1 / rate, nil
	}

	/* --- --- --- */

	var fromRate, toRate float64

	err := p.db.QueryRowContext(ctx, query, fromCurrency).Scan(&fromRate)
	if err != nil {
		return 0, err
	}

	err = p.db.QueryRowContext(ctx, query, toCurrency).Scan(&toRate)
	if err != nil {
		return 0, err
	}

	return toRate / fromRate, nil
}

func (p *Postgres) GetAllRates(ctx context.Context, baseCurrency string) (map[string]float64, error) {
	query := `
	SELECT currency_code, rate FROM rates WHERE base_currency = $1
  AND updated_at > NOW() - INTERVAL '1 day'
	`

	rows, err := p.db.QueryContext(ctx, query, baseCurrency)
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

func (p *Postgres) UpdateRates(ctx context.Context, baseCurrency string, rates map[string]float64) error {
	if len(rates) == 0 {
		return fmt.Errorf("rates map is empty")
	}

	/* --- --- --- */

	tx, err := p.db.BeginTx(ctx, nil)
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
