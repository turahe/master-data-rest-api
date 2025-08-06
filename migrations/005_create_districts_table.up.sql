CREATE TABLE IF NOT EXISTS tm_districts (
    id CHAR(36) NOT NULL PRIMARY KEY,
    city_id CHAR(36) NOT NULL,
    name VARCHAR(255) NOT NULL,
    code VARCHAR(10) NOT NULL,
    latitude VARCHAR(255) NULL,
    longitude VARCHAR(255) NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT fk_districts_city_id FOREIGN KEY (city_id) REFERENCES tm_cities(id) ON DELETE CASCADE
);

CREATE INDEX idx_city_id ON tm_districts(city_id);
CREATE INDEX idx_code ON tm_districts(code);
CREATE INDEX idx_name ON tm_districts(name);
CREATE INDEX idx_latitude ON tm_districts(latitude);
CREATE INDEX idx_longitude ON tm_districts(longitude); 