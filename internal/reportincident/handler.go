package reportincident

import (
	"alertly/internal/auth"
	"alertly/internal/database"
	"alertly/internal/response"
	"fmt"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

func ReportIncident(c *gin.Context) {
	var report Report
	var err error

	if err := c.ShouldBindJSON(&report); err != nil {
		log.Printf("Error al decodificar JSON: %v", err)
		response.Send(c, http.StatusBadRequest, true, "Invalid inputs. Please check the information and try again.", nil)
		return
	}

	report.AccountID, err = auth.GetUserFromContext(c)

	if err != nil {
		log.Printf("Error: %v", err)
		response.Send(c, http.StatusUnauthorized, true, "error", nil)
		return
	}

	repo := NewRepository(database.DB)
	service := NewService(repo)

	fmt.Printf("Report: %+v\n", report)

	err = service.ReportIncident(report)

	if err != nil {
		log.Printf("Error al reportar el incidente: %v", err)
		response.Send(c, http.StatusInternalServerError, true, "Something went wrong while reporting the incident. Please try again later.", nil)
		return
	}

	response.Send(c, http.StatusOK, false, "Incident reported successfully", nil)
}
