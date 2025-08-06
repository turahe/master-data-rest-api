CREATE TABLE IF NOT EXISTS tm_banks (
    id CHAR(36) NOT NULL PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    alias VARCHAR(255) NOT NULL,
    company VARCHAR(255) NOT NULL,
    code VARCHAR(255) NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_code ON tm_banks(code);
CREATE INDEX idx_name ON tm_banks(name);
CREATE INDEX idx_alias ON tm_banks(alias);
CREATE INDEX idx_company ON tm_banks(company); 