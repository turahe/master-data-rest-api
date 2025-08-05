CREATE TABLE IF NOT EXISTS tm_cities (
    id BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    province_id BIGINT UNSIGNED NOT NULL,
    name VARCHAR(255) NOT NULL,
    type VARCHAR(255) NULL,
    code VARCHAR(10) NULL,
    latitude VARCHAR(255) NULL,
    longitude VARCHAR(255) NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    FOREIGN KEY (province_id) REFERENCES tm_provinces(id) ON DELETE CASCADE,
    INDEX idx_province_id (province_id),
    INDEX idx_name (name),
    INDEX idx_code (code),
    INDEX idx_latitude (latitude),
    INDEX idx_longitude (longitude)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci; 