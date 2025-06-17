package saveclusteraccount

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

func ToggleSaveClusterAccount(c *gin.Context) {
	var accountID int64
	var inclID int64
	var err error

	accountID, err = auth.GetUserFromContext(c)
	if err != nil {
		fmt.Println("error1", err)
		response.Send(c, http.StatusInternalServerError, true, "error", err.Error())
		return
	}

	tpmInclID := c.Param("incl_id")

	if tpmInclID == "" {
		fmt.Println("error2", err)
		response.Send(c, http.StatusBadRequest, true, "bad request", nil)
		return
	}

	inclID, err = strconv.ParseInt(tpmInclID, 10, 64)

	if err != nil {
		fmt.Println("error3", err)
		response.Send(c, http.StatusInternalServerError, true, "error converting data", nil)
		return
	}

	repo := NewRepository(database.DB)
	service := NewService(repo)

	err = service.ToggleSaveClusterAccount(accountID, inclID)

	if err != nil {
		fmt.Println("error4", err)
		response.Send(c, http.StatusInternalServerError, true, "internal error. please try later", nil)
		return
	}

	response.Send(c, http.StatusOK, false, "success", nil)

}

func GetMyList(c *gin.Context) {
	var accountID int64
	var list []MyList
	var err error

	accountID, err = auth.GetUserFromContext(c)
	if err != nil {
		fmt.Println("error1", err)
		response.Send(c, http.StatusInternalServerError, true, "error", err.Error())
		return
	}

	repo := NewRepository(database.DB)
	service := NewService(repo)

	list, err = service.GetMyList(accountID)

	if err != nil {
		response.Send(c, http.StatusInternalServerError, true, "error", err.Error())
		return
	}
	response.Send(c, http.StatusOK, false, "success", list)
}

func DeleteFollowIncident(c *gin.Context) {
	var accountID int64
	var acsID int64
	var err error

	accountID, err = auth.GetUserFromContext(c)

	if err != nil {
		log.Printf("Error: %v", err)
		response.Send(c, http.StatusInternalServerError, true, "error", nil)
		return
	}

	unformattedAcsID := c.Param("acs_id")
	acsID, err = strconv.ParseInt(unformattedAcsID, 10, 64)

	if err != nil {
		log.Printf("Error: %v", err)
		response.Send(c, http.StatusInternalServerError, true, "error", nil)
		return
	}

	repo := NewRepository(database.DB)
	service := NewService(repo)

	err = service.DeleteFollowIncident(acsID, accountID)

	if err != nil {
		log.Printf("Error: %v", err)
		response.Send(c, http.StatusInternalServerError, true, "error", nil)
		return
	}

	response.Send(c, http.StatusOK, false, "success", nil)

}
