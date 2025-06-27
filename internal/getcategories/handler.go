package getcategories

import (
	"log"
	"net/http"

	"alertly/internal/database"
	"alertly/internal/response"

	"github.com/gin-gonic/gin"
)

func GetCategories(c *gin.Context) {

	repo := NewRepository(database.DB)
	service := NewService(repo)

	result, err := service.GetCategories()
	if err != nil {
		log.Printf("error al obtener las categorias. Por favor intentalo mas tarde: %v", err)
		response.Send(c, http.StatusInternalServerError, true, "Unable to load categories. Please try again later.", nil)
		return
	}
	response.Send(c, http.StatusOK, false, "Success", result)
}
