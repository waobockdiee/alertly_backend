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
		response.Send(c, http.StatusBadRequest, true, "The provided category ID is not valid. Please check and try again.", err.Error())
		return
	}

	repo := NewRepository(database.DB)
	service := NewService(repo)

	result, err := service.GetSubcategoriesByCategoryId(subcategoryID)
	if err != nil {
		log.Printf("We couldn’t load the subcategories. Please try again later: %v", err)
		response.Send(c, http.StatusInternalServerError, true, "We couldn’t load the subcategories. Please try again later.", err.Error())
		return
	}
	response.Send(c, http.StatusOK, false, "Success", result)
}
