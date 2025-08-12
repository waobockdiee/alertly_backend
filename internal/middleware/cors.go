package middleware

import (
	"os"
	"sync"

	"github.com/gin-gonic/gin"
)

// CORSConfig - Configuración CORS optimizada
type CORSConfig struct {
	allowedOrigins map[string]bool
	mu             sync.RWMutex
}

// NewCORSConfig crea una nueva configuración CORS
func NewCORSConfig() *CORSConfig {
	config := &CORSConfig{
		allowedOrigins: make(map[string]bool),
	}

	// ✅ Cargar dominios permitidos según ambiente
	config.loadAllowedOrigins()

	return config
}

// loadAllowedOrigins carga los dominios permitidos
func (c *CORSConfig) loadAllowedOrigins() {
	c.mu.Lock()
	defer c.mu.Unlock()

	// ✅ PRODUCCIÓN: Solo dominios específicos
	productionOrigins := []string{
		"https://alertly.ca",
		"https://www.alertly.ca",
	}

	// ✅ DESARROLLO: Incluir localhost
	developmentOrigins := []string{
		"http://localhost:3000",
		"http://localhost:3001",
		"http://127.0.0.1:3000",
		"http://127.0.0.1:3001",
	}

	// Agregar dominios de producción
	for _, origin := range productionOrigins {
		c.allowedOrigins[origin] = true
	}

	// Agregar dominios de desarrollo si no es producción
	if os.Getenv("GIN_MODE") != "release" {
		for _, origin := range developmentOrigins {
			c.allowedOrigins[origin] = true
		}
	}
}

// IsOriginAllowed verifica si un origin está permitido (O(1))
func (c *CORSConfig) IsOriginAllowed(origin string) bool {
	if origin == "" {
		return false
	}

	c.mu.RLock()
	allowed := c.allowedOrigins[origin]
	c.mu.RUnlock()

	return allowed
}

// AddOrigin agrega un nuevo origin permitido (útil para configuración dinámica)
func (c *CORSConfig) AddOrigin(origin string) {
	c.mu.Lock()
	c.allowedOrigins[origin] = true
	c.mu.Unlock()
}

// RemoveOrigin remueve un origin (útil para configuración dinámica)
func (c *CORSConfig) RemoveOrigin(origin string) {
	c.mu.Lock()
	delete(c.allowedOrigins, origin)
	c.mu.Unlock()
}

// GetAllowedOrigins retorna todos los origins permitidos
func (c *CORSConfig) GetAllowedOrigins() []string {
	c.mu.RLock()
	defer c.mu.RUnlock()

	origins := make([]string, 0, len(c.allowedOrigins))
	for origin := range c.allowedOrigins {
		origins = append(origins, origin)
	}

	return origins
}

// OptimizedCORSMiddleware - Middleware CORS optimizado
func OptimizedCORSMiddleware() gin.HandlerFunc {
	corsConfig := NewCORSConfig()

	return func(c *gin.Context) {
		origin := c.Request.Header.Get("Origin")

		// ✅ O(1) lookup en lugar de O(n) loop
		if corsConfig.IsOriginAllowed(origin) {
			c.Writer.Header().Set("Access-Control-Allow-Origin", origin)
		}

		// ✅ Headers estáticos (no cambian por request)
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT, DELETE")

		// ✅ Handle preflight
		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	}
}
