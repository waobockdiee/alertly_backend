package cjnewcluster

import "time"

// Notification representa una fila pendiente de procesar
type Notification struct {
	ID        int64
	ClusterID int64
	CreatedAt time.Time
	Processed bool
}

// Delivery registra a quién se envió la notificación
type Delivery struct {
	NotificationID int64
	ReferenceID    int64
	AccountID      int64
	SentAt         time.Time
}
