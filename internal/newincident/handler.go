package newincident

import (
	"alertly/internal/auth"
	"alertly/internal/database"
	"alertly/internal/media"
	"alertly/internal/response"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"

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

	// Crear un archivo temporal para el archivo original
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
	// Cerrar y eliminar el archivo temporal posteriormente
	tmpFile.Close()
	defer os.Remove(tmpFilePath)

	// Definir carpeta de salida para la imagen procesada
	uploadDir := "uploads"

	// Llamar a ProcessImage usando el archivo temporal
	processedFilePath, err := media.ProcessImage(tmpFilePath, uploadDir)
	if err != nil {
		log.Printf("Error processing image: %v", err)
		response.Send(c, http.StatusInternalServerError, true, "Error processing file", nil)
		return
	}

	// Asignar la ruta completa del archivo procesado al incidente
	incident.Media.Uri = processedFilePath

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

	log.Printf("success: %v", result)
	response.Send(c, http.StatusOK, false, "Thank you for your report! We’ve received your incident and will review it shortly.", result)
}
