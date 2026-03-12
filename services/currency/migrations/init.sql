CREATE TABLE rates (
    id SERIAL PRIMARY KEY,
    base_currency VARCHAR(3) NOT NULL,
    currency_code VARCHAR(3) NOT NULL,
    rate DECIMAL(20, 6) NOT NULL,
    updated_at TIMESTAMP NOT NULL DEFAULT NOW(),
    UNIQUE(base_currency, currency_code)
);

CREATE INDEX idx_rates_base_currency ON rates(base_currency);
CREATE INDEX idx_rates_updated_at ON rates(updated_at);