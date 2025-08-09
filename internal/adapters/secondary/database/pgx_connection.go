package database

import (
	"context"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/jackc/pgx/v5/tracelog"
	"github.com/turahe/master-data-rest-api/configs"
	"github.com/turahe/master-data-rest-api/pkg/logger"
)

// PgxConnection represents a PostgreSQL connection using pgx
type PgxConnection struct {
	pool   *pgxpool.Pool
	db     configs.DatabaseConfig
	logger *logger.Logger
}

// NewPgxConnection creates a new pgx connection
func NewPgxConnection(config configs.DatabaseConfig) *PgxConnection {
	return &PgxConnection{
		db: config,
	}
}

// NewPgxConnectionWithLogger creates a new pgx connection with logger
func NewPgxConnectionWithLogger(config configs.DatabaseConfig, appLogger *logger.Logger) *PgxConnection {
	return &PgxConnection{
		db:     config,
		logger: appLogger,
	}
}

// Connect establishes a connection to PostgreSQL using pgx
func (p *PgxConnection) Connect() error {
	dsn := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=%s&timezone=UTC",
		p.db.User,
		p.db.Password,
		p.db.Host,
		p.db.Port,
		p.db.Name,
		p.db.SSLMode,
	)

	config, err := pgxpool.ParseConfig(dsn)
	if err != nil {
		return fmt.Errorf("failed to parse database config: %w", err)
	}

	// Configure connection pool
	config.MaxConns = 30
	config.MinConns = 5
	config.MaxConnLifetime = time.Hour
	config.MaxConnIdleTime = time.Minute * 30
	config.HealthCheckPeriod = time.Minute

	// Configure logging if logger is available
	if p.logger != nil && p.db.LogQueries {
		pgxLogger := logger.NewPgxLogger(p.logger)
		logLevel := logger.GetPgxLogLevel(p.db.LogLevel)

		config.ConnConfig.Tracer = &tracelog.TraceLog{
			Logger:   pgxLogger,
			LogLevel: logLevel,
		}

		p.logger.WithComponent("database").Info("Database query logging enabled")
	}

	// Create connection pool
	pool, err := pgxpool.NewWithConfig(context.Background(), config)
	if err != nil {
		return fmt.Errorf("failed to create connection pool: %w", err)
	}

	// Test the connection
	if err := pool.Ping(context.Background()); err != nil {
		pool.Close()
		return fmt.Errorf("failed to ping database: %w", err)
	}

	p.pool = pool
	return nil
}

// GetPool returns the pgx connection pool
func (p *PgxConnection) GetPool() *pgxpool.Pool {
	return p.pool
}

// Close closes the database connection
func (p *PgxConnection) Close() {
	if p.pool != nil {
		p.pool.Close()
	}
}

// Ping tests the database connection
func (p *PgxConnection) Ping() error {
	if p.pool == nil {
		return fmt.Errorf("connection pool is not initialized")
	}
	return p.pool.Ping(context.Background())
}

// BeginTx starts a new transaction
func (p *PgxConnection) BeginTx(ctx context.Context) (pgx.Tx, error) {
	if p.pool == nil {
		return nil, fmt.Errorf("connection pool is not initialized")
	}
	return p.pool.Begin(ctx)
}

// Exec executes a query without returning any rows
func (p *PgxConnection) Exec(ctx context.Context, sql string, arguments ...interface{}) error {
	if p.pool == nil {
		return fmt.Errorf("connection pool is not initialized")
	}
	_, err := p.pool.Exec(ctx, sql, arguments...)
	return err
}

// Query executes a query that returns rows
func (p *PgxConnection) Query(ctx context.Context, sql string, args ...interface{}) (pgx.Rows, error) {
	if p.pool == nil {
		return nil, fmt.Errorf("connection pool is not initialized")
	}
	return p.pool.Query(ctx, sql, args...)
}

// QueryRow executes a query that is expected to return at most one row
func (p *PgxConnection) QueryRow(ctx context.Context, sql string, args ...interface{}) pgx.Row {
	return p.pool.QueryRow(ctx, sql, args...)
}
