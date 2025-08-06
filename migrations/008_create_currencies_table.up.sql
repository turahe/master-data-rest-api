CREATE TABLE IF NOT EXISTS tm_currencies (
    id CHAR(36) NOT NULL PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    code VARCHAR(3) NOT NULL,
    symbol VARCHAR(10) NULL,
    decimal_places INT NOT NULL DEFAULT 2,
    is_active BOOLEAN NOT NULL DEFAULT TRUE,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_code ON tm_currencies(code);
CREATE INDEX idx_name ON tm_currencies(name);
CREATE INDEX idx_symbol ON tm_currencies(symbol);
CREATE INDEX idx_is_active ON tm_currencies(is_active); 