package auth

import (
	"errors"

	"github.com/gin-gonic/gin"
)

// GetUserFromContext extrae el AccountId del contexto y lo retorna como int64.
// Maneja el caso en que el valor se almacene como int64 o float64.
func GetUserFromContext(c *gin.Context) (int64, error) {
	v, exists := c.Get("AccountId")
	if !exists {
		return 0, errors.New("account id not found in context")
	}

	// Intenta como int64
	if id, ok := v.(int64); ok {
		return id, nil
	}
	// Si no, prueba con float64 y conviértelo a int64
	if idFloat, ok := v.(float64); ok {
		return int64(idFloat), nil
	}
	return 0, errors.New("account id has invalid type")
}
