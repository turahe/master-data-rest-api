CREATE TABLE IF NOT EXISTS tm_banks (
    id CHAR(36) NOT NULL PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    alias VARCHAR(255) NOT NULL,
    company VARCHAR(255) NOT NULL,
    code VARCHAR(255) NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX IF NOT EXISTS idx_code_banks ON tm_banks(code);
CREATE INDEX IF NOT EXISTS idx_name_banks ON tm_banks(name);
CREATE INDEX IF NOT EXISTS idx_alias_banks ON tm_banks(alias);
CREATE INDEX IF NOT EXISTS idx_company_banks ON tm_banks(company); 