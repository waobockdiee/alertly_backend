package cjinactivityreminder

import (
	"alertly/internal/cronjobs/shared"
	"context"
	"database/sql"
	"fmt"
	"time"
)

// Repository maneja la persistencia para el cronjob de inactividad.
type Repository struct {
	db *sql.DB
}

// NewRepository crea una nueva instancia de Repository.
func NewRepository(db *sql.DB) *Repository {
	return &Repository{db: db}
}

const (
	// inactivityThresholdDays define cuántos días de inactividad deben pasar para generar una notificación.
	inactivityThresholdDays = 7
	// inactivityType corresponde al tipo de notificación para inactividad
	inactivityType = "inactivity_reminder"
)

// GenerateInactivityNotifications inserta en la tabla notifications todos los usuarios
// que llevan más de inactivityThresholdDays días sin actividad y aún no tienen
// una notificación de inactividad pendiente.
func (r *Repository) GenerateInactivityNotifications() error {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// Título y mensaje que se enviará
	title := "We miss you at Alertly"
	message := "It’s been a while since you last logged in. Come back and check what’s new!"

	query := fmt.Sprintf(`
        INSERT INTO notifications (owner_account_id, title, message, type, created_at, must_be_processed)
        SELECT a.account_id, '%s', '%s', '%s', NOW(), 1
          FROM account a
          JOIN (
             SELECT account_id, MAX(created_at) as last_login
             FROM account_session_history
             GROUP BY account_id
          ) ash ON a.account_id = ash.account_id
         WHERE ash.last_login <= NOW() - INTERVAL '%d' DAY
           AND NOT EXISTS (
             SELECT 1 FROM notifications n
              WHERE n.owner_account_id = a.account_id
                AND n.type = '%s'
                AND n.must_be_processed = 1
           )
    `, title, message, inactivityType, inactivityThresholdDays, inactivityType)

	_, err := r.db.ExecContext(ctx, query)
	return err
}

// FetchPending obtiene un lote (batch) de notificaciones pendientes.
// Utiliza keyset pagination para maximizar performance.
func (r *Repository) FetchPending(limit int, lastID int64) ([]Notification, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	sqlQuery := `
        SELECT n.noti_id, n.owner_account_id, dt.device_token, n.title, n.message
          FROM notifications n
          JOIN device_tokens dt ON n.owner_account_id = dt.account_id
         WHERE n.must_be_processed = 1
           AND n.type = $1
           AND n.noti_id > $2
         ORDER BY n.noti_id
         LIMIT $3`

	rows, err := r.db.QueryContext(ctx, sqlQuery, inactivityType, lastID, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var notis []Notification
	for rows.Next() {
		var n Notification
		if err := rows.Scan(&n.NotificationID, &n.AccountID, &n.DeviceToken, &n.Title, &n.Message); err != nil {
			return nil, err
		}
		notis = append(notis, n)
	}
	return notis, rows.Err()
}

// InsertDeliveries inserta en bloque los registros en notification_deliveries.
func (r *Repository) InsertDeliveries(deliveries []shared.Delivery) error {
	return shared.InsertDeliveries(r.db, deliveries)
}

// MarkProcessed actualiza must_be_processed a 0 y setea sent_at para las notificaciones enviadas.
func (r *Repository) MarkProcessed(ids []int64) error {
	return shared.MarkItemsAsProcessed(r.db, "notifications", "noti_id", ids)
}
