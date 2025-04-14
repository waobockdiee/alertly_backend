/*------------------------------------------------------------------------------
internal/common/notification.go depende de esto

Los tipos de notificaciones tienen la intencion de tanto enviar notificaciones a un usuario. Y si MustBeProcessed lo aplica. Debe enviar masivamente.

MustBeProcessed = true => entonces debe ser enviado masivamente a los respectivos usuarios involucrados (TIENE QUE PROCESARLO EL CRONJOB y enviara y se comunica con un tercero que hara el notification push)

MustBeProcessed = false => solo sera enviado a una sola persona(autor de la accion)

MustSendPush = true => debe ser enviado como notification push
MustSendInApp = true => debe mostrarse como popup o toast
------------------------------------------------------------------------------*/

package alerts

import "time"

type Alert struct {
	AcnoID          int64     `db:"acno_id" json:"acno_id"`
	AccountID       int64     `db:"account_id" json:"account_id"`
	Title           string    `db:"title" json:"title"`
	Message         string    `db:"message" json:"message"`
	Type            string    `db:"type" json:"type"`
	IsRead          bool      `db:"is_read" json:"is_read"`
	Link            string    `db:"link" json:"link"`
	CreatedAt       time.Time `db:"created_at" json:"created_at"`
	UpdatedAt       time.Time `db:"updated_at" json:"updated_at"`
	SentAt          time.Time `db:"sent_at" json:"sent_at"`
	MustSendPush    bool      `db:"must_send_as_notification_push" json:"must_send_as_notification_push"`
	MustSendInApp   bool      `db:"must_send_as_notification" json:"must_send_as_notification"`
	MustBeProcessed bool      `db:"must_be_processed" json:"must_be_processed"`
	ErrorMessage    string    `db:"error_message" json:"error_message "`
	RetryCount      int32     `db:"retry_count" json:"retry_count"`
	ReferenceID     int64     `db:"reference_id" json:"reference_id"`
}
