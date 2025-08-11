package redis

import (
	"context"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
	"github.com/sirupsen/logrus"
	"github.com/turahe/master-data-rest-api/configs"
)

// Manager handles Redis connections and operations
type Manager struct {
	client *redis.Client
	config *configs.RedisConfig
	log    *logrus.Logger
}

// NewManager creates a new Redis manager
func NewManager(config *configs.RedisConfig, log *logrus.Logger) *Manager {
	if !config.Enabled {
		return &Manager{
			client: nil,
			config: config,
			log:    log,
		}
	}

	client := redis.NewClient(&redis.Options{
		Addr:         fmt.Sprintf("%s:%s", config.Host, config.Port),
		Password:     config.Password,
		DB:           config.DB,
		PoolSize:     config.PoolSize,
		DialTimeout:  config.DialTimeout,
		ReadTimeout:  config.ReadTimeout,
		WriteTimeout: config.WriteTimeout,
		MinIdleConns: 5,
		MaxRetries:   3,
	})

	return &Manager{
		client: client,
		config: config,
		log:    log,
	}
}

// Connect establishes connection to Redis
func (m *Manager) Connect(ctx context.Context) error {
	if m.client == nil {
		return fmt.Errorf("Redis is not enabled")
	}

	// Test the connection
	_, err := m.client.Ping(ctx).Result()
	if err != nil {
		return fmt.Errorf("failed to connect to Redis: %w", err)
	}

	m.log.Info("Redis connection established successfully")
	return nil
}

// Close closes the Redis connection
func (m *Manager) Close() error {
	if m.client == nil {
		return nil
	}

	return m.client.Close()
}

// GetClient returns the Redis client
func (m *Manager) GetClient() *redis.Client {
	return m.client
}

// IsEnabled returns whether Redis is enabled
func (m *Manager) IsEnabled() bool {
	return m.config.Enabled && m.client != nil
}

// HealthCheck performs a health check on Redis
func (m *Manager) HealthCheck(ctx context.Context) error {
	if !m.IsEnabled() {
		return fmt.Errorf("Redis is not enabled")
	}

	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	_, err := m.client.Ping(ctx).Result()
	return err
}
