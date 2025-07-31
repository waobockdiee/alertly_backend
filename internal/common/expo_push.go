// internal/common/push.go
package common

import (
	"bytes"
	"encoding/json"
	"fmt"
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

// APNsNotification encapsula los campos para APNs directo
type APNsNotification struct {
	DeviceToken string
	Topic       string
	Payload     *payload.Payload
}

var (
	expoEndpoint = "https://exp.host/--/api/v2/push/send"
	apnsClient   *apns2.Client
)

func init() {
	env := os.Getenv("APNS_ENV")
	p12Path := os.Getenv("APNS_P12_PATH")
	p12Pass := os.Getenv("APNS_P12_PASS")

	if p12Path != "" && p12Pass != "" {
		cert, err := certificate.FromP12File(p12Path, p12Pass)
		if err != nil {
			panic(fmt.Errorf("APNs cert load error (%s): %w", p12Path, err))
		}
		if env == "production" {
			apnsClient = apns2.NewClient(cert).Production()
			fmt.Println("✅ APNs client initialized in Production mode")
		} else {
			apnsClient = apns2.NewClient(cert).Development()
			fmt.Println("✅ APNs client initialized in Development (sandbox) mode")
		}
	} else {
		fmt.Println("⚠️ APNs client not configured (missing APNS_P12_PATH or APNS_P12_PASS)")
	}
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

// SendAPNsPush envía directamente a APNs (si apnsClient está inicializado)
func SendAPNsPush(n APNsNotification) error {
	if apnsClient == nil {
		return fmt.Errorf("APNs client not configured")
	}
	notification := &apns2.Notification{
		DeviceToken: n.DeviceToken,
		Topic:       n.Topic,
		Payload:     n.Payload,
	}
	res, err := apnsClient.Push(notification)
	if err != nil {
		return fmt.Errorf("APNs push error: %w", err)
	}
	if !res.Sent() {
		return fmt.Errorf("APNs push failed: %v", res.Reason)
	}
	return nil
}
