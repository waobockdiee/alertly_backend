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
		response.Send(c, http.StatusInternalServerError, true, "We couldn’t verify your session. Please log in again.", nil)
		return
	}

	var inputs Inputs

	if err := c.ShouldBindUri(&inputs); err != nil {
		log.Printf("Error binding URI: %v", err)
		response.Send(c, http.StatusBadRequest, true, "Invalid URL data. Please check and try again.", nil)
		return
	}

	repo := NewRepository(database.DB)
	service := NewService(repo)

	data, err := service.GetReel(inputs, accountID)

	if err != nil {
		log.Printf("We couldn’t load the incidents. Please try again later: %v", err)
		response.Send(c, http.StatusOK, false, "We couldn’t load the incidents. Please try again later.", data)
		return
	}

	response.Send(c, http.StatusOK, false, "success", data)
}
