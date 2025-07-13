/*------------------------------------------------------------------------------
internal/common/notification.go depende de esto

Los tipos de notificaciones tienen la intencion de tanto enviar notificaciones a un usuario. Y si MustBeProcessed lo aplica. Debe enviar masivamente.

MustBeProcessed = true => entonces debe ser enviado masivamente a los respectivos usuarios involucrados (TIENE QUE PROCESARLO EL CRONJOB y enviara y se comunica con un tercero que hara el notification push)

MustBeProcessed = false => solo sera enviado a una sola persona(autor de la accion)

MustSendPush = true => debe ser enviado como notification push
MustSendInApp = true => debe mostrarse como popup o toast
------------------------------------------------------------------------------*/

package notifications

import "time"

type Notification struct {
	NotiID                     int64     `db:"noti_id" json:"noti_id"`
	OwnerAccountID             int64     `db:"owner_account_id" json:"owner_account_id"`
	Title                      string    `db:"title" json:"title"`
	Message                    string    `db:"message" json:"message"`
	Type                       string    `db:"type" json:"type"`
	Link                       string    `db:"link" json:"link"`
	CreatedAt                  time.Time `db:"created_at" json:"created_at"`
	SentAt                     time.Time `db:"sent_at" json:"sent_at"`
	MustSendAsNotificationPush bool      `db:"must_send_as_notification_push" json:"must_send_as_notification_push"`
	MustSendAsNotification     bool      `db:"must_send_as_notification" json:"must_send_as_notification"`
	MustBeProcessed            bool      `db:"must_be_processed" json:"must_be_processed"`
	ErrorMesssage              string    `db:"error_message" json:"error_message"`
	RetryCount                 int       `db:"retry_count" json:"retry_count"`
	ReferenceID                int64     `db:"reference_id" json:"reference_id"`
}
