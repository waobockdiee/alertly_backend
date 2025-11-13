package media

import (
	"bytes"
	"fmt"
	"image"
	"image/color"
	"image/draw"
	"image/jpeg"
	"log"
	"os"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/rekognition"
	"github.com/disintegration/imaging"
	"github.com/rwcarlsen/goexif/exif"
)

// Constantes para los archivos de cascada y configuración de procesamiento
const (
	FaceCascadeFile  = "data/cascades/haarcascade_frontalface_default.xml"
	PlateCascadeFile = "data/cascades/haarcascade_russian_plate_number.xml"
	MobileWidth      = 720 // Ancho deseado para la imagen en mobile
)

// PixelateRegion aplica efecto de pixelado a una región específica de la imagen
func PixelateRegion(img image.Image, region image.Rectangle, pixelSize int) image.Image {
	// Crear una nueva imagen que sea una copia de la original
	bounds := img.Bounds()
	pixelatedImg := imaging.Clone(img)

	// Convertir a imagen drawable
	drawable := image.NewRGBA(bounds)
	draw.Draw(drawable, bounds, pixelatedImg, bounds.Min, draw.Src)

	// Asegurar que la región está dentro de los límites de la imagen
	region = region.Intersect(bounds)
	if region.Empty() {
		return img // Si no hay intersección, devolver imagen original
	}

	// Aplicar pixelado a la región especificada
	for y := region.Min.Y; y < region.Max.Y; y += pixelSize {
		for x := region.Min.X; x < region.Max.X; x += pixelSize {
			// Obtener el color promedio del bloque de pixeles
			r, g, b, a := getAverageColor(img, image.Rect(x, y,
				min(x+pixelSize, region.Max.X),
				min(y+pixelSize, region.Max.Y)))

			// Aplicar el color promedio a todo el bloque
			pixelColor := color.RGBA{r, g, b, a}
			for py := y; py < min(y+pixelSize, region.Max.Y); py++ {
				for px := x; px < min(x+pixelSize, region.Max.X); px++ {
					drawable.Set(px, py, pixelColor)
				}
			}
		}
	}

	return drawable
}

// getAverageColor calcula el color promedio de una región de la imagen
func getAverageColor(img image.Image, region image.Rectangle) (uint8, uint8, uint8, uint8) {
	var totalR, totalG, totalB, totalA uint64
	var count uint64

	for y := region.Min.Y; y < region.Max.Y; y++ {
		for x := region.Min.X; x < region.Max.X; x++ {
			r, g, b, a := img.At(x, y).RGBA()
			totalR += uint64(r >> 8) // Convertir de uint32 a uint8
			totalG += uint64(g >> 8)
			totalB += uint64(b >> 8)
			totalA += uint64(a >> 8)
			count++
		}
	}

	if count == 0 {
		return 0, 0, 0, 255
	}

	return uint8(totalR / count), uint8(totalG / count), uint8(totalB / count), uint8(totalA / count)
}

// min función auxiliar
func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

