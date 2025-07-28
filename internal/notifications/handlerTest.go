package notifications

import (
	"net/http"

	"alertly/internal/common" // donde está tu SendExpoPush y el struct ExpoPushMessage

	"github.com/gin-gonic/gin"
)

type testReq struct {
	DeviceToken string `json:"deviceToken" binding:"required"`
}

// TestPushHandler recibe un token y envía una notificación de prueba
func TestPushHandler(c *gin.Context) {
	var req testReq
	if err := c.BindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	msg := common.ExpoPushMessage{
		To:    req.DeviceToken,
		Title: "Prueba de Notificación",
		Body:  "¡Funciona tu sistema de push!",
		Data:  map[string]interface{}{"test": true},
	}

	if err := common.SendExpoPush(msg); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "sent"})
}
