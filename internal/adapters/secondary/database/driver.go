package database

import (
	"database/sql"
	"fmt"
)

// DriverType represents the type of database driver
type DriverType string

const (
	MySQL      DriverType = "mysql"
	PostgreSQL DriverType = "postgres"
	SQLite     DriverType = "sqlite"
)

// Driver interface defines the contract for database drivers
type Driver interface {
	// Connect establishes a connection to the database
	Connect(config *Config) (*sql.DB, error)

	// GetDriverName returns the driver name for database/sql
	GetDriverName() string

	// GetDSN returns the Data Source Name for the connection
	GetDSN(config *Config) string

	// GetMigrationDriver returns the migration driver name
	GetMigrationDriver() string
}

// Config holds database configuration
type Config struct {
	Driver    DriverType
	Host      string
	Port      int
	Username  string
	Password  string
	Database  string
	SSLMode   string
	Charset   string
	ParseTime bool
	Loc       string
}

// NewDriver creates a new database driver based on the driver type
func NewDriver(driverType DriverType) (Driver, error) {
	switch driverType {
	case MySQL:
		return &MySQLDriver{}, nil
	case PostgreSQL:
		return &PostgreSQLDriver{}, nil
	case SQLite:
		return &SQLiteDriver{}, nil
	default:
		return nil, fmt.Errorf("unsupported database driver: %s", driverType)
	}
}

// MySQLDriver implements the Driver interface for MySQL
type MySQLDriver struct{}

func (d *MySQLDriver) Connect(config *Config) (*sql.DB, error) {
	dsn := d.GetDSN(config)
	return sql.Open(d.GetDriverName(), dsn)
}

func (d *MySQLDriver) GetDriverName() string {
	return "mysql"
}

func (d *MySQLDriver) GetDSN(config *Config) string {
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s",
		config.Username,
		config.Password,
		config.Host,
		config.Port,
		config.Database)

	if config.Charset != "" {
		dsn += "?charset=" + config.Charset
	}

	if config.ParseTime {
		if config.Charset != "" {
			dsn += "&parseTime=true"
		} else {
			dsn += "?parseTime=true"
		}
	}

	if config.Loc != "" {
		if config.Charset != "" || config.ParseTime {
			dsn += "&loc=" + config.Loc
		} else {
			dsn += "?loc=" + config.Loc
		}
	}

	return dsn
}

func (d *MySQLDriver) GetMigrationDriver() string {
	return "mysql"
}

// PostgreSQLDriver implements the Driver interface for PostgreSQL
type PostgreSQLDriver struct{}

func (d *PostgreSQLDriver) Connect(config *Config) (*sql.DB, error) {
	dsn := d.GetDSN(config)
	return sql.Open(d.GetDriverName(), dsn)
}

func (d *PostgreSQLDriver) GetDriverName() string {
	return "postgres"
}

func (d *PostgreSQLDriver) GetDSN(config *Config) string {
	dsn := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s",
		config.Host,
		config.Port,
		config.Username,
		config.Password,
		config.Database)

	if config.SSLMode != "" {
		dsn += " sslmode=" + config.SSLMode
	}

	return dsn
}

func (d *PostgreSQLDriver) GetMigrationDriver() string {
	return "postgres"
}

// SQLiteDriver implements the Driver interface for SQLite
type SQLiteDriver struct{}

func (d *SQLiteDriver) Connect(config *Config) (*sql.DB, error) {
	dsn := d.GetDSN(config)
	return sql.Open(d.GetDriverName(), dsn)
}

func (d *SQLiteDriver) GetDriverName() string {
	return "sqlite3"
}

func (d *SQLiteDriver) GetDSN(config *Config) string {
	return config.Database
}

func (d *SQLiteDriver) GetMigrationDriver() string {
	return "sqlite3"
}
