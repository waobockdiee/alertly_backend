package getclusterby

import (
	"log"
	"net/http"
	"strconv"

	"alertly/internal/auth"
	"alertly/internal/database"
	"alertly/internal/response"

	"github.com/gin-gonic/gin"
)

func View(c *gin.Context) {

	accountID, err := auth.GetUserFromContext(c)
	if err != nil {
		response.Send(c, http.StatusInternalServerError, true, "We couldn’t verify your session. Please log in again.", err.Error())
		return
	}

	idStr := c.Param("incl_id")
	inclId, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		log.Printf("ID inválido: %v", err)
		response.Send(c, http.StatusBadRequest, true, "The provided incident ID is not valid. Please check and try again.", err.Error())
		return
	}

	repo := NewRepository(database.DB)
	service := NewService(repo)

	result, err := service.GetIncidentBy(inclId, accountID)
	if err != nil {
		log.Printf("error fetching incident. Please try later: %v", err)
		response.Send(c, http.StatusInternalServerError, true, "We couldn’t load the incident details. Please try again later.", err.Error())
		return
	}
	response.Send(c, http.StatusOK, false, "Success", result)
}
