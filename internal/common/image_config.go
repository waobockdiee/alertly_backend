package common

import (
	"os"
)

// GetImageBaseURL retorna la URL base para las imágenes
// En desarrollo: http://192.168.1.66:8080
// En producción: https://cdn.alertly.ca o similar
func GetImageBaseURL() string {
	baseURL := os.Getenv("IMAGE_BASE_URL")
	if baseURL == "" {
		// Fallback para desarrollo
		baseURL = "http://192.168.1.66:8080"
	}
	return baseURL
}

// GetImageURL construye la URL completa para una imagen
func GetImageURL(filename string) string {
	return GetImageBaseURL() + "/uploads/" + filename
}

// IsProductionEnvironment verifica si estamos en producción
func IsProductionEnvironment() bool {
	return os.Getenv("NODE_ENV") == "production"
}
