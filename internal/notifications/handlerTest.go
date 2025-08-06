package notifications

import (
	"log"
	"net/http"

	"alertly/internal/common" // donde está tu SendExpoPush y el struct ExpoPushMessage
	"alertly/internal/response"

	"github.com/gin-gonic/gin"
	"github.com/sideshow/apns2/payload"
)

type testReq struct {
	DeviceToken string `json:"deviceToken" binding:"required"`
}

// TestPushHandler recibe un token y envía una notificación de prueba
func TestPushHandler(c *gin.Context) {
	var req testReq
	if err := c.BindJSON(&req); err != nil {
		response.Send(c, http.StatusBadRequest, true, "Bad request", err.Error())
		return
	}

	title := "Alertly Test Notification"
	message := "This is a test message from Alertly."
	deviceToken := req.DeviceToken // <<< aquí

	// Envía la notificación
	err := common.SendPush(
		common.ExpoPushMessage{Title: title, Body: message},
		deviceToken,
		payload.NewPayload().AlertTitle(title).AlertBody(message),
	)
	if err != nil {
		log.Printf("TestPushHandler error sending to %s: %v", deviceToken, err)
		response.Send(c, http.StatusInternalServerError, true, "Unauthorized", err.Error())
		return
	}

	response.Send(c, http.StatusOK, true, "success", nil)
}
