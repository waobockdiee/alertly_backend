package common

import (
	"log"
	"os"

	"github.com/sideshow/apns2"
	"github.com/sideshow/apns2/certificate"
)

var APNSClient *apns2.Client

func init() {
	env := os.Getenv("APNS_ENV") // debe ser "production" o "development"
	p12Path := os.Getenv("APNS_P12_PATH")
	p12Pass := os.Getenv("APNS_P12_PASS")

	// Solo en producción y si tengo la ruta y contraseña, cargo el certificado
	if env == "production" {
		if p12Path == "" || p12Pass == "" {
			log.Fatalf("APNs production mode requires APNS_P12_PATH and APNS_P12_PASS")
		}
		cert, err := certificate.FromP12File(p12Path, p12Pass)
		if err != nil {
			log.Fatalf("APNs cert load error (%s): %v", p12Path, err)
		}
		APNSClient = apns2.NewClient(cert).Production()
		log.Println("APNs client initialized in Production mode")
		return
	}

	// En cualquier otro caso (sandbox / local), no inicializo APNSClient
	log.Printf("ℹSkipping APNs init (APNS_ENV=%s)", env)
}
