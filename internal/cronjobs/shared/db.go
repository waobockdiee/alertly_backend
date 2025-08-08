package shared

import (
	"context"
	"database/sql"
	"fmt"
	"strings"
	"time"
)

// Delivery representa un registro en la tabla notification_deliveries.
// Esta struct se puede mover a un paquete de modelo si es necesario.
type Delivery struct {
	NotificationID int64
	AccountID      int64
	Title          string
	Message        string
}

// MarkItemsAsProcessed marca un conjunto de registros como procesados en una tabla específica.
// Utiliza un nombre de tabla y un nombre de columna de ID dinámicos para mayor flexibilidad.
func MarkItemsAsProcessed(db *sql.DB, tableName string, idColumn string, ids []int64) error {
	if len(ids) == 0 {
		return nil
	}
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// Crear placeholders para los IDs
	placeholders := strings.Repeat("?,", len(ids)-1) + "?"

	// Construir la consulta dinámicamente
	query := fmt.Sprintf(
		"UPDATE %s SET must_be_processed = 0, sent_at = NOW() WHERE %s IN (%s)",
		tableName, idColumn, placeholders,
	)

	// Convertir los IDs a []interface{} para ExecContext
	args := make([]interface{}, len(ids))
	for i, id := range ids {
		args[i] = id
	}

	_, err := db.ExecContext(ctx, query, args...)
	return err
}

// InsertDeliveries inserta en bloque los registros en notification_deliveries.
func InsertDeliveries(db *sql.DB, deliveries []Delivery) error {
	if len(deliveries) == 0 {
		return nil
	}
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	base := `INSERT INTO notification_deliveries (noti_id, to_account_id, title, message, created_at) VALUES `
	var placeholders []string
	var args []interface{}
	now := time.Now()

	for _, d := range deliveries {
		placeholders = append(placeholders, "(?,?,?,?,?)")
		args = append(args, d.NotificationID, d.AccountID, d.Title, d.Message, now)
	}

	query := base + strings.Join(placeholders, ",")
	_, err := db.ExecContext(ctx, query, args...)
	return err
}

// GetDeviceTokensForAccount returns all device tokens for a given account
func GetDeviceTokensForAccount(db *sql.DB, accountID int64) ([]string, error) {
	rows, err := db.Query(
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
