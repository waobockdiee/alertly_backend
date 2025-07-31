// internal/common/push.go
package common

import "os"

// SendPush selecciona entre Expo Push Service o APNs directo
func SendPush(msg ExpoPushMessage, apnsNotif APNsNotification) error {
	// Si apnsClient está inicializado y estamos en producción, usa APNs
	if os.Getenv("APNS_ENV") == "production" && apnsClient != nil {
		return SendAPNsPush(apnsNotif)
	}
	// En cualquier otro caso (sandbox, falta de cert), envía a Expo
	return SendExpoPush(msg)
}
