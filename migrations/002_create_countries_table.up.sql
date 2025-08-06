CREATE TABLE IF NOT EXISTS tm_countries (
    id CHAR(36) NOT NULL PRIMARY KEY,
    capital VARCHAR(255) NULL,
    citizenship VARCHAR(255) NULL,
    country_code VARCHAR(3) NOT NULL,
    currency_name VARCHAR(255) NULL,
    currency_code VARCHAR(255) NULL,
    currency_sub_unit VARCHAR(255) NULL,
    currency_symbol VARCHAR(3) NULL,
    full_name VARCHAR(255) NULL,
    iso_3166_2 VARCHAR(2) NOT NULL,
    iso_3166_3 VARCHAR(3) NOT NULL,
    name VARCHAR(255) NOT NULL,
    region_code VARCHAR(3) NULL,
    sub_region_code VARCHAR(3) NULL,
    eea BOOLEAN NOT NULL DEFAULT FALSE,
    calling_code VARCHAR(3) NOT NULL,
    flag VARCHAR(6) NULL,
    latitude VARCHAR(255) NULL,
    longitude VARCHAR(255) NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_country_code ON tm_countries(country_code);
CREATE INDEX idx_iso_3166_2 ON tm_countries(iso_3166_2);
CREATE INDEX idx_iso_3166_3 ON tm_countries(iso_3166_3);
CREATE INDEX idx_calling_code ON tm_countries(calling_code);
CREATE INDEX idx_latitude ON tm_countries(latitude);
CREATE INDEX idx_longitude ON tm_countries(longitude); 