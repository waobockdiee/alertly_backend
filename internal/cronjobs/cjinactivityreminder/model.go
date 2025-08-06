package cjinactivityreminder

// Notification representa la estructura de datos para una notificación pendiente.
type Notification struct {
	NotificationID int64
	AccountID      int64
	DeviceToken    string
	Title          string
	Message        string
}