package profile

import (
	"alertly/internal/auth"
	"alertly/internal/database"
	"alertly/internal/response"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

// si viene un id por la url entonces obtiene ese usuario
// si no, toma el id de la session del usuario
func GetById(c *gin.Context) {
	id := c.Param("account_id")
	var accountID int64
	var err error

	// si id == "" o id == 0 entonces se da por entendido que es el perfil del owner el que se necesita consultar.
	if id == "" || id == "0" {
		accountID, err = auth.GetUserFromContext(c)
		if err != nil {
			response.Send(c, http.StatusInternalServerError, true, "error", err.Error())
			return
		}
	} else {
		accountID, err = strconv.ParseInt(id, 10, 64)
	}

	if err != nil {
		response.Send(c, http.StatusBadRequest, true, "Error converting param", nil)
		return
	}

	repo := NewRepository(database.DB)
	service := NewService(repo)

	result, err := service.GetById(accountID)

	if err != nil {
		log.Printf("error en el handler: %v", err)
		response.Send(c, http.StatusInternalServerError, true, "error fetching data, pls try later", err.Error())
		return
	}

	fmt.Println("RES:", result)

	response.Send(c, http.StatusOK, false, "Success", result)
}

func InviteFriend(c *gin.Context) {

}
