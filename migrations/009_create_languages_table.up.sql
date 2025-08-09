CREATE TABLE IF NOT EXISTS tm_languages (
    id CHAR(36) NOT NULL PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    code VARCHAR(10) NOT NULL,
    is_active BOOLEAN NOT NULL DEFAULT TRUE,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX IF NOT EXISTS idx_code_languages ON tm_languages(code);
CREATE INDEX IF NOT EXISTS idx_name_languages ON tm_languages(name);
CREATE INDEX IF NOT EXISTS idx_is_active_languages ON tm_languages(is_active); 