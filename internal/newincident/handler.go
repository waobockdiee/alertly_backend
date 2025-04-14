package newincident

import (
	"alertly/internal/database"
	"alertly/internal/media"
	"alertly/internal/response"
	"fmt"
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
	if err := c.Request.ParseMultipartForm(10 << 20); err != nil {
		log.Printf("Error parsing multipart form: %v", err)
		response.Send(c, http.StatusBadRequest, true, "Error parsing form", err.Error())
		return
	}

	// Vincular campos del formulario al struct IncidentReport
	var incident IncidentReport
	if err := c.ShouldBind(&incident); err != nil {
		log.Printf("Error al bindear formulario: %v", err)
		response.Send(c, http.StatusBadRequest, true, "Wrong data in", err.Error())
		return
	}

	// Validar el struct
	if err := validate.Struct(incident); err != nil {
		log.Printf("Error de validación: %v", err)
		response.Send(c, http.StatusBadRequest, true, "Bad request", err.Error())
		return
	}

	fmt.Println("DEBUGING!!!", incident)

	// Procesar el archivo enviado (campo "file")
	file, header, err := c.Request.FormFile("file")
	if err != nil {
		log.Printf("Error retrieving file: %v", err)
		response.Send(c, http.StatusBadRequest, true, "Error fetching file", err.Error())
		return
	}
	defer file.Close()

	// Crear un archivo temporal para el archivo original
	ext := filepath.Ext(header.Filename)
	tmpFile, err := os.CreateTemp("", "orig_*"+ext)
	if err != nil {
		log.Printf("Error creating temp file: %v", err)
		response.Send(c, http.StatusInternalServerError, true, "Error creating folder", err.Error())
		return
	}
	tmpFilePath := tmpFile.Name()
	// Copiar el contenido al archivo temporal
	if _, err := io.Copy(tmpFile, file); err != nil {
		log.Printf("Error saving temp file: %v", err)
		response.Send(c, http.StatusInternalServerError, true, "Error saving tmpl file", err.Error())
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
		response.Send(c, http.StatusInternalServerError, true, "Error processing file", err.Error())
		return
	}

	// Asignar la ruta completa del archivo procesado al incidente
	incident.Media.Uri = processedFilePath

	// Continuar con la lógica original de guardado en la base de datos
	repo := NewRepository(database.DB)
	service := NewService(repo)

	result, err := service.Save(incident)
	if err != nil {
		log.Printf("error saving incident: %v", err)
		response.Send(c, http.StatusInternalServerError, true, "error saving incident. please try again", err.Error())
		return
	}

	response.Send(c, http.StatusOK, false, "File uploaded successfully", result)
}
