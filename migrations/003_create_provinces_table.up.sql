CREATE TABLE IF NOT EXISTS tm_provinces (
    id CHAR(36) NOT NULL PRIMARY KEY,
    country_id CHAR(36) NOT NULL,
    name VARCHAR(255) NOT NULL,
    code VARCHAR(10) NOT NULL,
    latitude VARCHAR(255) NULL,
    longitude VARCHAR(255) NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT fk_provinces_country_id FOREIGN KEY (country_id) REFERENCES tm_countries(id) ON DELETE CASCADE
);

CREATE INDEX idx_country_id ON tm_provinces(country_id);
CREATE INDEX idx_code ON tm_provinces(code);
CREATE INDEX idx_name ON tm_provinces(name);
CREATE INDEX idx_latitude ON tm_provinces(latitude);
CREATE INDEX idx_longitude ON tm_provinces(longitude); 