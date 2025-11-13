package notifications

import (
	"alertly/internal/common"
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"time"
)

type Repository interface {
	Save(n Notification) (int64, error)
	SaveDeviceToken(accountID int64, token string) error
	DeleteDeviceToken(accountID int64, deviceToken string) error
	GetNotifications(accountID int64, limit, offset int) ([]NotificationDelivery, error)
	GetUnreadCount(accountID int64) (int64, error)
	MarkAsRead(accountID, notificationID int64) error
	MarkAllAsRead(accountID int64) error
	DeleteNotification(accountID, notificationID int64) error
}

type mysqlRepository struct {
	db *sql.DB
}

// NotificationDelivery representa una notificación entregada a un usuario
type NotificationDelivery struct {
	NodeID      int64           `db:"node_id" json:"node_id"`
	CreatedAt   common.NullTime `db:"created_at" json:"created_at"`
	IsRead      sql.NullInt64   `db:"is_read" json:"is_read"`
	ToAccountID int64           `db:"to_account_id" json:"to_account_id"`
	NotiID      int64           `db:"noti_id" json:"noti_id"`
	Title       string          `db:"title" json:"title"`
	Message     string          `db:"message" json:"message"`
	Type        string          `db:"type" json:"type"`
	ReferenceID sql.NullInt64   `db:"reference_id" json:"reference_id"`
}

// IsReadBool retorna el valor booleano del campo IsRead
func (nd *NotificationDelivery) IsReadBool() bool {
	return nd.IsRead.Valid && nd.IsRead.Int64 == 1
}

// MarshalJSON implementa la serialización JSON personalizada
func (nd *NotificationDelivery) MarshalJSON() ([]byte, error) {
	var referenceID *int64
	if nd.ReferenceID.Valid {
		referenceID = &nd.ReferenceID.Int64
	}
	return json.Marshal(map[string]interface{}{
		"node_id":       nd.NodeID,
		"created_at":    nd.CreatedAt,
		"is_read":       nd.IsReadBool(),
		"to_account_id": nd.ToAccountID,
		"noti_id":       nd.NotiID,
		"title":         nd.Title,
		"message":       nd.Message,
		"type":          nd.Type,
		"reference_id":  referenceID,
	})
}

func NewRepository(db *sql.DB) Repository {
	return &mysqlRepository{db: db}
}

func (r *mysqlRepository) Save(n Notification) (int64, error) {
	var notiID int64

	query := `INSERT INTO notifications(noti_id, owner_account_id, title, message, type, link, must_send_as_notification_push, must_send_as_notification, must_be_processed, error_message, reference_id)
	VALUES(?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`

	result, err := r.db.Exec(query, n.OwnerAccountID, n.Title, n.Message, n.Type, n.Link, n.MustSendAsNotificationPush, n.MustSendAsNotification, n.MustBeProcessed, n.ErrorMesssage, n.ReferenceID)

	if err != nil {
		log.Println("Error saving notification...")
		return notiID, err
	}

	notiID, err = result.LastInsertId()

	if err != nil {
		log.Println("Error getting noti_id notification...")
		return notiID, err
	}

	log.Println("Notification has been saved succesfully...")
	return notiID, nil
}

func (r *mysqlRepository) SaveDeviceToken(accountID int64, token string) error {
	query := `
        INSERT INTO device_tokens (account_id, device_token)
        VALUES (?, ?)
        ON DUPLICATE KEY UPDATE updated_at = CURRENT_TIMESTAMP;
    `
	if _, err := r.db.Exec(query, accountID, token); err != nil {
		return fmt.Errorf("SaveDeviceToken: %w", err)
	}
	return nil
}

func (r *mysqlRepository) DeleteDeviceToken(accountID int64, deviceToken string) error {
	_, err := r.db.Exec(`
	  DELETE FROM device_tokens 
	  WHERE account_id = ? AND device_token = ?`,
		accountID, deviceToken,
	)
	return err
}

