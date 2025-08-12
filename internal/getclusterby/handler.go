package getclusterby

import (
	"fmt"
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
		response.Send(c, http.StatusInternalServerError, true, "We couldn't verify your session. Please log in again.", nil)
		return
	}

	idStr := c.Param("incl_id")
	inclId, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		log.Printf("ID inválido: %v", err)
		response.Send(c, http.StatusBadRequest, true, "The provided incident ID is not valid. Please check and try again.", nil)
		return
	}

	repo := NewRepository(database.DB)
	service := NewService(repo)

	result, err := service.GetIncidentBy(inclId, accountID)
	if err != nil {
		log.Printf("error fetching incident. Please try later: %v", err)
		response.Send(c, http.StatusInternalServerError, true, "We couldn't load the incident details. Please try again later.", nil)
		return
	}
	response.Send(c, http.StatusOK, false, "Success", result)
}

// ViewPublic - Endpoint público para landing pages web (sin autenticación)
func ViewPublic(c *gin.Context) {
	idStr := c.Param("incl_id")

	// ✅ VALIDACIÓN: Solo IDs numéricos válidos
	inclId, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil || inclId <= 0 {
		log.Printf("ID inválido para endpoint público: %v", idStr)
		response.Send(c, http.StatusBadRequest, true, "Invalid incident ID format.", nil)
		return
	}

	// ✅ VALIDACIÓN: Limitar rango de IDs para prevenir enumeration
	if inclId > 999999999 { // Límite razonable
		log.Printf("ID fuera de rango para endpoint público: %v", inclId)
		response.Send(c, http.StatusBadRequest, true, "Incident ID out of range.", nil)
		return
	}

	repo := NewRepository(database.DB)
	service := NewService(repo)

	// ✅ SEGURIDAD: Para endpoint público, usamos accountID = 0 (sin usuario autenticado)
	// Esto significa que NO se incluirán datos sensibles como get_account_already_voted, etc.
	result, err := service.GetIncidentBy(inclId, 0)
	if err != nil {
		log.Printf("error fetching incident for public view (ID: %d): %v", inclId, err)

		// ✅ SEGURIDAD: No revelar detalles del error
		response.Send(c, http.StatusNotFound, true, "Incident not found or not available for public viewing.", nil)
		return
	}

	// ✅ CACHE: Agregar headers de cache para reducir load
	c.Header("Cache-Control", "public, max-age=300") // 5 minutos
	c.Header("ETag", fmt.Sprintf("\"%d-%d\"", inclId, result.CounterTotalVotes))

	response.Send(c, http.StatusOK, false, "Success", result)
}
