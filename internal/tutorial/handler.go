package tutorial

import (
	"alertly/internal/auth"
	"alertly/internal/database"
	"alertly/internal/response"
	"net/http"

	"github.com/gin-gonic/gin"
)

type CompleteRequest struct {
	Latitude  *float32 `json:"latitude"`
	Longitude *float32 `json:"longitude"`
}

func CompleteHandler(c *gin.Context) {
	accountID, err := auth.GetUserFromContext(c)
	if err != nil {
		response.Send(c, http.StatusUnauthorized, true, "Invalid session", nil)
		return
	}

	var req CompleteRequest
	// Parse JSON, but don't fail if empty or malformed - coordinates are optional
	c.BindJSON(&req)

	repo := NewRepository(database.DB)
	service := NewService(repo)

	if err := service.FinishTutorial(accountID, req.Latitude, req.Longitude); err != nil {
		response.Send(c, http.StatusInternalServerError, true, "Could not update tutorial status", err.Error())
		return
	}

	response.Send(c, http.StatusOK, false, "Tutorial completed successfully", nil)
}