// GetNotifications obtiene las notificaciones del usuario con paginación
func (r *mysqlRepository) GetNotifications(accountID int64, limit, offset int) ([]NotificationDelivery, error) {
	query := `
		SELECT
			nd.node_id,
			nd.created_at,
			nd.is_read,
			nd.to_account_id,
			nd.noti_id,
			nd.title,
			nd.message,
			n.type,
			n.reference_id
		FROM notification_deliveries nd
		LEFT JOIN notifications n ON nd.noti_id = n.noti_id
		WHERE nd.to_account_id = ?
		ORDER BY nd.created_at DESC
		LIMIT ? OFFSET ?
	`

	rows, err := r.db.Query(query, accountID, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("GetNotifications query error: %w", err)
	}
	defer rows.Close()

	var notifications []NotificationDelivery
	for rows.Next() {
		var nd NotificationDelivery
		var createdAt sql.NullTime
		err := rows.Scan(
			&nd.NodeID,
			&createdAt,
			&nd.IsRead,
			&nd.ToAccountID,
			&nd.NotiID,
			&nd.Title,
			&nd.Message,
			&nd.Type,
			&nd.ReferenceID,
		)
		if err != nil {
			log.Printf("Error scanning notification: %v", err)
			continue
		}

		// Convertir sql.NullTime a common.NullTime
		if createdAt.Valid {
			nd.CreatedAt = common.NullTime{NullTime: sql.NullTime{Time: createdAt.Time, Valid: true}}
		} else {
			nd.CreatedAt = common.NullTime{NullTime: sql.NullTime{Time: time.Time{}, Valid: false}}
		}

		notifications = append(notifications, nd)
	}

	return notifications, nil
}

// GetUnreadCount obtiene el conteo de notificaciones no leídas
func (r *mysqlRepository) GetUnreadCount(accountID int64) (int64, error) {
	query := `
		SELECT COUNT(*)
		FROM notification_deliveries
		WHERE to_account_id = ? AND (is_read = 0 OR is_read IS NULL)
	`

	var count int64
	err := r.db.QueryRow(query, accountID).Scan(&count)
	if err != nil {
		return 0, fmt.Errorf("GetUnreadCount query error: %w", err)
	}

	return count, nil
}

// MarkAsRead marca una notificación como leída
func (r *mysqlRepository) MarkAsRead(accountID, notificationID int64) error {
	// First, let's check if the notification exists and get its current status
	var isRead sql.NullInt64
	checkQuery := `SELECT is_read FROM notification_deliveries WHERE to_account_id = ? AND node_id = ?`
	err := r.db.QueryRow(checkQuery, accountID, notificationID).Scan(&isRead)
	if err != nil {
		if err == sql.ErrNoRows {
			log.Printf("MarkAsRead: Notification not found - node_id=%d, account_id=%d", notificationID, accountID)
			return fmt.Errorf("notification not found")
		}
		log.Printf("MarkAsRead: Error checking notification status: %v", err)
		return fmt.Errorf("error checking notification status: %w", err)
	}

	// Check if already read
	if isRead.Valid && isRead.Int64 == 1 {
		log.Printf("MarkAsRead: Notification already read - node_id=%d, account_id=%d", notificationID, accountID)
		return fmt.Errorf("notification already read")
	}

	// Update the notification as read
	updateQuery := `UPDATE notification_deliveries SET is_read = 1 WHERE to_account_id = ? AND node_id = ?`
	result, err := r.db.Exec(updateQuery, accountID, notificationID)
	if err != nil {
		log.Printf("MarkAsRead: Update query error: %v", err)
		return fmt.Errorf("update query error: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		log.Printf("MarkAsRead: Rows affected error: %v", err)
		return fmt.Errorf("rows affected error: %w", err)
	}

	if rowsAffected == 0 {
		log.Printf("MarkAsRead: No rows updated - unexpected - node_id=%d, account_id=%d", notificationID, accountID)
		return fmt.Errorf("unexpected error: no rows updated")
	}

	log.Printf("MarkAsRead: Successfully marked notification as read - node_id=%d, account_id=%d", notificationID, accountID)
	return nil
}

// MarkAllAsRead marca todas las notificaciones como leídas
func (r *mysqlRepository) MarkAllAsRead(accountID int64) error {
	query := `
		UPDATE notification_deliveries
		SET is_read = 1
		WHERE to_account_id = ? AND (is_read = 0 OR is_read IS NULL)
	`

	_, err := r.db.Exec(query, accountID)
	if err != nil {
		return fmt.Errorf("MarkAllAsRead query error: %w", err)
	}

	return nil
}

// DeleteNotification elimina una notificación
func (r *mysqlRepository) DeleteNotification(accountID, notificationID int64) error {
	query := `
		DELETE FROM notification_deliveries
		WHERE to_account_id = ? AND node_id = ?
	`

	result, err := r.db.Exec(query, accountID, notificationID)
	if err != nil {
		return fmt.Errorf("DeleteNotification query error: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("DeleteNotification rows affected error: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("notification not found")
	}

	return nil
}
