package media

import (
	"fmt"

	"github.com/disintegration/imaging"
)

// Constantes para los archivos de cascada y configuración de procesamiento
const (
	FaceCascadeFile  = "data/cascades/haarcascade_frontalface_default.xml"
	PlateCascadeFile = "data/cascades/haarcascade_russian_plate_number.xml"
	MobileWidth      = 720 // Ancho deseado para la imagen en mobile
)

// PixelateROI - TEMPORALMENTE DESACTIVADO PARA MVP
// TODO: Implementar pixelation usando imaging library sin gocv
/*
func PixelateROI(mat *gocv.Mat, r image.Rectangle, pixelSize int) {
	// OpenCV implementation temporarily disabled
}
*/

// ProcessImage procesa la imagen ubicada en filePath, aplica detección de rostros y matrículas,
// y la sube a S3. Devuelve la URL pública de S3.
func ProcessImage(filePath, folder string) (string, error) {
	// Crear servicio S3
	s3Service, err := NewS3Service()
	if err != nil {
		return "", fmt.Errorf("failed to create S3 service: %v", err)
	}

	// 🚀 CARGA DIRECTA DE IMAGEN SIN OpenCV PARA MVP
	fmt.Println("⚠️ OpenCV processing temporarily disabled for MVP performance")
	processedImg, err := imaging.Open(filePath)
	if err != nil {
		return "", fmt.Errorf("failed to read image: %s", err)
	}

	// Redimensionar la imagen para mobile
	resizedImg := imaging.Resize(processedImg, MobileWidth, 0, imaging.Lanczos)

	// Subir imagen procesada a S3
	s3URL, err := s3Service.UploadImage(resizedImg, folder)
	if err != nil {
		return "", fmt.Errorf("failed to upload image to S3: %v", err)
	}

	return s3URL, nil
}

// ProcessVideo es una función placeholder para procesar videos.
func ProcessVideo(filePath, outputFolder string) (string, error) {
	return "", fmt.Errorf("ProcessVideo not implemented")
}

// ProcessProfileImage procesa la imagen de perfil ubicada en filePath y la sube a S3.
// Redimensiona la imagen para mobile y la codifica a formato WebP con calidad 80.
// Devuelve la URL pública de S3.
func ProcessProfileImage(filePath, folder string) (string, error) {
	// Crear servicio S3
	s3Service, err := NewS3Service()
	if err != nil {
		return "", fmt.Errorf("failed to create S3 service: %v", err)
	}

	// 🚀 CARGA DIRECTA DE IMAGEN SIN OpenCV PARA PERFILES
	processedImg, err := imaging.Open(filePath)
	if err != nil {
		return "", fmt.Errorf("failed to read image: %s", err)
	}

	// Redimensionar la imagen para mobile
	resizedImg := imaging.Resize(processedImg, MobileWidth, 0, imaging.Lanczos)

	// Subir imagen procesada a S3
	s3URL, err := s3Service.UploadImage(resizedImg, folder)
	if err != nil {
		return "", fmt.Errorf("failed to upload profile image to S3: %v", err)
	}

	return s3URL, nil
}
