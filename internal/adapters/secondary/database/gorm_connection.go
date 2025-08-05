package database

import (
	"fmt"
	"log"
	"time"

	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// GORMConnectionManager manages GORM database connections and migrations
type GORMConnectionManager struct {
	db     *gorm.DB
	driver GORMDriver
	config *Config
}

// NewGORMConnectionManager creates a new GORM database connection manager
func NewGORMConnectionManager(config *Config) (*GORMConnectionManager, error) {
	driver, err := NewGORMDriver(config.Driver)
	if err != nil {
		return nil, fmt.Errorf("failed to create GORM database driver: %w", err)
	}

	db, err := driver.Connect(config)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	// Configure GORM connection pool
	sqlDB, err := db.DB()
	if err != nil {
		return nil, fmt.Errorf("failed to get underlying sql.DB: %w", err)
	}

	// Configure connection pool
	sqlDB.SetMaxOpenConns(25)
	sqlDB.SetMaxIdleConns(25)
	sqlDB.SetConnMaxLifetime(5 * time.Minute)

	// Test the connection
	if err := sqlDB.Ping(); err != nil {
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	return &GORMConnectionManager{
		db:     db,
		driver: driver,
		config: config,
	}, nil
}

// GetDB returns the GORM database connection
func (cm *GORMConnectionManager) GetDB() *gorm.DB {
	return cm.db
}

// GetDriver returns the GORM database driver
func (cm *GORMConnectionManager) GetDriver() GORMDriver {
	return cm.driver
}

// GetConfig returns the database configuration
func (cm *GORMConnectionManager) GetConfig() *Config {
	return cm.config
}

// Close closes the database connection
func (cm *GORMConnectionManager) Close() error {
	if cm.db != nil {
		sqlDB, err := cm.db.DB()
		if err != nil {
			return fmt.Errorf("failed to get underlying sql.DB: %w", err)
		}
		return sqlDB.Close()
	}
	return nil
}

// AutoMigrate runs GORM auto-migrations for the given models
func (cm *GORMConnectionManager) AutoMigrate(models ...interface{}) error {
	if err := cm.db.AutoMigrate(models...); err != nil {
		return fmt.Errorf("failed to auto-migrate models: %w", err)
	}

	log.Println("GORM auto-migrations completed successfully")
	return nil
}

// SetLogLevel sets the GORM log level
func (cm *GORMConnectionManager) SetLogLevel(level logger.LogLevel) {
	cm.db.Logger = cm.db.Logger.LogMode(level)
}

// Begin starts a new transaction
func (cm *GORMConnectionManager) Begin() *gorm.DB {
	return cm.db.Begin()
}

// Transaction executes a function within a database transaction
func (cm *GORMConnectionManager) Transaction(fc func(tx *gorm.DB) error) error {
	return cm.db.Transaction(fc)
}
