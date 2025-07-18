package cjnewcluster

import (
	"database/sql"
	"fmt"
	"time"
)

// Repository encapsula acceso a BD
type Repository struct {
	db *sql.DB
}

func NewRepository(db *sql.DB) *Repository {
	return &Repository{db: db}
}

// FetchPending obtiene notificaciones no procesadas en batch
func (r *Repository) FetchPending(limit int64) ([]Notification, error) {
	rows, err := r.db.Query(
		`SELECT id, cluster_id, created_at FROM notifications
         WHERE type = 'new_cluster' AND processed = FALSE
         ORDER BY created_at
         LIMIT ?`, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var list []Notification
	for rows.Next() {
		var n Notification
		if err := rows.Scan(&n.ID, &n.ClusterID, &n.CreatedAt); err != nil {
			return nil, err
		}
		list = append(list, n)
	}
	return list, nil
}

// MarkProcessed actualiza processed=true
func (r *Repository) MarkProcessed(ids []int64) error {
	query := "UPDATE notifications SET processed = TRUE WHERE id IN ("
	params := make([]interface{}, len(ids))
	for i, id := range ids {
		query += fmt.Sprintf("?%s", ",")
		params[i] = id
	}
	query = query[:len(query)-1] + ")"
	_, err := r.db.Exec(query, params...)
	return err
}

// FetchRecipients obtiene cuentas a notificar para múltiples clusters
func (r *Repository) FetchRecipients(clusterIDs []int64) (map[int64][]int64, error) {
	// Devolver map[clusterID] -> []accountID
	// Ejemplo: JOIN accounts_clusters
	in := placeholders(len(clusterIDs))
	args := interfaceSlice(clusterIDs)
	q := fmt.Sprintf(
		`SELECT cluster_id, account_id FROM accounts_clusters
         WHERE cluster_id IN (%s)`, in)

	rows, err := r.db.Query(q, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	m := make(map[int64][]int64)
	for rows.Next() {
		var cid, aid int64
		rows.Scan(&cid, &aid)
		m[cid] = append(m[cid], aid)
	}
	return m, nil
}

// InsertDeliveries inserta en batch registros de envío
func (r *Repository) InsertDeliveries(delivs []Delivery) error {
	tx, err := r.db.Begin()
	if err != nil {
		return err
	}
	stmt, _ := tx.Prepare(`INSERT INTO notifications_deliveries
        (notification_id, account_id, sent_at) VALUES (?,?,?)`)
	defer stmt.Close()
	now := time.Now()
	for _, d := range delivs {
		if _, err := stmt.Exec(d.NotificationID, d.AccountID, now); err != nil {
			tx.Rollback()
			return err
		}
	}
	return tx.Commit()
}

// helpers
func placeholders(n int) string {
	s := ""
	for i := 0; i < n; i++ {
		s += "?,"
	}
	return s[:len(s)-1]
}
func interfaceSlice(ids []int64) []interface{} {
	out := make([]interface{}, len(ids))
	for i, v := range ids {
		out[i] = v
	}
	return out
}

// GetDeviceTokensForAccount returns all device tokens for a given account
func (r *Repository) GetDeviceTokensForAccount(accountID int64) ([]string, error) {
	rows, err := r.db.Query(
		`SELECT device_token FROM device_tokens WHERE account_id = ?`,
		accountID,
	)
	if err != nil {
		return nil, fmt.Errorf("GetDeviceTokensForAccount: %w", err)
	}
	defer rows.Close()

	var tokens []string
	for rows.Next() {
		var t string
		if err := rows.Scan(&t); err != nil {
			return nil, fmt.Errorf("scanning device_token: %w", err)
		}
		tokens = append(tokens, t)
	}
	return tokens, nil
}
