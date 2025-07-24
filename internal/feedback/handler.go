package feedback

import (
	"alertly/internal/auth"
	"alertly/internal/database"
	"alertly/internal/response"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

func SendFeedback(c *gin.Context) {
	var feedback Feedback

	var err error

	if err := c.ShouldBindJSON(&feedback); err != nil {
		log.Printf("Error: %v", err)
		response.Send(c, http.StatusBadRequest, true, "Invalid entry data", nil)
		return
	}

	feedback.AccountID, err = auth.GetUserFromContext(c)

	if err != nil {
		log.Printf("Error: %v", err)
		response.Send(c, http.StatusUnauthorized, true, "Unauthorized", nil)
		return
	}

	repo := NewRepository(database.DB)
	service := NewService(repo)

	err = service.SendFeedback(feedback)

	if err != nil {
		log.Printf("Error: %v", err)
		response.Send(c, http.StatusInternalServerError, true, "Error sending feedback", err.Error())
	}

	response.Send(c, http.StatusOK, false, "Feedback sent successfully", nil)
}
