CREATE TABLE IF NOT EXISTS tm_villages (
    id CHAR(36) NOT NULL PRIMARY KEY,
    district_id CHAR(36) NOT NULL,
    name VARCHAR(255) NOT NULL,
    code VARCHAR(10) NOT NULL,
    latitude VARCHAR(255) NULL,
    longitude VARCHAR(255) NULL,
    postal_code VARCHAR(10) NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT fk_villages_district_id FOREIGN KEY (district_id) REFERENCES tm_districts(id) ON DELETE CASCADE
);

CREATE INDEX idx_district_id ON tm_villages(district_id);
CREATE INDEX idx_code ON tm_villages(code);
CREATE INDEX idx_name ON tm_villages(name);
CREATE INDEX idx_latitude ON tm_villages(latitude);
CREATE INDEX idx_longitude ON tm_villages(longitude);
CREATE INDEX idx_postal_code ON tm_villages(postal_code); 