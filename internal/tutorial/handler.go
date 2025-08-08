package tutorial

import (
	"alertly/internal/auth"
	"alertly/internal/database"
	"alertly/internal/response"
	"net/http"

	"github.com/gin-gonic/gin"
)

func CompleteHandler(c *gin.Context) {
	accountID, err := auth.GetUserFromContext(c)
	if err != nil {
		response.Send(c, http.StatusUnauthorized, true, "Invalid session", nil)
		return
	}

	repo := NewRepository(database.DB)
	service := NewService(repo)

	if err := service.FinishTutorial(accountID); err != nil {
		response.Send(c, http.StatusInternalServerError, true, "Could not update tutorial status", err.Error())
		return
	}

	response.Send(c, http.StatusOK, false, "Tutorial completed successfully", nil)
}
