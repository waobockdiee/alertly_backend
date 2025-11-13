package media

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

// ReprocessImageHandler maneja la solicitud para reprocesar una imagen existente con pixelado
func ReprocessImageHandler(c *gin.Context) {
	// Obtener ID del incidente de los parámetros
	inclIdStr := c.Param("incl_id")
	inclId, err := strconv.ParseInt(inclIdStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   true,
			"message": "Invalid incident ID",
		})
		return
	}

	// TODO: Implementar lógica para:
	// 1. Obtener la URL de la imagen actual del incidente
	// 2. Descargar la imagen desde S3
	// 3. Procesarla con pixelado usando ProcessImage
	// 4. Actualizar la URL en la base de datos

	c.JSON(http.StatusOK, gin.H{
		"error":   false,
		"message": fmt.Sprintf("Reprocessing incident %d - Feature under development", inclId),
	})
}

// TestPixelationHandler endpoint de prueba para verificar que la detección funciona
func TestPixelationHandler(c *gin.Context) {
	// Este endpoint serviría para subir una imagen de prueba y ver los resultados
	// de detección sin afectar la base de datos

	c.JSON(http.StatusOK, gin.H{
		"error":   false,
		"message": "Test pixelation endpoint - Upload a test image to see detection results",
		"status":  "Pure Go face and license plate detection enabled",
	})
}
