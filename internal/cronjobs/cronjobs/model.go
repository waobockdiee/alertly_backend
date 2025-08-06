package cronjobs

import "time"

// type Cluster struct {
// 	InclID int64 `db:"incl_id" json:"incl_id"`
// }

type Notification struct {
	ID          int64
	AccountID   int64
	ClusterID   int64
	CreatedAt   time.Time
	Processed   bool
	DeviceToken string
}

// Delivery registra a quién se envió la notificación
type Delivery struct {
	NotificationID int64
	ReferenceID    int64
	AccountID      int64
	SentAt         time.Time
}
