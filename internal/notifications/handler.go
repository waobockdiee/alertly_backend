package notifications

import (
	"alertly/internal/auth"
	"alertly/internal/database"
	"alertly/internal/response"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

func SaveDeviceToken(c *gin.Context) {

	var accountID int64
	var err error

	var req struct {
		DeviceToken string `json:"deviceToken" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "deviceToken is required"})
		return
	}

	accountID, err = auth.GetUserFromContext(c)

	if err != nil {
		log.Printf("Error: %v", err)
		response.Send(c, http.StatusUnauthorized, true, "Unauthorized", nil)
		return
	}

	repo := NewRepository(database.DB)
	if err := repo.SaveDeviceToken(accountID, req.DeviceToken); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "could not save device token"})
		return
	}
	c.Status(http.StatusNoContent)
}
