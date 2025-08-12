package config

import (
	"log"
	"os"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

// ProductionConfig - Configuración segura para producción
type ProductionConfig struct {
	// Database
	DBUser string
	DBPass string
	DBHost string
	DBPort string
	DBName string

	// Server
	Port    string
	GinMode string

	// Security
	JWTSecret string
	APIKey    string

	// External Services
	SMTPHost string
	SMTPPort string
	SMTPUser string
	SMTPPass string

	// Performance
	MaxDBConnections    int
	DBConnectionTimeout time.Duration
	RateLimitEnabled    bool
}

// LoadProductionConfig carga configuración desde variables de ambiente
func LoadProductionConfig() *ProductionConfig {
	config := &ProductionConfig{
		// ✅ Database - REQUERIDAS
		DBUser: getEnvRequired("DB_USER"),
		DBPass: getEnvRequired("DB_PASS"),
		DBHost: getEnvRequired("DB_HOST"),
		DBPort: getEnv("DB_PORT", "3306"),
		DBName: getEnvRequired("DB_NAME"),

		// ✅ Server
		Port:    getEnv("PORT", "8080"),
		GinMode: getEnv("GIN_MODE", "release"),

		// ✅ Security - REQUERIDAS
		JWTSecret: getEnvRequired("JWT_SECRET"),
		APIKey:    getEnvRequired("API_KEY"),

		// ✅ SMTP - REQUERIDAS para notificaciones
		SMTPHost: getEnvRequired("SMTP_HOST"),
		SMTPPort: getEnv("SMTP_PORT", "587"),
		SMTPUser: getEnvRequired("SMTP_USER"),
		SMTPPass: getEnvRequired("SMTP_PASS"),

		// ✅ Performance
		MaxDBConnections:    getEnvAsInt("MAX_DB_CONNECTIONS", 100),
		DBConnectionTimeout: time.Duration(getEnvAsInt("DB_TIMEOUT_SECONDS", 30)) * time.Second,
		RateLimitEnabled:    getEnvAsBool("RATE_LIMIT_ENABLED", true),
	}

	// ✅ Configurar Gin Mode automáticamente
	gin.SetMode(config.GinMode)

	// ✅ Validar configuración crítica
	config.validate()

	return config
}

// validate verifica que la configuración sea válida
func (c *ProductionConfig) validate() {
	log.Println("🔍 Validating production configuration...")

	// ✅ Verificar que JWT_SECRET sea suficientemente fuerte
	if len(c.JWTSecret) < 32 {
		log.Fatal("❌ JWT_SECRET must be at least 32 characters long")
	}

	// ✅ Verificar que API_KEY sea suficientemente fuerte
	if len(c.APIKey) < 16 {
		log.Fatal("❌ API_KEY must be at least 16 characters long")
	}

	// ✅ Verificar conexión de base de datos
	if c.DBHost == "" || c.DBUser == "" || c.DBPass == "" {
		log.Fatal("❌ Database configuration incomplete")
	}

	log.Println("✅ Production configuration validated successfully")
}

// getEnvRequired obtiene variable de ambiente requerida
func getEnvRequired(key string) string {
	value := os.Getenv(key)
	if value == "" {
		log.Fatalf("❌ Required environment variable %s is not set", key)
	}
	return value
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

// getEnvAsBool obtiene variable de ambiente como bool
func getEnvAsBool(key string, defaultValue bool) bool {
	if value := os.Getenv(key); value != "" {
		return value == "true" || value == "1"
	}
	return defaultValue
}
