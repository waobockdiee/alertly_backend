package middleware

import (
	"database/sql"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

// PremiumMiddleware verifica que el usuario tenga suscripción premium activa y válida
// DEBE usarse DESPUÉS de TokenAuthMiddleware() para que AccountId esté disponible en el contexto
func PremiumMiddleware(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Obtener account_id del contexto (ya validado por TokenAuthMiddleware)
		accountIDInterface, exists := c.Get("AccountId")
		if !exists {
			log.Printf("⚠️ PremiumMiddleware: AccountId not found in context")
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error":   "Authentication required",
				"message": "Please log in to access this feature",
			})
			return
		}

		accountID, ok := accountIDInterface.(int64)
		if !ok {
			log.Printf("❌ PremiumMiddleware: Invalid AccountId type in context")
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
				"error": "Invalid account ID format",
			})
			return
		}

		// Consultar estado premium del usuario en la base de datos
		var isPremium bool
		var premiumExpiresAt sql.NullTime

		query := `
			SELECT is_premium, premium_expired_date
			FROM account
			WHERE account_id = $1
		`

		err := db.QueryRow(query, accountID).Scan(&isPremium, &premiumExpiresAt)
		if err != nil {
			if err == sql.ErrNoRows {
				log.Printf("⚠️ PremiumMiddleware: Account %d not found", accountID)
				c.AbortWithStatusJSON(http.StatusNotFound, gin.H{
					"error": "User account not found",
				})
				return
			}

			log.Printf("❌ PremiumMiddleware: Error checking premium status for account %d: %v", accountID, err)
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
				"error": "Error verifying premium status",
			})
			return
		}

		// Validar que el usuario tiene premium activo
		if !isPremium {
			log.Printf("🔒 PremiumMiddleware: Account %d attempted to access premium feature (not premium)", accountID)
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{
				"error":   "Premium subscription required",
				"message": "This feature requires an active premium subscription. Upgrade to premium to unlock this feature.",
				"code":    "PREMIUM_REQUIRED",
			})
			return
		}

		// Validar que la suscripción no haya expirado (si hay fecha de expiración)
		if premiumExpiresAt.Valid && premiumExpiresAt.Time.Before(time.Now()) {
			log.Printf("🔒 PremiumMiddleware: Account %d premium subscription expired at %v", accountID, premiumExpiresAt.Time)
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{
				"error":      "Premium subscription expired",
				"message":    "Your premium subscription has expired. Please renew to continue using premium features.",
				"code":       "PREMIUM_EXPIRED",
				"expired_at": premiumExpiresAt.Time.Format(time.RFC3339),
			})
			return
		}

		// ✅ Usuario tiene premium válido, continuar con el request
		log.Printf("✅ PremiumMiddleware: Account %d has valid premium, allowing access", accountID)
		c.Set("IsPremium", true)
		c.Next()
	}
}
