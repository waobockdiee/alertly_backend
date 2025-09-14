// internal/common/push.go
package common

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/sideshow/apns2"
	"github.com/sideshow/apns2/certificate"
	"github.com/sideshow/apns2/payload"
)

// ExpoPushMessage representa el payload para Expo Push Service
type ExpoPushMessage struct {
	To    string                 `json:"to"`
	Title string                 `json:"title"`
	Body  string                 `json:"body"`
	Data  map[string]interface{} `json:"data,omitempty"`
}

var (
	// Expo endpoint
	expoEndpoint = "https://exp.host/--/api/v2/push/send"

	// APNs client e información de tópico
	APNSClient *apns2.Client
	apnsTopic  string
)

func init() {
	// Carga configuración de entorno
	env := os.Getenv("APNS_ENV") // "production" o "development"
	p12Pass := os.Getenv("APNS_P12_PASS")
	apnsTopic = os.Getenv("APNS_TOPIC")

	// Inicializa APNs solo en producción
	if env == "production" {
		// ✅ AWS Lambda: Usar certificado desde variable de entorno base64
		p12Base64 := os.Getenv("APNS_P12_BASE64")
		if p12Base64 == "" || p12Pass == "" || apnsTopic == "" {
			log.Printf("⚠️ APNs disabled: missing APNS_P12_BASE64, APNS_P12_PASS or APNS_TOPIC")
			return
		}

		// Decodificar certificado de base64
		certData, err := base64.StdEncoding.DecodeString(p12Base64)
		if err != nil {
			log.Printf("⚠️ APNs cert decode error: %v", err)
			return
		}

		cert, err := certificate.FromP12Bytes(certData, p12Pass)
		if err != nil {
			log.Printf("⚠️ APNs cert load error: %v", err)
			return
		}
		APNSClient = apns2.NewClient(cert).Production()
		log.Println("✅ APNs client initialized in Production mode")
		return
	}

	// En sandbox o sin configuración, no inicializamos APNSClient
	log.Printf("ℹ️ Skipping APNs init (APNS_ENV=%s)", env)
}

// SendExpoPush envía vía Expo Push Service
func SendExpoPush(msg ExpoPushMessage) error {
	payloadBytes, err := json.Marshal(msg)
	if err != nil {
		return fmt.Errorf("failed to marshal ExpoPushMessage: %w", err)
	}

	resp, err := http.Post(expoEndpoint, "application/json", bytes.NewBuffer(payloadBytes))
	if err != nil {
		return fmt.Errorf("error sending push to Expo: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("expo push failed with status: %s", resp.Status)
	}
	return nil
}

// SendAPNsPush envía directamente a APNs
type APNsNotification struct {
	DeviceToken string
	Payload     *payload.Payload
}

func SendAPNsPush(n APNsNotification) error {
	if APNSClient == nil {
		return fmt.Errorf("APNs client not configured")
	}
	notification := &apns2.Notification{
		DeviceToken: n.DeviceToken,
		Topic:       apnsTopic,
		Payload:     n.Payload,
	}
	res, err := APNSClient.Push(notification)
	if err != nil {
		return fmt.Errorf("APNs push error: %w", err)
	}
	if !res.Sent() {
		return fmt.Errorf("APNs push failed: %v", res.Reason)
	}
	return nil
}

// SendPush selecciona entre Expo o APNs según configuración
// - expoMsg: argumentos para Expo Push Service
// - deviceToken: token de destino (igual para ambos servicios)
// - apnsPayload: payload para APNs directo
func SendPush(expoMsg ExpoPushMessage, deviceToken string, apnsPayload *payload.Payload) error {
	// Inyecta el token en la petición Expo
	expoMsg.To = deviceToken

	// Si estamos en producción y APNSClient listo, enviamos por APNs
	if os.Getenv("APNS_ENV") == "production" && APNSClient != nil {
		return SendAPNsPush(APNsNotification{DeviceToken: deviceToken, Payload: apnsPayload})
	}

	// En cualquier otro caso, enviamos por Expo
	return SendExpoPush(expoMsg)
}
