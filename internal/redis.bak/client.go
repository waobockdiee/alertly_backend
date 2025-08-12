package redis

import (
	"context"
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/go-redis/redis/v8"
)

// Config - Configuración de Redis
type Config struct {
	Host         string
	Port         int
	Password     string
	DB           int
	PoolSize     int
	MinIdleConns int
	MaxRetries   int
	DialTimeout  time.Duration
	ReadTimeout  time.Duration
	WriteTimeout time.Duration
}

// NewConfig crea configuración por defecto
func NewConfig() *Config {
	return &Config{
		Host:         getEnv("REDIS_HOST", "localhost"),
		Port:         getEnvAsInt("REDIS_PORT", 6379),
		Password:     getEnv("REDIS_PASSWORD", ""),
		DB:           getEnvAsInt("REDIS_DB", 0),
		PoolSize:     getEnvAsInt("REDIS_POOL_SIZE", 100),
		MinIdleConns: getEnvAsInt("REDIS_MIN_IDLE_CONNS", 10),
		MaxRetries:   getEnvAsInt("REDIS_MAX_RETRIES", 3),
		DialTimeout:  time.Duration(getEnvAsInt("REDIS_DIAL_TIMEOUT", 5)) * time.Second,
		ReadTimeout:  time.Duration(getEnvAsInt("REDIS_READ_TIMEOUT", 3)) * time.Second,
		WriteTimeout: time.Duration(getEnvAsInt("REDIS_WRITE_TIMEOUT", 3)) * time.Second,
	}
}

// NewClient crea un cliente Redis optimizado
func NewClient(config *Config) (*redis.Client, error) {
	rdb := redis.NewClient(&redis.Options{
		Addr:         fmt.Sprintf("%s:%d", config.Host, config.Port),
		Password:     config.Password,
		DB:           config.DB,
		PoolSize:     config.PoolSize,
		MinIdleConns: config.MinIdleConns,
		MaxRetries:   config.MaxRetries,
		DialTimeout:  config.DialTimeout,
		ReadTimeout:  config.ReadTimeout,
		WriteTimeout: config.WriteTimeout,

		// ✅ Optimizaciones para alta concurrencia
		PoolTimeout: 30 * time.Second,
		IdleTimeout: 5 * time.Minute,

		// ✅ Health checks
		MaxConnAge: 30 * time.Minute,
	})

	// ✅ Test conexión
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err := rdb.Ping(ctx).Result()
	if err != nil {
		return nil, fmt.Errorf("failed to connect to Redis: %w", err)
	}

	return rdb, nil
}

// getEnv obtiene variable de ambiente con valor por defecto
func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

// getEnvAsInt obtiene variable de ambiente como int
func getEnvAsInt(key string, defaultValue int) int {
	if value := os.Getenv(key); value != "" {
		if intValue, err := strconv.Atoi(value); err == nil {
			return intValue
		}
	}
	return defaultValue
}

// HealthCheck verifica la salud de la conexión Redis
func HealthCheck(client *redis.Client) error {
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	_, err := client.Ping(ctx).Result()
	return err
}
