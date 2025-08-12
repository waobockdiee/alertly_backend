package middleware

import (
	"os"
	"sync"

	"github.com/gin-gonic/gin"
)

// SafeCORSConfig - CORS optimizado pero compatible con la lógica actual
type SafeCORSConfig struct {
	allowedOrigins map[string]bool
	mu             sync.RWMutex
	fallbackMode   bool // ✅ Modo de compatibilidad
}

// NewSafeCORSConfig crea configuración CORS compatible
func NewSafeCORSConfig() *SafeCORSConfig {
	config := &SafeCORSConfig{
		allowedOrigins: make(map[string]bool),
		fallbackMode:   os.Getenv("CORS_STRICT_MODE") != "true", // ✅ Opt-in, no opt-out
	}

	config.loadAllowedOrigins()
	return config
}

func (c *SafeCORSConfig) loadAllowedOrigins() {
	c.mu.Lock()
	defer c.mu.Unlock()

	// ✅ Mismos dominios que antes, pero en map para O(1)
	productionOrigins := []string{
		"https://alertly.ca",
		"https://www.alertly.ca",
	}

	developmentOrigins := []string{
		"http://localhost:3000",
		"http://localhost:3001",
		"http://127.0.0.1:3000",
		"http://127.0.0.1:3001",
	}

	for _, origin := range productionOrigins {
		c.allowedOrigins[origin] = true
	}

	if os.Getenv("GIN_MODE") != "release" {
		for _, origin := range developmentOrigins {
			c.allowedOrigins[origin] = true
		}
	}
}

func (c *SafeCORSConfig) IsOriginAllowed(origin string) bool {
	if origin == "" {
		return c.fallbackMode // ✅ En modo compatible, permitir requests sin origin
	}

	c.mu.RLock()
	allowed := c.allowedOrigins[origin]
	c.mu.RUnlock()

	return allowed
}

// SafeCORSMiddleware - Drop-in replacement del middleware actual
func SafeCORSMiddleware() gin.HandlerFunc {
	corsConfig := NewSafeCORSConfig()

	return func(c *gin.Context) {
		origin := c.Request.Header.Get("Origin")

		// ✅ COMPATIBILIDAD: Mantener comportamiento actual por defecto
		if corsConfig.fallbackMode {
			// Modo compatible: permitir origins conocidos O usar wildcard como antes
			if origin != "" && corsConfig.IsOriginAllowed(origin) {
				c.Writer.Header().Set("Access-Control-Allow-Origin", origin)
			} else if origin == "" || !corsConfig.IsOriginAllowed(origin) {
				// ✅ FALLBACK: Comportamiento original para compatibilidad
				c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
			}
		} else {
			// Modo estricto: solo origins permitidos (cuando CORS_STRICT_MODE=true)
			if corsConfig.IsOriginAllowed(origin) {
				c.Writer.Header().Set("Access-Control-Allow-Origin", origin)
			}
		}

		// ✅ Headers exactamente iguales que antes
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT, DELETE")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	}
}
