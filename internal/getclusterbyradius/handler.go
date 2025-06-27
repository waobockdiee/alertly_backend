package getclusterbyradius

import (
	"alertly/internal/database"
	"alertly/internal/response"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

var validate = validator.New()

func GetByRadius(c *gin.Context) {
	var inputs Inputs
	if err := c.ShouldBindUri(&inputs); err != nil {
		log.Printf("Error al bindear URI: %v", err)
		response.Send(c, http.StatusBadRequest, true, "Invalid data in the URL. Please check and try again.", nil)
		return
	}
	if err := validate.Struct(inputs); err != nil {
		log.Printf("Error de validación: %v", err)
		response.Send(c, http.StatusBadRequest, true, "Some fields are missing or incorrect. Please review and try again.", nil)
		return
	}
	if err := c.ShouldBindQuery(&inputs); err != nil {
		log.Printf("Error en query params: %v", err)
		response.Send(c, http.StatusBadRequest, true, "Invalid query parameters. Please check and try again.", nil)
		return
	}

	repo := NewRepository(database.DB)
	service := NewService(repo)

	result, err := service.GetClustersByRadius(inputs)
	if err != nil {
		log.Printf("Error loading clusters by radius: %v", err)
		response.Send(c, http.StatusInternalServerError, true, "Could not load clusters. Please try again later.", nil)
		return
	}
	response.Send(c, http.StatusOK, false, "Success", result)
}
