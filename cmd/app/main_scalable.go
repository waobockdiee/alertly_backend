package main

import (
	"log"
	"os"

	"alertly/internal/getclusterby"
	"alertly/internal/middleware"
	redisClient "alertly/internal/redis"

	"github.com/gin-gonic/gin"
)

func setupScalableMiddleware(router *gin.Engine) error {
	// ✅ 1. Configurar Redis (opcional, fallback a in-memory)
	redisConfig := redisClient.NewConfig()
	redis, err := redisClient.NewClient(redisConfig)

	if err != nil {
		log.Printf("⚠️  Redis not available, falling back to in-memory rate limiting: %v", err)

		// ✅ Fallback: Rate limiting in-memory (menos escalable pero funcional)
		if os.Getenv("GIN_MODE") == "release" {
			router.Use(middleware.RateLimitMiddlewareStrict())
		} else {
			router.Use(middleware.RateLimitMiddleware())
		}
	} else {
		log.Println("✅ Redis connected successfully")

		// ✅ Rate limiting distribuido con Redis
		router.Use(middleware.RedisRateLimitMiddleware(redis))
	}

	// ✅ 2. CORS optimizado (O(1) lookup)
	router.Use(middleware.OptimizedCORSMiddleware())

	return nil
}

// Ejemplo de uso en main():
func exampleMain() {
	router := gin.Default()

	// ✅ Setup middleware escalable
	if err := setupScalableMiddleware(router); err != nil {
		log.Fatalf("Failed to setup middleware: %v", err)
	}

	// ✅ Endpoints públicos con rate limiting específico
	publicRoutes := router.Group("/public")

	// Si Redis está disponible, usar rate limiting Redis
	// Si no, usar el middleware in-memory existente
	if redis, err := redisClient.NewClient(redisClient.NewConfig()); err == nil {
		publicRoutes.Use(middleware.RedisPublicRateLimitMiddleware(redis))
	} else {
		publicRoutes.Use(middleware.RateLimitMiddlewarePublic())
	}

	publicRoutes.GET("/cluster/getbyid/:incl_id", getclusterby.ViewPublic)

	// ✅ Resto de rutas...

	router.Run(":8080")
}

/*
✅ VENTAJAS DE ESTA IMPLEMENTACIÓN:

1. 📊 ESCALABILIDAD:
   - Redis distribuido para múltiples instancias
   - O(1) CORS lookup en lugar de O(n)
   - Connection pooling optimizado

2. 🛡️ RESILENCIA:
   - Fallback automático si Redis falla
   - Health checks integrados
   - Timeouts configurables

3. ⚡ PERFORMANCE:
   - Pipeline Redis para operaciones atómicas
   - Conexiones reutilizables
   - Headers estáticos cacheados

4. 🔧 CONFIGURABILIDAD:
   - Variables de ambiente
   - Configuración por ambiente (dev/prod)
   - Rate limits dinámicos

📈 MÉTRICAS ESPERADAS:
- Desarrollo: 200 req/h por IP
- Producción: 60 req/h por IP
- Endpoints públicos: 10-50 req/h por IP
- Latencia CORS: ~0.1ms (vs ~1ms anterior)
- Memoria: Constante (vs crecimiento lineal)
*/
