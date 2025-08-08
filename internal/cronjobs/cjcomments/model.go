package cjcomments

import (
	"database/sql"
)

// CommentNotification representa una notificación de comentario pendiente.
type CommentNotification struct {
	NotificationID int64
	CommentID      int64
	CreatedAt      sql.NullTime
}

// CommentDetails contiene la información relevante del comentario y el clúster.
type CommentDetails struct {
	CommentID   int64
	ClusterID   int64
	CommentText string
	CommenterID int64
	SubcategoryName string
}

// Recipient representa un usuario que debe recibir la notificación.
type Recipient struct {
	AccountID   int64
	DeviceToken string
}
