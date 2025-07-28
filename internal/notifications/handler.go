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
		log.Printf("Error bindjson: %v", err)
		response.Send(c, http.StatusBadRequest, true, "error.", nil)
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
		log.Printf("Error bindjson: %v", err)
		response.Send(c, http.StatusUnauthorized, true, "could not save device token", nil)
		return
	}
	response.Send(c, http.StatusOK, false, "Success", nil)
}

func DeleteDeviceToken(c *gin.Context) {

	var accountID int64
	var err error

	var req struct {
		DeviceToken string `json:"deviceToken" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		log.Printf("Error bindjson: %v", err)
		response.Send(c, http.StatusBadRequest, true, "error.", nil)
		return
	}

	accountID, err = auth.GetUserFromContext(c)

	if err != nil {
		log.Printf("Error: %v", err)
		response.Send(c, http.StatusUnauthorized, true, "Unauthorized", nil)
		return
	}

	repo := NewRepository(database.DB)
	if err := repo.DeleteDeviceToken(accountID, req.DeviceToken); err != nil {
		log.Printf("Error bindjson: %v", err)
		response.Send(c, http.StatusUnauthorized, true, "could not delete device token", nil)
		return
	}

	log.Printf("no error deleting device token: TOKEN: %v ACCOUNT_ID: %v", req.DeviceToken, accountID)
	response.Send(c, http.StatusOK, false, "Success", nil)
}
