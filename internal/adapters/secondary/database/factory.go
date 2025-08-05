package database

import (
	"strconv"

	"github.com/turahe/master-data-rest-api/configs"
)

// Factory creates database connections based on configuration
type Factory struct{}

// NewFactory creates a new database factory
func NewFactory() *Factory {
	return &Factory{}
}

// CreateConnectionManager creates a database connection manager from application config
func (f *Factory) CreateConnectionManager(config *configs.Config) (*ConnectionManager, error) {
	dbConfig := f.convertConfig(config.Database)
	return NewConnectionManager(dbConfig)
}

// CreateGORMConnectionManager creates a GORM database connection manager from application config
func (f *Factory) CreateGORMConnectionManager(config *configs.Config) (*GORMConnectionManager, error) {
	dbConfig := f.convertConfig(config.Database)
	return NewGORMConnectionManager(dbConfig)
}

// convertConfig converts application database config to driver config
func (f *Factory) convertConfig(dbConfig configs.DatabaseConfig) *Config {
	port, _ := strconv.Atoi(dbConfig.Port)
	parseTime := dbConfig.ParseTime == "true"

	return &Config{
		Driver:    DriverType(dbConfig.Driver),
		Host:      dbConfig.Host,
		Port:      port,
		Username:  dbConfig.User,
		Password:  dbConfig.Password,
		Database:  dbConfig.Name,
		SSLMode:   dbConfig.SSLMode,
		Charset:   dbConfig.Charset,
		ParseTime: parseTime,
		Loc:       dbConfig.Loc,
	}
}
