package media

import (
	"fmt"
	"image"
	"os"
	"path/filepath"
	"time"

	"github.com/chai2010/webp"
	"github.com/disintegration/imaging"
	"gocv.io/x/gocv"
)

// Constantes para los archivos de cascada y configuración de procesamiento
const (
	FaceCascadeFile  = "data/cascades/haarcascade_frontalface_default.xml"
	PlateCascadeFile = "data/cascades/haarcascade_russian_plate_number.xml"
	MobileWidth      = 720 // Ancho deseado para la imagen en mobile
)

// PixelateROI aplica un efecto pixelado a la región de interés (ROI) en el Mat.
// pixelSize determina el nivel de pixelación (por ejemplo, 10).
func PixelateROI(mat *gocv.Mat, r image.Rectangle, pixelSize int) {
	// Extraer la ROI
	roi := mat.Region(r)
	defer roi.Close()

	origWidth := roi.Cols()
	origHeight := roi.Rows()

	// Calcular dimensiones reducidas
	smallWidth := origWidth / pixelSize
	smallHeight := origHeight / pixelSize
	if smallWidth < 1 {
		smallWidth = 1
	}
	if smallHeight < 1 {
		smallHeight = 1
	}

	// Crear un Mat temporal para la versión pixelada
	pixelated := gocv.NewMat()
	defer pixelated.Close()

	// Reducir el tamaño: Interpolación lineal
	gocv.Resize(roi, &pixelated, image.Pt(smallWidth, smallHeight), 0, 0, gocv.InterpolationLinear)
	// Escalar nuevamente al tamaño original usando interpolación de vecinos cercanos para crear el efecto pixelado
	gocv.Resize(pixelated, &roi, image.Pt(origWidth, origHeight), 0, 0, gocv.InterpolationNearestNeighbor)
}

// ProcessImage procesa la imagen ubicada en filePath y la guarda en outputFolder.
// Aplica detección de rostros y matrículas, a las que se les aplica un efecto de pixelado.
// Luego, redimensiona la imagen para mobile y la codifica a formato WebP con calidad 80.
// Devuelve la ruta completa del archivo procesado o un error.
func ProcessImage(filePath, outputFolder string) (string, error) {
	// Cargar la imagen en un Mat usando gocv
	img := gocv.IMRead(filePath, gocv.IMReadColor)
	if img.Empty() {
		return "", fmt.Errorf("failed to read image: %s", filePath)
	}
	defer img.Close()

	// Cargar la cascada para detección de rostros
	faceCascade := gocv.NewCascadeClassifier()
	defer faceCascade.Close()
	if !faceCascade.Load(FaceCascadeFile) {
		return "", fmt.Errorf("error loading face cascade file: %s", FaceCascadeFile)
	} else {
		fmt.Println("Face cascade loaded successfully.")
	}

	// Detectar rostros usando parámetros ajustados (ajusta según necesidad)
	faceRects := faceCascade.DetectMultiScaleWithParams(img, 1.1, 3, 0, image.Pt(30, 30), image.Pt(0, 0))
	fmt.Printf("Detected %d faces.\n", len(faceRects))
	for _, r := range faceRects {
		// Aplicar pixelado en lugar de blur; pixelSize ajustado a 10 (puedes modificarlo)
		PixelateROI(&img, r, 20)
	}

	// Cargar la cascada para matrículas
	plateCascade := gocv.NewCascadeClassifier()
	defer plateCascade.Close()
	if plateCascade.Load(PlateCascadeFile) {
		fmt.Println("Plate cascade loaded successfully.")
		plateRects := plateCascade.DetectMultiScale(img)
		fmt.Printf("Detected %d plates.\n", len(plateRects))
		for _, r := range plateRects {
			PixelateROI(&img, r, 20)
		}
	} else {
		fmt.Println("No se pudo cargar la cascada de matrículas, omitiendo detección de placas")
	}

	// Convertir el Mat procesado a image.Image
	processedImg, err := img.ToImage()
	if err != nil {
		return "", fmt.Errorf("error converting Mat to image: %v", err)
	}

	// Redimensionar la imagen para mobile
	resizedImg := imaging.Resize(processedImg, MobileWidth, 0, imaging.Lanczos)

	// Generar un nuevo nombre de archivo con timestamp y extensión .webp
	timestamp := time.Now().UnixNano()
	newFilename := fmt.Sprintf("alerty_%d.webp", timestamp)
	outputPath := filepath.Join(outputFolder, newFilename)

	// Crear la carpeta de salida si no existe
	if _, err := os.Stat(outputFolder); os.IsNotExist(err) {
		if err := os.MkdirAll(outputFolder, os.ModePerm); err != nil {
			return "", fmt.Errorf("failed to create output folder: %v", err)
		}
	}

	// Crear el archivo de salida
	outFile, err := os.Create(outputPath)
	if err != nil {
		return "", fmt.Errorf("failed to create output file: %v", err)
	}
	defer outFile.Close()

	// Codificar la imagen en formato WebP con calidad 80
	if err := webp.Encode(outFile, resizedImg, &webp.Options{Quality: 80}); err != nil {
		return "", fmt.Errorf("failed to encode image to webp: %v", err)
	}

	ipserver := os.Getenv("IP_SERVER")
	fullurl := "http://" + ipserver + ":8080/" + outputPath

	return fullurl, nil
}

// ProcessVideo es una función placeholder para procesar videos.
func ProcessVideo(filePath, outputFolder string) (string, error) {
	return "", fmt.Errorf("ProcessVideo not implemented")
}
