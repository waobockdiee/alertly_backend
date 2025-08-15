package cjincidentupdate

import (
	"database/sql"
)

// IncidentUpdateNotification representa una notificación de actualización de incidente pendiente.
type IncidentUpdateNotification struct {
	NotificationID int64
	ClusterID      int64
	ReporterID     int64
	CreatedAt      sql.NullTime
}

// ClusterDetails contiene la información relevante del cluster actualizado.
type ClusterDetails struct {
	ClusterID       int64
	SubcategoryName string
	Description     string
	City            string
	ReporterID      int64
}

// Recipient representa un usuario que debe recibir la notificación.
type Recipient struct {
	AccountID   int64
	DeviceToken string
}

