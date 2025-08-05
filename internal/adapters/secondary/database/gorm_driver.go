package database

import (
	"fmt"

	"gorm.io/driver/mysql"
	"gorm.io/driver/postgres"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// GORMDriver interface defines the contract for GORM database drivers
type GORMDriver interface {
	// Connect establishes a GORM connection to the database
	Connect(config *Config) (*gorm.DB, error)

	// GetDriverName returns the driver name for database/sql
	GetDriverName() string

	// GetDSN returns the Data Source Name for the connection
	GetDSN(config *Config) string

	// GetMigrationDriver returns the migration driver name
	GetMigrationDriver() string
}

// NewGORMDriver creates a new GORM database driver based on the driver type
func NewGORMDriver(driverType DriverType) (GORMDriver, error) {
	switch driverType {
	case MySQL:
		return &GORMMySQLDriver{}, nil
	case PostgreSQL:
		return &GORMPostgreSQLDriver{}, nil
	case SQLite:
		return &GORMSQLiteDriver{}, nil
	default:
		return nil, fmt.Errorf("unsupported database driver: %s", driverType)
	}
}

// GORMMySQLDriver implements the GORMDriver interface for MySQL
type GORMMySQLDriver struct{}

func (d *GORMMySQLDriver) Connect(config *Config) (*gorm.DB, error) {
	dsn := d.GetDSN(config)
	return gorm.Open(mysql.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})
}

func (d *GORMMySQLDriver) GetDriverName() string {
	return "mysql"
}

func (d *GORMMySQLDriver) GetDSN(config *Config) string {
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

func (d *GORMMySQLDriver) GetMigrationDriver() string {
	return "mysql"
}

// GORMPostgreSQLDriver implements the GORMDriver interface for PostgreSQL
type GORMPostgreSQLDriver struct{}

func (d *GORMPostgreSQLDriver) Connect(config *Config) (*gorm.DB, error) {
	dsn := d.GetDSN(config)
	return gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})
}

func (d *GORMPostgreSQLDriver) GetDriverName() string {
	return "postgres"
}

func (d *GORMPostgreSQLDriver) GetDSN(config *Config) string {
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

func (d *GORMPostgreSQLDriver) GetMigrationDriver() string {
	return "postgres"
}

// GORMSQLiteDriver implements the GORMDriver interface for SQLite
type GORMSQLiteDriver struct{}

func (d *GORMSQLiteDriver) Connect(config *Config) (*gorm.DB, error) {
	dsn := d.GetDSN(config)
	return gorm.Open(sqlite.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})
}

func (d *GORMSQLiteDriver) GetDriverName() string {
	return "sqlite3"
}

func (d *GORMSQLiteDriver) GetDSN(config *Config) string {
	return config.Database
}

func (d *GORMSQLiteDriver) GetMigrationDriver() string {
	return "sqlite3"
}
