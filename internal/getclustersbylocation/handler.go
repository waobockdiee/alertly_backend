package getclustersbylocation

import (
	"alertly/internal/database"
	"alertly/internal/response"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

var validate = validator.New()

func Get(c *gin.Context) {
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

	result, err := service.GetClustersByLocation(inputs)
	if err != nil {
		log.Printf("We couldn’t load the categories. Please try again later: %v", err)
		response.Send(c, http.StatusInternalServerError, true, "We couldn’t load the categories. Please try again later.", err.Error())
		return
	}
	response.Send(c, http.StatusOK, false, "Success", result)
}
