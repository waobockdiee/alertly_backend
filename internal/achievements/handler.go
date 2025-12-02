package achievements

import (
	"alertly/internal/auth"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	service Service
}

func NewHandler(service Service) *Handler {
	return &Handler{service: service}
}

// GetPending obtiene todos los achievements pendientes de mostrar (show_in_modal = 1)
func (h *Handler) GetPending(c *gin.Context) {
	// Obtener account_id del JWT (middleware.TokenAuthMiddleware() ya lo validó)
	accountID, err := auth.GetUserFromContext(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	achievements, err := h.service.GetPendingByAccountID(accountID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch achievements"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success":      true,
		"achievements": achievements,
		"count":        len(achievements),
	})
}

// MarkAsShown marca un achievement como mostrado (show_in_modal = 0)
func (h *Handler) MarkAsShown(c *gin.Context) {
	// Obtener account_id del JWT
	accountID, err := auth.GetUserFromContext(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	// Obtener acac_id del parámetro de ruta
	acacIDStr := c.Param("id")
	acacID, err := strconv.ParseInt(acacIDStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid achievement ID"})
		return
	}

	// Marcar como mostrado
	err = h.service.MarkAsShown(acacID, accountID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to mark achievement as shown"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Achievement marked as shown",
	})
}
