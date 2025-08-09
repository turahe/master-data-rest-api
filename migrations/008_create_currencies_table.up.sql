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

CREATE INDEX IF NOT EXISTS idx_code_currencies ON tm_currencies(code);
CREATE INDEX IF NOT EXISTS idx_name_currencies ON tm_currencies(name);
CREATE INDEX IF NOT EXISTS idx_symbol_currencies ON tm_currencies(symbol);
CREATE INDEX IF NOT EXISTS idx_is_active_currencies ON tm_currencies(is_active); 