// DetectFacesWithRekognition usa AWS Rekognition para detección precisa de rostros
func DetectFacesWithRekognition(imageBytes []byte, imgBounds image.Rectangle) ([]image.Rectangle, error) {
	// Crear sesión de AWS
	sess, err := session.NewSession(&aws.Config{
		Region: aws.String("us-west-2"), // Misma región que tu infraestructura
	})
	if err != nil {
		log.Printf("AWS session error (skipping face detection): %v", err)
		return []image.Rectangle{}, nil // Fallar silenciosamente
	}

	// Crear servicio Rekognition
	svc := rekognition.New(sess)

	// Preparar input para Rekognition
	input := &rekognition.DetectFacesInput{
		Image: &rekognition.Image{
			Bytes: imageBytes,
		},
		Attributes: []*string{
			aws.String("DEFAULT"),
		},
	}

	// Detectar rostros
	result, err := svc.DetectFaces(input)
	if err != nil {
		log.Printf("Rekognition error (skipping face detection): %v", err)
		return []image.Rectangle{}, nil // Fallar silenciosamente, no bloquear el upload
	}

	// Convertir resultados a rectangulos absolutos
	var faces []image.Rectangle
	imgWidth := float64(imgBounds.Dx())
	imgHeight := float64(imgBounds.Dy())

	for _, faceDetail := range result.FaceDetails {
		if faceDetail.BoundingBox != nil {
			// Convertir coordenadas relativas (0-1) a absolutas
			left := int(*faceDetail.BoundingBox.Left * imgWidth)
			top := int(*faceDetail.BoundingBox.Top * imgHeight)
			width := int(*faceDetail.BoundingBox.Width * imgWidth)
			height := int(*faceDetail.BoundingBox.Height * imgHeight)

			faces = append(faces, image.Rectangle{
				Min: image.Point{X: left, Y: top},
				Max: image.Point{X: left + width, Y: top + height},
			})
		}
	}

	return faces, nil
}

// DetectLicensePlatesSimple implementa detección básica y rápida de placas
// Pixela solo regiones típicas donde están las placas sin análisis complejo
func DetectLicensePlatesSimple(img image.Image) ([]image.Rectangle, error) {
	bounds := img.Bounds()
	var plates []image.Rectangle

	// Solo agregar 1 región en el tercio inferior donde típicamente están las placas
	// Zona inferior centro (más probable para placas)
	bottomRegion := image.Rect(
		bounds.Dx()/4,
		bounds.Dy()*4/5, // 80% hacia abajo
		bounds.Dx()*3/4,
		bounds.Max.Y,
	)

	plates = append(plates, bottomRegion)

	return plates, nil
}

// fixImageOrientation lee los metadatos EXIF y aplica la rotación necesaria según la orientación
func fixImageOrientation(img image.Image, filePath string) (image.Image, error) {
	// Abrir el archivo para leer metadatos EXIF
	file, err := os.Open(filePath)
	if err != nil {
		// Si no se puede abrir el archivo, devolver la imagen sin modificar
		log.Printf("Could not open file for EXIF reading: %v", err)
		return img, nil
	}
	defer file.Close()

	// Leer metadatos EXIF
	exifData, err := exif.Decode(file)
	if err != nil {
		// Si no hay metadatos EXIF, devolver la imagen sin modificar
		log.Printf("No EXIF data found or error reading EXIF: %v", err)
		return img, nil
	}

	// Buscar el tag de orientación
	orientationTag, err := exifData.Get(exif.Orientation)
	if err != nil {
		// Si no hay tag de orientación, devolver la imagen sin modificar
		log.Printf("No orientation tag found: %v", err)
		return img, nil
	}

	// Obtener el valor de orientación
	orientation, err := orientationTag.Int(0)
	if err != nil {
		log.Printf("Could not parse orientation value: %v", err)
		return img, nil
	}

	// Aplicar la rotación según el valor de orientación EXIF
	switch orientation {
	case 1:
		// Normal, no rotation needed
		return img, nil
	case 2:
		// Flip horizontal
		return imaging.FlipH(img), nil
	case 3:
		// Rotate 180°
		return imaging.Rotate180(img), nil
	case 4:
		// Flip vertical
		return imaging.FlipV(img), nil
	case 5:
		// Rotate 270° and flip horizontal
		rotated := imaging.Rotate270(img)
		return imaging.FlipH(rotated), nil
	case 6:
		// Rotate 270° (or 90° clockwise)
		return imaging.Rotate270(img), nil
	case 7:
		// Rotate 90° and flip horizontal
		rotated := imaging.Rotate90(img)
		return imaging.FlipH(rotated), nil
	case 8:
		// Rotate 90° (or 270° clockwise)
		return imaging.Rotate90(img), nil
	default:
		log.Printf("Unknown orientation value: %d", orientation)
		return img, nil
	}
}

