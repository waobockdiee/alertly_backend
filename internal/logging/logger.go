package logging

import (
	"log"
	"os"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

// SetupProductionLogging configura logging para producción
func SetupProductionLogging() {
	// ✅ En producción, solo errores y warnings
	if os.Getenv("GIN_MODE") == "release" {
		gin.SetMode(gin.ReleaseMode)

		// ✅ Configurar formato de logs
		gin.DefaultWriter = os.Stdout
		gin.DefaultErrorWriter = os.Stderr

		log.SetFlags(log.LstdFlags | log.Lshortfile)
		log.SetPrefix("[ALERTLY] ")
	}
}

// ProductionLogger - Middleware de logging para producción
func ProductionLogger() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		path := c.Request.URL.Path
		
		// Process request
		c.Next()

		// ✅ Solo log errores y requests lentos en producción
		latency := time.Since(start)
		statusCode := c.Writer.Status()

		// ✅ Log errores (4xx, 5xx)
		if statusCode >= 400 {
			log.Printf("ERROR %d %s %s %v %s",
				statusCode,
				c.Request.Method,
				path,
				latency,
				c.ClientIP(),
			)
		}

		// ✅ Log requests lentos (>1s)
		if latency > time.Second {
			log.Printf("SLOW %d %s %s %v %s",
				statusCode,
				c.Request.Method,
				path,
				latency,
				c.ClientIP(),
			)
		}

		// ✅ Log requests públicos (para monitoreo)
		if strings.HasPrefix(path, "/public/") {
			log.Printf("PUBLIC %d %s %s %v %s",
				statusCode,
				c.Request.Method,
				path,
				latency,
				c.ClientIP(),
			)
		}
	}
}
