package newincident

import (
	"alertly/internal/auth"
	"alertly/internal/database"
	"alertly/internal/media"
	"alertly/internal/response"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"alertly/internal/common"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

var validate = validator.New()

func Create(c *gin.Context) {
	// Parse multipart form (límite 10 MB)
	// var accountID int64
	accountID, err := auth.GetUserFromContext(c)

	if err != nil {
		log.Printf("error getting accountID: %v", err.Error())
		response.Send(c, http.StatusInternalServerError, true, "error", err.Error())
		return
	}

	if err := c.Request.ParseMultipartForm(10 << 20); err != nil {
		log.Printf("Error parsing multipart form: %v", err)
		response.Send(c, http.StatusBadRequest, true, "Error parsing form", err.Error())
		return
	}

	// Vincular campos del formulario al struct IncidentReport
	var incident IncidentReport
	log.Printf("New Incident bind: %+v", incident)
	if err := c.ShouldBind(&incident); err != nil {
		log.Printf("Error al bindear formulario: %v", err)
		response.Send(c, http.StatusBadRequest, true, "Wrong data in", err.Error())
		return
	}

	// Validar el struct
	if err := validate.Struct(incident); err != nil {
		log.Printf("Error de validación: %v", err)
		response.Send(c, http.StatusBadRequest, true, "Bad request", nil)
		return
	}

	// Procesar el archivo enviado (campo "file")
	file, header, err := c.Request.FormFile("file")
	if err != nil {
		log.Printf("Error retrieving file: %v", err)
		response.Send(c, http.StatusBadRequest, true, "Error fetching file", nil)
		return
	}
	defer file.Close()

	// ✅ OPTIMIZACIÓN: Procesamiento de imágenes asíncrono
	// 1. Crear archivo temporal para el archivo original
	ext := filepath.Ext(header.Filename)
	tmpFile, err := os.CreateTemp("", "orig_*"+ext)
	if err != nil {
		log.Printf("Error creating temp file: %v", err)
		response.Send(c, http.StatusInternalServerError, true, "Error creating folder", nil)
		return
	}
	tmpFilePath := tmpFile.Name()

	// Copiar el contenido al archivo temporal
	if _, err := io.Copy(tmpFile, file); err != nil {
		log.Printf("Error saving temp file: %v", err)
		response.Send(c, http.StatusInternalServerError, true, "Error saving tmpl file", nil)
		return
	}
	tmpFile.Close()

	// ✅ SOLUCIÓN: Guardar imagen original inmediatamente para mostrar en frontend
	uploadDir := "uploads"
	originalFileName := fmt.Sprintf("alerty_%d%s", time.Now().UnixNano(), ext)
	originalFilePath := filepath.Join(uploadDir, originalFileName)

	// Crear carpeta uploads si no existe
	if err := os.MkdirAll(uploadDir, 0755); err != nil {
		log.Printf("Error creating uploads directory: %v", err)
		response.Send(c, http.StatusInternalServerError, true, "Error creating uploads folder", nil)
		return
	}

	// Copiar archivo temporal a carpeta uploads
	if err := copyFile(tmpFilePath, originalFilePath); err != nil {
		log.Printf("Error copying file to uploads: %v", err)
		response.Send(c, http.StatusInternalServerError, true, "Error saving image", nil)
		return
	}

	// ✅ OPTIMIZACIÓN: Guardar incidente con imagen original inmediatamente
	// Usar la URL completa que el frontend puede acceder
	incident.Media.Uri = common.GetImageURL(originalFileName)

	// Continuar con la lógica original de guardado en la base de datos
	repo := NewRepository(database.DB)
	service := NewService(repo)
	incident.AccountId = accountID

	result, err := service.Save(incident)
	if err != nil {
		log.Printf("error saving incident: %v", err)
		response.Send(c, http.StatusInternalServerError, true, "error saving incident. please try again", nil)
		return
	}

	// ✅ OPTIMIZACIÓN: Procesar imagen optimizada en background
	go func() {
		// ✅ Asegurar que el archivo temporal se elimine al final
		defer os.Remove(tmpFilePath)

		// Llamar a ProcessImage usando el archivo temporal
		processedFilePath, err := media.ProcessImage(tmpFilePath, uploadDir)
		if err != nil {
			log.Printf("⚠️ Error processing image for incident %d: %v", result.InreId, err)
			return
		}

		// ✅ Convertir ruta absoluta a URL completa para el frontend
		processedFileName := filepath.Base(processedFilePath)
		processedFullURL := common.GetImageURL(processedFileName)

		// ✅ Actualizar la ruta del archivo procesado en la base de datos
		if err := repo.UpdateIncidentMediaPath(result.InreId, processedFullURL); err != nil {
			log.Printf("⚠️ Error updating media path for incident %d: %v", result.InreId, err)
		} else {
			log.Printf("✅ Image processed successfully for incident %d: %s", result.InreId, processedFullURL)
		}

		// ✅ Actualizar también el cluster si existe
		if result.InclId != 0 {
			if err := repo.UpdateClusterMediaPath(result.InclId, processedFullURL); err != nil {
				log.Printf("⚠️ Error updating cluster media path for %d: %v", result.InclId, err)
			} else {
				log.Printf("✅ Cluster media updated for %d: %s", result.InclId, processedFullURL)
			}
		}
	}()

	log.Printf("success: %v", result)
	response.Send(c, http.StatusOK, false, "Thank you for your report! We've received your incident and will review it shortly.", result)
}

// ✅ Función auxiliar para copiar archivos
func copyFile(src, dst string) error {
	sourceFile, err := os.Open(src)
	if err != nil {
		return err
	}
	defer sourceFile.Close()

	destFile, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer destFile.Close()

	_, err = io.Copy(destFile, sourceFile)
	return err
}
