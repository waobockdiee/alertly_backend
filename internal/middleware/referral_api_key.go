package middleware

import (
	"net/http"
	"os"
	"strings"

	"github.com/gin-gonic/gin"
)

// ReferralAPIKeyMiddleware valida el API Key para endpoints de referrals
// Este middleware protege los endpoints que el backend web consume
func ReferralAPIKeyMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Obtener el API Key esperado de las variables de entorno
		expectedAPIKey := os.Getenv("REFERRAL_API_KEY")

		if expectedAPIKey == "" {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
				"error": "Server configuration error: REFERRAL_API_KEY not set",
			})
			return
		}

		// Obtener el header Authorization
		authHeader := c.GetHeader("Authorization")

		if authHeader == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error": "Missing or invalid authorization header",
			})
			return
		}

		// Verificar formato: "Bearer <api_key>"
		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error": "Invalid authorization format. Expected: Bearer <api_key>",
			})
			return
		}

		providedAPIKey := parts[1]

		// Comparar el API Key proporcionado con el esperado
		if providedAPIKey != expectedAPIKey {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error": "Invalid API key",
			})
			return
		}

		// API Key válido, continuar con el siguiente handler
		c.Next()
	}
}
