package account

import (
	"alertly/internal/auth"
	"alertly/internal/database"
	"alertly/internal/response"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

func GetHistory(c *gin.Context) {
	var accountID int64
	var err error

	accountID, err = auth.GetUserFromContext(c)

	if err != nil {
		log.Printf("Error: %v", err)
		response.Send(c, http.StatusUnauthorized, true, "Unauthorized", nil)
		return
	}

	repo := NewRepository(database.DB)
	service := NewService(repo)

	data, err := service.GetHistory(accountID)

	if err != nil {
		log.Printf("Error: %v", err)
		response.Send(c, http.StatusInternalServerError, true, "Error getting history", err)
		return
	}

	response.Send(c, http.StatusOK, false, "success", data)

}

func ClearHistory(c *gin.Context) {
	var accountID int64
	var err error

	accountID, err = auth.GetUserFromContext(c)

	if err != nil {
		log.Printf("Error: %v", err)
		response.Send(c, http.StatusUnauthorized, true, "Unauthorized", nil)
		return
	}

	repo := NewRepository(database.DB)
	service := NewService(repo)

	err = service.ClearHistory(accountID)

	if err != nil {
		log.Printf("Error: %v", err)
		response.Send(c, http.StatusInternalServerError, true, "Error", err)
		return
	}

	response.Send(c, http.StatusOK, false, "success", nil)

}

func DeleteAccount(c *gin.Context) {
	response.Send(c, http.StatusOK, false, "success", nil)
}
