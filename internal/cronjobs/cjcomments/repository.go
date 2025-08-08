package cjcomments

import (
	"alertly/internal/cronjobs/shared"
	"database/sql"
	"fmt"
	"strings"
)

// Repository encapsula el acceso a la base de datos para el cronjob de comentarios.
type Repository struct {
	db *sql.DB
}

// NewRepository crea una nueva instancia de Repository.
func NewRepository(db *sql.DB) *Repository {
	return &Repository{db: db}
}

// FetchPendingCommentNotifications obtiene notificaciones de comentarios pendientes.
func (r *Repository) FetchPendingCommentNotifications(limit int64) ([]CommentNotification, error) {
	query := `
        SELECT noti_id, reference_id, created_at
        FROM notifications
        WHERE type = 'new_comment' AND must_be_processed = 1
        ORDER BY created_at
        LIMIT ?
    `
	rows, err := r.db.Query(query, limit)
	if err != nil {
		return nil, fmt.Errorf("FetchPendingCommentNotifications: %w", err)
	}
	defer rows.Close()

	var notifications []CommentNotification
	for rows.Next() {
		var cn CommentNotification
		if err := rows.Scan(&cn.NotificationID, &cn.CommentID, &cn.CreatedAt); err != nil {
			return nil, fmt.Errorf("scanning comment notification: %w", err)
		}
		notifications = append(notifications, cn)
	}
	return notifications, nil
}

// GetCommentDetails obtiene los detalles de un comentario específico.
func (r *Repository) GetCommentDetails(commentID int64) (*CommentDetails, error) {
	query := `
        SELECT
            inc.inco_id, inc.incl_id, inc.comment, inc.account_id, ic.subcategory_name
        FROM
            incident_comments inc
        JOIN
            incident_clusters ic ON inc.incl_id = ic.incl_id
        WHERE
            inc.inco_id = ?
    `
	var cd CommentDetails
	err := r.db.QueryRow(query, commentID).Scan(&cd.CommentID, &cd.ClusterID, &cd.CommentText, &cd.CommenterID, &cd.SubcategoryName)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("comment with ID %d not found", commentID)
		}
		return nil, fmt.Errorf("GetCommentDetails: %w", err)
	}
	return &cd, nil
}

// GetIncidentCreators obtiene los account_id de los creadores de incidentes para un cluster.
func (r *Repository) GetIncidentCreators(clusterID int64) ([]int64, error) {
	query := `
        SELECT DISTINCT account_id
        FROM incident_reports
        WHERE incl_id = ?
    `
	rows, err := r.db.Query(query, clusterID)
	if err != nil {
		return nil, fmt.Errorf("GetIncidentCreators: %w", err)
	}
	defer rows.Close()

	var creators []int64
	for rows.Next() {
		var accountID int64
		if err := rows.Scan(&accountID); err != nil {
			return nil, fmt.Errorf("scanning incident creator: %w", err)
		}
		creators = append(creators, accountID)
	}
	return creators, nil
}

// GetSavedClusterUsers obtiene los account_id de los usuarios que guardaron un cluster.
func (r *Repository) GetSavedClusterUsers(clusterID int64) ([]int64, error) {
	query := `
        SELECT account_id
        FROM account_cluster_saved
        WHERE incl_id = ?
    `
	rows, err := r.db.Query(query, clusterID)
	if err != nil {
		return nil, fmt.Errorf("GetSavedClusterUsers: %w", err)
	}
	defer rows.Close()

	var users []int64
	for rows.Next() {
		var accountID int64
		if err := rows.Scan(&accountID); err != nil {
			return nil, fmt.Errorf("scanning saved cluster user: %w", err)
		}
		users = append(users, accountID)
	}
	return users, nil
}

// GetDeviceTokensForAccounts obtiene los tokens de dispositivo para una lista de account_ids.
func (r *Repository) GetDeviceTokensForAccounts(accountIDs []int64) ([]Recipient, error) {
	if len(accountIDs) == 0 {
		return nil, nil
	}
	// Construir placeholders para la cláusula IN
	placeholders := strings.Repeat("?,", len(accountIDs)-1) + "?"
	query := fmt.Sprintf(`
        SELECT dt.account_id, dt.device_token
        FROM device_tokens dt
        JOIN account a ON dt.account_id = a.account_id
        WHERE dt.account_id IN (%s) AND a.status = 'active' AND a.receive_notifications = 1
    `, placeholders)

	// Convertir []int64 a []interface{} para Query
	args := make([]interface{}, len(accountIDs))
	for i, id := range accountIDs {
		args[i] = id
	}

	rows, err := r.db.Query(query, args...)
	if err != nil {
		return nil, fmt.Errorf("GetDeviceTokensForAccounts: %w", err)
	}
	defer rows.Close()

	var recipients []Recipient
	for rows.Next() {
		var rec Recipient
		if err := rows.Scan(&rec.AccountID, &rec.DeviceToken); err != nil {
			return nil, fmt.Errorf("scanning recipient: %w", err)
		}
		recipients = append(recipients, rec)
	}
	return recipients, nil
}

// InsertNotificationDeliveries inserta en bloque los registros de envío de notificaciones.
func (r *Repository) InsertNotificationDeliveries(deliveries []shared.Delivery) error {
	return shared.InsertDeliveries(r.db, deliveries)
}

// MarkNotificationProcessed marca una notificación como procesada.
func (r *Repository) MarkNotificationProcessed(notificationID int64) error {
	return shared.MarkItemsAsProcessed(r.db, "notifications", "noti_id", []int64{notificationID})
}
