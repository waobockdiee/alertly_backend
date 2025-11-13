package middleware

import (
	"net/http"
	"os"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"golang.org/x/time/rate"
)

// RateLimiter almacena los limitadores por IP
type RateLimiter struct {
	limiters map[string]*rate.Limiter
	mu       sync.RWMutex
	r        rate.Limit
	b        int
}

// NewRateLimiter crea un nuevo rate limiter
func NewRateLimiter(r rate.Limit, b int) *RateLimiter {
	return &RateLimiter{
		limiters: make(map[string]*rate.Limiter),
		r:        r,
		b:        b,
	}
}

// getLimiter obtiene o crea un limiter para una IP específica
func (rl *RateLimiter) getLimiter(ip string) *rate.Limiter {
	rl.mu.Lock()
	defer rl.mu.Unlock()

	limiter, exists := rl.limiters[ip]
	if !exists {
		limiter = rate.NewLimiter(rl.r, rl.b)
		rl.limiters[ip] = limiter
	}

	return limiter
}

// RateLimitMiddleware crea un middleware de rate limiting
func RateLimitMiddleware() gin.HandlerFunc {
	// ✅ Configuración: 50 requests por segundo por IP (más permisivo para apps interactivas)
	rl := NewRateLimiter(rate.Every(time.Second/50), 100)

	return func(c *gin.Context) {
		ip := c.ClientIP()
		limiter := rl.getLimiter(ip)

		if !limiter.Allow() {
			c.AbortWithStatusJSON(http.StatusTooManyRequests, gin.H{
				"error":       "Too many requests. Please try again later.",
				"retry_after": "1 second",
			})
			return
		}

		c.Next()
	}
}

// RateLimitMiddlewareStrict crea un rate limiter más estricto para endpoints críticos
func RateLimitMiddlewareStrict() gin.HandlerFunc {
	// ✅ Configuración más estricta: 20 requests por segundo por IP
	rl := NewRateLimiter(rate.Every(time.Second/20), 40)

	return func(c *gin.Context) {
		ip := c.ClientIP()
		limiter := rl.getLimiter(ip)

		if !limiter.Allow() {
			c.AbortWithStatusJSON(http.StatusTooManyRequests, gin.H{
				"error":       "Rate limit exceeded. Please slow down your requests.",
				"retry_after": "1 second",
			})
			return
		}

		c.Next()
	}
}

// RateLimitMiddlewarePublic crea un rate limiter para endpoints públicos
func RateLimitMiddlewarePublic() gin.HandlerFunc {
	// ✅ Configuración dinámica según ambiente
	var rateLimiter *RateLimiter

	if os.Getenv("GIN_MODE") == "release" {
		// ✅ PRODUCCIÓN: 10 requests por segundo, burst de 20 (para mapas interactivos)
		rateLimiter = NewRateLimiter(rate.Every(time.Second/10), 20)
	} else {
		// ✅ DESARROLLO: 20 requests por segundo, burst de 40
		rateLimiter = NewRateLimiter(rate.Every(time.Second/20), 40)
	}

	return func(c *gin.Context) {
		ip := c.ClientIP()
		limiter := rateLimiter.getLimiter(ip)

		if !limiter.Allow() {
			c.AbortWithStatusJSON(http.StatusTooManyRequests, gin.H{
				"error":       "Too many requests. Please slow down.",
				"retry_after": "1 second",
				"remaining":   "0",
			})
			return
		}

		c.Next()
	}
}
