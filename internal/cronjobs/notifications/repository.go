package notifications

import (
	"database/sql"
	"log"
)

type Repository interface {
	GetUnprocessedNotificationsPush() ([]Notification, error)
	BatchSaveNewNotificationDeliveries(nd []NotificationDelivery) error
	SaveNotificationDelivery(nd NotificationDelivery) error
	UpdateNotificationAsProcessed(notiID int64) error
	GetProcessWelcomeToAppAccounts(n Notification) ([]Account, error)
	GetDB() *sql.DB
}

type pgRepository struct {
	db *sql.DB
}

func NewRepository(db *sql.DB) Repository {
	return &pgRepository{db: db}
}

// GetUnprocessedNotificationsPush obtiene las notificaciones pendientes de procesar.
func (r *pgRepository) GetUnprocessedNotificationsPush() ([]Notification, error) {
	var notifications []Notification
	query := `
	SELECT
		t1.noti_id, t1.owner_account_id, t1.title, t1.message, t1.type, t1.link,
		t1.created_at, t1.must_send_as_notification_push, t1.must_send_as_notification,
		t1.must_be_processed, t1.retry_count, t1.reference_id, t2.nickname, COALESCE(t2.thumbnail_url, '') as thumbnail_url
	FROM notifications t1
	INNER JOIN account t2 ON t1.owner_account_id = t2.account_id
	WHERE t1.must_be_processed = 1`

	rows, err := r.db.Query(query)
	if err != nil {
		log.Printf("Error ejecutando query GetUnprocessedNotificationsPush: %v", err)
		return notifications, err
	}
	defer rows.Close()

	for rows.Next() {
		var n Notification
		if err := rows.Scan(
			&n.NotiID, &n.AccountID, &n.Title, &n.Message, &n.Type, &n.Link,
			&n.CreatedAt, &n.MustSendPush, &n.MustSendInApp, &n.MustBeProcessed,
			&n.RetryCount, &n.ReferenceID, &n.Nickname, &n.ThumbnailURL,
		); err != nil {
			log.Printf("Error escaneando fila en GetUnprocessedNotificationsPush: %v", err)
			continue
		}
		notifications = append(notifications, n)
	}
	return notifications, nil
}

// BatchSaveNewNotificationDeliveries realiza una inserción en batch de las notification deliveries.
func (r *pgRepository) BatchSaveNewNotificationDeliveries(nd []NotificationDelivery) error {
	if len(nd) == 0 {
		return nil
	}

	// Usar una transacción para insertar en batch de manera eficiente
	tx, err := r.db.Begin()
	if err != nil {
		log.Printf("Error starting transaction: %v", err)
		return err
	}
	defer tx.Rollback()

	stmt, err := tx.Prepare(`INSERT INTO notification_deliveries (to_account_id, noti_id, title, message) VALUES ($1, $2, $3, $4)`)
	if err != nil {
		log.Printf("Error preparing statement: %v", err)
		return err
	}
	defer stmt.Close()

	for _, delivery := range nd {
		_, err := stmt.Exec(delivery.ToAccountID, delivery.NotiID, delivery.Title, delivery.Message)
		if err != nil {
			log.Printf("Error executing statement for delivery: %v", err)
			return err
		}
	}

	if err := tx.Commit(); err != nil {
		log.Printf("Error committing transaction: %v", err)
		return err
	}

	return nil
}

// SaveNotificationDelivery guarda una sola notification delivery
func (r *pgRepository) SaveNotificationDelivery(nd NotificationDelivery) error {
	query := "INSERT INTO notification_deliveries (to_account_id, noti_id, title, message) VALUES ($1, $2, $3, $4)"
	_, err := r.db.Exec(query, nd.ToAccountID, nd.NotiID, nd.Title, nd.Message)
	if err != nil {
		log.Printf("Error saving individual notification delivery: %v", err)
		return err
	}
	return nil
}

func (r *pgRepository) UpdateNotificationAsProcessed(notiID int64) error {
	query := "UPDATE notifications SET must_be_processed = 0 WHERE noti_id = $1"
	_, err := r.db.Exec(query, notiID)
	if err != nil {
		log.Printf("Error actualizando notificación (ID: %d) como procesada: %v", notiID, err)
	}
	return err
}

func (r *pgRepository) GetProcessWelcomeToAppAccounts(n Notification) ([]Account, error) {
	var accounts []Account

	// Suponiendo que para la notificación 'welcome_to_app' se quieren obtener todas las cuentas
	// excepto el dueño (podrías ajustar la condición según la lógica del negocio).
	query := `
	SELECT
		a.account_id, a.email, a.nickname, COALESCE(a.thumbnail_url, '') as thumbnail_url
	FROM account a
	WHERE a.account_id != $1`
	rows, err := r.db.Query(query, n.AccountID)
	if err != nil {
		log.Printf("Error ejecutando query GetProcessWelcomeToAppAccounts: %v", err)
		return accounts, err
	}
	defer rows.Close()

	for rows.Next() {
		var account Account
		if err := rows.Scan(&account.AccountID, &account.Email, &account.Nickname, &account.Thumbnail); err != nil {
			log.Printf("Error escaneando fila en GetProcessWelcomeToAppAccounts: %v", err)
			continue
		}
		accounts = append(accounts, account)
	}
	return accounts, nil
}

// GetDB returns the database connection
func (r *pgRepository) GetDB() *sql.DB {
	return r.db
}
