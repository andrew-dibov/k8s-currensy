-- history/migrations/init.sql
CREATE TABLE conversion_history (
    id            VARCHAR(36) PRIMARY KEY,
    from_currency VARCHAR(3) NOT NULL,
    to_currency   VARCHAR(3) NOT NULL,
    amount        DECIMAL(20,6) NOT NULL,
    result        DECIMAL(20,6) NOT NULL,
    rate          DECIMAL(20,6) NOT NULL,
    created_at    TIMESTAMP NOT NULL
);

CREATE INDEX idx_history_from ON conversion_history(from_currency);
CREATE INDEX idx_history_to ON conversion_history(to_currency);
CREATE INDEX idx_history_created ON conversion_history(created_at);