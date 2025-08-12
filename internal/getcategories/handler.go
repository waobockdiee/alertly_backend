package getcategories

import (
	"alertly/internal/common"
	"log"
	"net/http"
	"time"

	"alertly/internal/database"
	"alertly/internal/response"

	"github.com/gin-gonic/gin"
)

func GetCategories(c *gin.Context) {
	// ✅ OPTIMIZACIÓN: Cache de categorías por 5 minutos
	cacheKey := "categories_all"
	if cached, found := common.GlobalCache.Get(cacheKey); found {
		log.Println("✅ Categories served from cache")
		response.Send(c, http.StatusOK, false, "Success", cached)
		return
	}

	repo := NewRepository(database.DB)
	service := NewService(repo)

	result, err := service.GetCategories()
	if err != nil {
		log.Printf("error al obtener las categorias. Por favor intentalo mas tarde: %v", err)
		response.Send(c, http.StatusInternalServerError, true, "Unable to load categories. Please try again later.", nil)
		return
	}

	// ✅ Guardar en cache por 5 minutos
	common.GlobalCache.Set(cacheKey, result, 5*time.Minute)
	log.Println("✅ Categories cached for 5 minutes")

	response.Send(c, http.StatusOK, false, "Success", result)
}