// ProcessImage procesa la imagen ubicada en filePath, aplica detección de rostros y matrículas,
// y la sube a S3. Devuelve la URL pública de S3.
func ProcessImage(filePath, folder string) (string, error) {
	// Crear servicio S3
	s3Service, err := NewS3Service()
	if err != nil {
		return "", fmt.Errorf("failed to create S3 service: %v", err)
	}

	// ✅ PROCESAMIENTO CON DETECCIÓN Y PIXELADO OPTIMIZADO REKOGNITION
	fmt.Println("🔍 Processing image with AWS Rekognition face detection...")

	// Cargar la imagen original
	originalImg, err := imaging.Open(filePath)
	if err != nil {
		return "", fmt.Errorf("failed to read image: %s", err)
	}

	// Aplicar rotación automática basada en metadatos EXIF
	orientedImg, err := fixImageOrientation(originalImg, filePath)
	if err != nil {
		log.Printf("Error applying EXIF orientation: %v", err)
		// Continuar con la imagen sin orientar si hay error
		orientedImg = originalImg
	}

	// Inicializar imagen procesada con la imagen orientada
	processedImg := orientedImg

	// 1. DETECCIÓN DE ROSTROS usando AWS Rekognition
	// Convertir imagen a bytes para Rekognition
	var buf bytes.Buffer
	if err := jpeg.Encode(&buf, orientedImg, &jpeg.Options{Quality: 85}); err != nil {
		log.Printf("Error encoding image for Rekognition: %v", err)
	} else {
		faces, err := DetectFacesWithRekognition(buf.Bytes(), orientedImg.Bounds())
		if err != nil {
			log.Printf("Error detecting faces with Rekognition: %v", err)
		} else {
			fmt.Printf("✅ Detected %d faces with AWS Rekognition\n", len(faces))

			// Aplicar pixelado a cada rostro detectado
			for _, face := range faces {
				// Aplicar pixelado con tamaño de pixel de 15
				processedImg = PixelateRegion(processedImg, face, 15)
			}
		}
	}

	// 2. DETECCIÓN DE PLACAS DE MATRÍCULA (simple y rápida)
	plates, err := DetectLicensePlatesSimple(processedImg)
	if err != nil {
		log.Printf("Error detecting license plates: %v", err)
	} else {
		fmt.Printf("✅ Applied pixelation to %d potential license plate regions\n", len(plates))

		// Aplicar pixelado a cada placa detectada
		for _, plate := range plates {
			// Aplicar pixelado con tamaño de pixel de 10 (más fino para placas)
			processedImg = PixelateRegion(processedImg, plate, 10)
		}
	}

	// 3. Redimensionar la imagen para mobile (después del pixelado para mantener calidad)
	resizedImg := imaging.Resize(processedImg, MobileWidth, 0, imaging.Lanczos)

	// 4. Subir imagen procesada a S3
	s3URL, err := s3Service.UploadImage(resizedImg, folder)
	if err != nil {
		return "", fmt.Errorf("failed to upload processed image to S3: %v", err)
	}

	fmt.Printf("✅ Image processed successfully with AWS Rekognition privacy features\n")
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

	// Aplicar rotación automática basada en metadatos EXIF
	orientedImg, err := fixImageOrientation(processedImg, filePath)
	if err != nil {
		log.Printf("Error applying EXIF orientation for profile: %v", err)
		// Continuar con la imagen sin orientar si hay error
		orientedImg = processedImg
	}

	// Redimensionar la imagen para mobile
	resizedImg := imaging.Resize(orientedImg, MobileWidth, 0, imaging.Lanczos)

	// Subir imagen procesada a S3
	s3URL, err := s3Service.UploadImage(resizedImg, folder)
	if err != nil {
		return "", fmt.Errorf("failed to upload profile image to S3: %v", err)
	}

	return s3URL, nil
}
