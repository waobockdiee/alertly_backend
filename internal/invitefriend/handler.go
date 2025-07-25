package invitefriend

import (
	"alertly/internal/auth"
	"alertly/internal/database"
	"alertly/internal/response"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

func Save(c *gin.Context) {
	var invitation Invitation
	var err error

	invitation.AccountID, err = auth.GetUserFromContext(c)

	if err != nil {
		log.Printf("Error: %v", err)
		response.Send(c, http.StatusUnauthorized, true, "error", nil)
		return
	}

	repo := NewRepository(database.DB)
	service := NewService(repo)

	err = service.Save(invitation)

	if err != nil {
		response.Send(c, http.StatusInternalServerError, true, "error saving invitation", nil)
	}

	response.Send(c, http.StatusOK, false, "Invitation saved successfully", nil)
}
