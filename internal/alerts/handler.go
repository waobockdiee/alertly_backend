package alerts

import (
	"alertly/internal/auth"
	"alertly/internal/database"
	"alertly/internal/response"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

func GetAlerts(c *gin.Context) {

}

// retorna el contador de los nuevas nuevas alertas para el usuario.
func GetNewAlertsCount(c *gin.Context) {

	var accountID int64
	var count int64
	var err error

	accountID, err = auth.GetUserFromContext(c)

	if err != nil {
		fmt.Println("error1", err)
		response.Send(c, http.StatusInternalServerError, true, "error", nil)
		return
	}

	repo := NewRepository(database.DB)
	service := NewService(repo)
	count, err = service.GetNewAlertsCount(accountID)

	if err != nil {
		fmt.Println("error2", err)
		response.Send(c, http.StatusInternalServerError, true, "error", nil)
		return
	}

	response.Send(c, http.StatusOK, false, "success", count)
}
