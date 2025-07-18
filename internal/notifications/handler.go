package notifications

import (
	"alertly/internal/database"
	"net/http"

	"github.com/gin-gonic/gin"
)

func RegisterDeviceTokenHandler(c *gin.Context) {
	var req struct {
		DeviceToken string `json:"deviceToken" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "deviceToken is required"})
		return
	}
	accountID, exists := c.Get("accountID") // según tu middleware de auth
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	repo := NewRepository(database.DB)
	if err := repo.SaveDeviceToken(accountID.(int64), req.DeviceToken); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "could not save device token"})
		return
	}
	c.Status(http.StatusNoContent)
}
