CREATE TABLE IF NOT EXISTS tm_villages (
    id BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    district_id BIGINT UNSIGNED NOT NULL,
    name VARCHAR(255) NOT NULL,
    code VARCHAR(10) NULL,
    postal_code VARCHAR(255) NULL,
    latitude VARCHAR(255) NULL,
    longitude VARCHAR(255) NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    FOREIGN KEY (district_id) REFERENCES tm_districts(id) ON DELETE CASCADE,
    INDEX idx_district_id (district_id),
    INDEX idx_name (name),
    INDEX idx_code (code),
    INDEX idx_postal_code (postal_code),
    INDEX idx_latitude (latitude),
    INDEX idx_longitude (longitude)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci; 