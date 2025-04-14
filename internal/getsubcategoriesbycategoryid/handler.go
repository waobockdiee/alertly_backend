package getsubcategoriesbycategoryid

import (
	"log"
	"net/http"
	"strconv"

	"alertly/internal/database"
	"alertly/internal/response"

	"github.com/gin-gonic/gin"
)

func GetSubcategoriesByCategoryId(c *gin.Context) {

	idStr := c.Param("id")
	subcategoryID, err := strconv.Atoi(idStr)
	if err != nil {
		log.Printf("ID inválido: %v", err)
		response.Send(c, http.StatusBadRequest, true, "El ID proporcionado no es válido", err.Error())
		return
	}

	repo := NewRepository(database.DB)
	service := NewService(repo)

	result, err := service.GetSubcategoriesByCategoryId(subcategoryID)
	if err != nil {
		log.Printf("error al obtener las subcategorias. Por favor intentalo mas tarde: %v", err)
		response.Send(c, http.StatusInternalServerError, true, "error al obtener las subcategorias. Por favor intentalo mas tarde", err.Error())
		return
	}
	response.Send(c, http.StatusOK, false, "Success", result)
}
