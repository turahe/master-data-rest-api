CREATE TABLE IF NOT EXISTS tm_districts (
    id BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    city_id BIGINT UNSIGNED NOT NULL,
    name VARCHAR(255) NOT NULL,
    code VARCHAR(10) NULL,
    latitude VARCHAR(255) NULL,
    longitude VARCHAR(255) NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    FOREIGN KEY (city_id) REFERENCES tm_cities(id) ON DELETE CASCADE,
    INDEX idx_city_id (city_id),
    INDEX idx_name (name),
    INDEX idx_code (code),
    INDEX idx_latitude (latitude),
    INDEX idx_longitude (longitude)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci; 