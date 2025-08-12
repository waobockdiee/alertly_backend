package middleware

import (
	"os"

	"github.com/gin-gonic/gin"
)

// SecurityHeadersMiddleware - Headers de seguridad para producción
func SecurityHeadersMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// ✅ Solo aplicar en producción
		if os.Getenv("GIN_MODE") == "release" {
			// Prevent clickjacking
			c.Header("X-Frame-Options", "DENY")

			// XSS Protection
			c.Header("X-XSS-Protection", "1; mode=block")

			// Prevent MIME type sniffing
			c.Header("X-Content-Type-Options", "nosniff")

			// Referrer Policy
			c.Header("Referrer-Policy", "strict-origin-when-cross-origin")

			// HSTS (HTTPS only)
			if c.Request.TLS != nil {
				c.Header("Strict-Transport-Security", "max-age=31536000; includeSubDomains")
			}

			// Content Security Policy (básico)
			c.Header("Content-Security-Policy", "default-src 'self'; script-src 'self'; style-src 'self' 'unsafe-inline'")

			// Remove server information
			c.Header("Server", "Alertly")
		}

		c.Next()
	}
}

// RateLimitHeadersMiddleware - Agregar headers informativos de rate limiting
func RateLimitHeadersMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// ✅ Headers informativos (no afectan funcionalidad)
		c.Header("X-RateLimit-Policy", "60 requests per hour per IP")

		c.Next()
	}
}
