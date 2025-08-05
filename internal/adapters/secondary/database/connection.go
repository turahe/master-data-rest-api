package database

import (
	"database/sql"
	"fmt"
	"log"
	"time"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/mysql"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/database/sqlite3"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

// ConnectionManager manages database connections and migrations
type ConnectionManager struct {
	db     *sql.DB
	driver Driver
	config *Config
}

// NewConnectionManager creates a new database connection manager
func NewConnectionManager(config *Config) (*ConnectionManager, error) {
	driver, err := NewDriver(config.Driver)
	if err != nil {
		return nil, fmt.Errorf("failed to create database driver: %w", err)
	}

	db, err := driver.Connect(config)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	// Configure connection pool
	db.SetMaxOpenConns(25)
	db.SetMaxIdleConns(25)
	db.SetConnMaxLifetime(5 * time.Minute)

	// Test the connection
	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	return &ConnectionManager{
		db:     db,
		driver: driver,
		config: config,
	}, nil
}

// GetDB returns the database connection
func (cm *ConnectionManager) GetDB() *sql.DB {
	return cm.db
}

// GetDriver returns the database driver
func (cm *ConnectionManager) GetDriver() Driver {
	return cm.driver
}

// GetConfig returns the database configuration
func (cm *ConnectionManager) GetConfig() *Config {
	return cm.config
}

// Close closes the database connection
func (cm *ConnectionManager) Close() error {
	if cm.db != nil {
		return cm.db.Close()
	}
	return nil
}

// RunMigrations runs database migrations
func (cm *ConnectionManager) RunMigrations(migrationsPath string) error {
	driverName := cm.driver.GetMigrationDriver()
	dsn := cm.driver.GetDSN(cm.config)

	migrationURL := fmt.Sprintf("%s://%s", driverName, dsn)

	m, err := migrate.New(
		fmt.Sprintf("file://%s", migrationsPath),
		migrationURL,
	)
	if err != nil {
		return fmt.Errorf("failed to create migration instance: %w", err)
	}
	defer m.Close()

	if err := m.Up(); err != nil && err != migrate.ErrNoChange {
		return fmt.Errorf("failed to run migrations: %w", err)
	}

	log.Println("Database migrations completed successfully")
	return nil
}

// RollbackMigrations rolls back database migrations
func (cm *ConnectionManager) RollbackMigrations(migrationsPath string, steps int) error {
	driverName := cm.driver.GetMigrationDriver()
	dsn := cm.driver.GetDSN(cm.config)

	migrationURL := fmt.Sprintf("%s://%s", driverName, dsn)

	m, err := migrate.New(
		fmt.Sprintf("file://%s", migrationsPath),
		migrationURL,
	)
	if err != nil {
		return fmt.Errorf("failed to create migration instance: %w", err)
	}
	defer m.Close()

	if err := m.Steps(-steps); err != nil {
		return fmt.Errorf("failed to rollback migrations: %w", err)
	}

	log.Printf("Database migrations rolled back %d steps successfully", steps)
	return nil
}

// GetMigrationVersion returns the current migration version
func (cm *ConnectionManager) GetMigrationVersion(migrationsPath string) (uint, error) {
	driverName := cm.driver.GetMigrationDriver()
	dsn := cm.driver.GetDSN(cm.config)

	migrationURL := fmt.Sprintf("%s://%s", driverName, dsn)

	m, err := migrate.New(
		fmt.Sprintf("file://%s", migrationsPath),
		migrationURL,
	)
	if err != nil {
		return 0, fmt.Errorf("failed to create migration instance: %w", err)
	}
	defer m.Close()

	version, dirty, err := m.Version()
	if err != nil {
		return 0, fmt.Errorf("failed to get migration version: %w", err)
	}

	if dirty {
		return 0, fmt.Errorf("database is in dirty state")
	}

	return version, nil
}

// ForceMigrationVersion forces the migration version
func (cm *ConnectionManager) ForceMigrationVersion(migrationsPath string, version uint) error {
	driverName := cm.driver.GetMigrationDriver()
	dsn := cm.driver.GetDSN(cm.config)

	migrationURL := fmt.Sprintf("%s://%s", driverName, dsn)

	m, err := migrate.New(
		fmt.Sprintf("file://%s", migrationsPath),
		migrationURL,
	)
	if err != nil {
		return fmt.Errorf("failed to create migration instance: %w", err)
	}
	defer m.Close()

	if err := m.Force(int(version)); err != nil {
		return fmt.Errorf("failed to force migration version: %w", err)
	}

	log.Printf("Migration version forced to %d", version)
	return nil
}
