CREATE TABLE IF NOT EXISTS tm_cities (
    id CHAR(36) NOT NULL PRIMARY KEY,
    province_id CHAR(36) NOT NULL,
    name VARCHAR(255) NOT NULL,
    code VARCHAR(10) NOT NULL,
    latitude VARCHAR(255) NULL,
    longitude VARCHAR(255) NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT fk_cities_province_id FOREIGN KEY (province_id) REFERENCES tm_provinces(id) ON DELETE CASCADE
);

CREATE INDEX idx_province_id ON tm_cities(province_id);
CREATE INDEX idx_code ON tm_cities(code);
CREATE INDEX idx_name ON tm_cities(name);
CREATE INDEX idx_latitude ON tm_cities(latitude);
CREATE INDEX idx_longitude ON tm_cities(longitude); 