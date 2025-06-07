package middleware

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
)

type Claims struct {
	AccountID int64 `json:"account_id"`
	jwt.RegisteredClaims
}

var jwtSecret = []byte("mi_clave_secreta")

func TokenAuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Se espera que el header tenga el formato: "Bearer <token>"
		authHeader := c.GetHeader("Authorization")
		fmt.Println("authHeader", authHeader)
		if authHeader == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Authorization token is missing."})
			return
		}

		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Invalid token format. Expected format: Bearer <token>."})
			return
		}

		tokenStr := parts[1]
		token, err := jwt.ParseWithClaims(tokenStr, &Claims{}, func(token *jwt.Token) (interface{}, error) {
			// Verifica el método de firma
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("invalid token. Please log in again: %v", token.Header["alg"])
			}
			return jwtSecret, nil
		})

		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Your session is not valid. Please log in again.: " + err.Error()})
			return
		}

		if claims, ok := token.Claims.(*Claims); ok && token.Valid {
			c.Set("AccountId", claims.AccountID) // claims.AccountID debe ser int64
			// c.Set("account_id", claims.AccountID)
			c.Next()
		} else {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Your session is not valid. Please log in again."})
			return
		}

	}
}
