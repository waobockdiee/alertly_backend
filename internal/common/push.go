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
	"github.com/sideshow/apns2/payload"
	"github.com/sideshow/apns2/token"
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
	apnsTopic = os.Getenv("APNS_TOPIC")

	// Inicializa APNs solo en producción
	if env == "production" {
		// ✅ Usar APNs Auth Key (.p8) - más robusto y no expira
		authKeyBase64 := os.Getenv("APNS_AUTH_KEY")
		keyID := os.Getenv("APNS_KEY_ID")
		teamID := os.Getenv("APNS_TEAM_ID")

		if authKeyBase64 == "" || keyID == "" || teamID == "" || apnsTopic == "" {
			log.Printf("⚠️ APNs disabled: missing APNS_AUTH_KEY, APNS_KEY_ID, APNS_TEAM_ID or APNS_TOPIC")
			return
		}

		// Decodificar auth key de base64
		authKeyData, err := base64.StdEncoding.DecodeString(authKeyBase64)
		if err != nil {
			log.Printf("⚠️ APNs auth key decode error: %v", err)
			return
		}

		authKey, err := token.AuthKeyFromBytes(authKeyData)
		if err != nil {
			log.Printf("⚠️ APNs auth key load error: %v", err)
			return
		}

		apnsToken := &token.Token{
			AuthKey: authKey,
			KeyID:   keyID,
			TeamID:  teamID,
		}

		APNSClient = apns2.NewTokenClient(apnsToken).Production()
		log.Printf("✅ APNs client initialized in Production mode (Key ID: %s, Team ID: %s)", keyID, teamID)
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

	// DEBUG: Log payload
	payloadJSON, _ := n.Payload.MarshalJSON()
	log.Printf("📤 Sending APNs push to %s... with payload: %s", n.DeviceToken[:20], string(payloadJSON))

	res, err := APNSClient.Push(notification)
	if err != nil {
		return fmt.Errorf("APNs push error: %w", err)
	}
	if !res.Sent() {
		return fmt.Errorf("APNs push failed: %v", res.Reason)
	}
	log.Printf("✅ APNs push sent successfully to %s...", n.DeviceToken[:20])
	return nil
}

// SendPush selecciona entre Expo o APNs según el tipo de token
// - expoMsg: argumentos para Expo Push Service
// - deviceToken: token de destino
// - apnsPayload: payload para APNs directo (solo para tokens nativos)
func SendPush(expoMsg ExpoPushMessage, deviceToken string, apnsPayload *payload.Payload) error {
	// Inyecta el token en la petición Expo
	expoMsg.To = deviceToken

	// ✅ DETECCIÓN DE TIPO DE TOKEN:
	// - ExponentPushToken[...] = Token de Expo → Usar Expo Push Service
	// - Token hex de 64 chars = Token nativo iOS → Usar APNs directo
	isExpoToken := len(deviceToken) > 18 && deviceToken[:18] == "ExponentPushToken["

	// Si es un token de Expo, SIEMPRE usar Expo Push Service
	if isExpoToken {
		return SendExpoPush(expoMsg)
	}

	// Si es un token nativo de iOS Y estamos en producción con APNs configurado
	if os.Getenv("APNS_ENV") == "production" && APNSClient != nil {
		return SendAPNsPush(APNsNotification{DeviceToken: deviceToken, Payload: apnsPayload})
	}

	// Fallback: intentar enviar por Expo
	return SendExpoPush(expoMsg)
}
