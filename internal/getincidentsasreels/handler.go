package getincidentsasreels

import (
	"alertly/internal/auth"
	"alertly/internal/database"
	"alertly/internal/response"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

func GetReel(c *gin.Context) {

	accountID, err := auth.GetUserFromContext(c)
	if err != nil {
		log.Printf("Error handler 1 reel: %v", err)
		response.Send(c, http.StatusInternalServerError, true, "error", err.Error())
		return
	}

	var inputs Inputs

	if err := c.ShouldBindUri(&inputs); err != nil {
		log.Printf("Error binding URI: %v", err)
		response.Send(c, http.StatusBadRequest, true, "Bad request", nil)
		return
	}

	repo := NewRepository(database.DB)
	service := NewService(repo)

	data, err := service.GetReel(inputs, accountID)

	if err != nil {
		log.Printf("Error getting reel: %v", err)
		response.Send(c, http.StatusOK, false, "error getting incidents", data)
		return
	}

	response.Send(c, http.StatusOK, false, "success", data)
}
