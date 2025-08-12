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
	// ✅ Configuración: 10 requests por segundo por IP
	rl := NewRateLimiter(rate.Every(time.Second), 10)

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
	// ✅ Configuración más estricta: 5 requests por segundo por IP
	rl := NewRateLimiter(rate.Every(time.Second), 5)

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
		// ✅ PRODUCCIÓN: Más estricto - 3 requests por minuto, burst de 2
		rateLimiter = NewRateLimiter(rate.Every(20*time.Second), 2)
	} else {
		// ✅ DESARROLLO: Más permisivo - 10 requests por minuto, burst de 5
		rateLimiter = NewRateLimiter(rate.Every(6*time.Second), 5)
	}

	return func(c *gin.Context) {
		ip := c.ClientIP()
		limiter := rateLimiter.getLimiter(ip)

		if !limiter.Allow() {
			retryAfter := "20 seconds"
			if os.Getenv("GIN_MODE") != "release" {
				retryAfter = "6 seconds"
			}

			c.AbortWithStatusJSON(http.StatusTooManyRequests, gin.H{
				"error":       "Public API rate limit exceeded. Please wait before making another request.",
				"retry_after": retryAfter,
				"remaining":   "0",
			})
			return
		}

		c.Next()
	}
}
