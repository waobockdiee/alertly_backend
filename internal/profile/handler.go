package profile

import (
	"alertly/internal/auth"
	"alertly/internal/database"
	"alertly/internal/response"
	"log"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

// GET /api/account/profile/get_by_id/:account_id
func GetById(c *gin.Context) {
	idParam := c.Param("account_id")
	if idParam == "" {
		response.Send(c, http.StatusBadRequest, true, "account_id is required", nil)
		return
	}

	var accountID int64
	var err error

	if idParam == "0" {
		// Perfil propio: extraer AccountID del context
		accountID, err = auth.GetUserFromContext(c)
		if err != nil {
			// Si no hay token o es inválido, devolvemos 401
			response.Send(c, http.StatusUnauthorized, true, "unauthorized", nil)
			return
		}
	} else {
		// Perfil de otro usuario: convertir param a entero
		accountID, err = strconv.ParseInt(idParam, 10, 64)
		if err != nil {
			response.Send(c, http.StatusBadRequest, true, "invalid account_id parameter", nil)
			return
		}
	}

	// Llamada al servicio
	repo := NewRepository(database.DB)
	service := NewService(repo)

	profileData, err := service.GetById(accountID)
	if err != nil {
		log.Printf("error fetching profile for accountID %d: %v", accountID, err)
		response.Send(c, http.StatusInternalServerError, true,
			"error fetching profile, please try later", nil)
		return
	}

	response.Send(c, http.StatusOK, false, "Success", profileData)
}

func InviteFriend(c *gin.Context) {

}

func BlockAccount(c *gin.Context) {
	idParam := c.Param("account_id")
	if idParam == "" {
		response.Send(c, http.StatusBadRequest, true, "account_id is required", nil)
		return
	}

	blockerID, err := auth.GetUserFromContext(c)
	if err != nil {
		response.Send(c, http.StatusUnauthorized, true, "unauthorized", nil)
		return
	}

	blockedID, err := strconv.ParseInt(idParam, 10, 64)
	if err != nil {
		response.Send(c, http.StatusBadRequest, true, "invalid account_id parameter", nil)
		return
	}

	if blockerID == blockedID {
		response.Send(c, http.StatusBadRequest, true, "You cannot block your own account.", nil)
		return
	}

	repo := NewRepository(database.DB)
	service := NewService(repo)

	if err := service.BlockAccount(BlockAccountInput{BlockerID: blockerID, BlockedID: blockedID}); err != nil {
		log.Printf("Error blocking account: %v", err)
		response.Send(c, http.StatusInternalServerError, true, "error processing block, please try later", nil)
		return
	}

	response.Send(c, http.StatusOK, false, "User blocked successfully.", nil)
}

func UnblockAccount(c *gin.Context) {
	idParam := c.Param("account_id")
	if idParam == "" {
		response.Send(c, http.StatusBadRequest, true, "account_id is required", nil)
		return
	}

	blockerID, err := auth.GetUserFromContext(c)
	if err != nil {
		response.Send(c, http.StatusUnauthorized, true, "unauthorized", nil)
		return
	}

	blockedID, err := strconv.ParseInt(idParam, 10, 64)
	if err != nil {
		response.Send(c, http.StatusBadRequest, true, "invalid account_id parameter", nil)
		return
	}

	repo := NewRepository(database.DB)
	service := NewService(repo)

	if err := service.UnblockAccount(blockerID, blockedID); err != nil {
		log.Printf("Error unblocking account: %v", err)
		response.Send(c, http.StatusInternalServerError, true, "error processing unblock, please try later", nil)
		return
	}

	response.Send(c, http.StatusOK, false, "User unblocked successfully.", nil)
}

func ReportAccount(c *gin.Context) {

	idParam := c.Param("account_id")
	if idParam == "" {
		log.Printf("Error ReportAccount: %v", idParam)
		response.Send(c, http.StatusBadRequest, true, "account_id is required", nil)
		return
	}

	var accountID int64
	var accountIDWhosReporting int64
	var err error
	var report ReportAccountInput

	if err := c.ShouldBind(&report); err != nil {
		log.Printf("Error al bindear formulario: %v", err)
		response.Send(c, http.StatusBadRequest, true, "Wrong data in", nil)
		return
	}

	accountIDWhosReporting, err = auth.GetUserFromContext(c)

	if err != nil {
		log.Printf("Error: %v", err)
		response.Send(c, http.StatusUnauthorized, true, "unauthorized", nil)
		return
	}

	accountID, err = strconv.ParseInt(idParam, 10, 64)
	if err != nil {
		log.Printf("Error: %v", err)
		response.Send(c, http.StatusBadRequest, true, "invalid account_id parameter", nil)
		return
	}

	// if accountIDWhosReporting == accountID {
	// 	response.Send(c, http.StatusBadRequest, true, "You cannot report your own account.", nil)
	// 	return
	// }

	report.AccountIDWhosReporting = accountIDWhosReporting
	report.AccountID = accountID

	repo := NewRepository(database.DB)
	service := NewService(repo)

	err = service.ReportAccount(report)

	if err != nil {
		log.Printf("Error: %v", err)
		response.Send(c, http.StatusInternalServerError, true, "error saving report, please try later", nil)
		return
	}

	response.Send(c, http.StatusOK, false, "Thank you for your contribution! Your report strengthens trust in our platform and helps safeguard the entire community.", nil)
}
