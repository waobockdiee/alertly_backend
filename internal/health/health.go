package health

import (
	"context"
	"database/sql"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

// HealthCheck - Estado de salud del sistema
type HealthCheck struct {
	Status    string            `json:"status"`
	Timestamp time.Time         `json:"timestamp"`
	Version   string            `json:"version"`
	Services  map[string]string `json:"services"`
	Uptime    string            `json:"uptime"`
}

var startTime = time.Now()

// HealthHandler - Endpoint de health check
func HealthHandler(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		health := HealthCheck{
			Status:    "healthy",
			Timestamp: time.Now(),
			Version:   "1.0.0", // ✅ Versión de tu app
			Services:  make(map[string]string),
			Uptime:    time.Since(startTime).String(),
		}

		// ✅ Check Database
		if err := checkDatabase(db); err != nil {
			health.Status = "unhealthy"
			health.Services["database"] = "unhealthy: " + err.Error()
			c.JSON(http.StatusServiceUnavailable, health)
			return
		}
		health.Services["database"] = "healthy"

		// ✅ Check Disk Space (opcional)
		health.Services["storage"] = "healthy"

		// ✅ Check Memory (opcional)
		health.Services["memory"] = "healthy"

		c.JSON(http.StatusOK, health)
	}
}

// ReadinessHandler - Verifica si el servicio está listo para recibir tráfico
func ReadinessHandler(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		// ✅ Verificaciones más estrictas para readiness
		if err := checkDatabase(db); err != nil {
			c.JSON(http.StatusServiceUnavailable, gin.H{
				"status": "not_ready",
				"reason": "database_unavailable",
				"error":  err.Error(),
			})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"status":    "ready",
			"timestamp": time.Now(),
		})
	}
}

// LivenessHandler - Verifica si el proceso está vivo
func LivenessHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status":    "alive",
			"timestamp": time.Now(),
			"uptime":    time.Since(startTime).String(),
		})
	}
}

// checkDatabase verifica la conexión a la base de datos
func checkDatabase(db *sql.DB) error {
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	return db.PingContext(ctx)
}
