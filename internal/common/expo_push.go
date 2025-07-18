package common

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

// ExpoPushMessage represents the payload for Expo Push Service
type ExpoPushMessage struct {
	To    string                 `json:"to"`
	Title string                 `json:"title"`
	Body  string                 `json:"body"`
	Data  map[string]interface{} `json:"data,omitempty"`
}

// SendExpoPush sends a push notification via Expo Push Service
func SendExpoPush(msg ExpoPushMessage) error {
	payload, err := json.Marshal(msg)
	if err != nil {
		return fmt.Errorf("failed to marshal ExpoPushMessage: %w", err)
	}

	resp, err := http.Post(
		"https://exp.host/--/api/v2/push/send",
		"application/json",
		bytes.NewBuffer(payload),
	)
	if err != nil {
		return fmt.Errorf("error sending push to Expo: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("expo push failed with status: %s", resp.Status)
	}
	return nil
}
