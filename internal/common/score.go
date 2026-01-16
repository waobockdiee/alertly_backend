package common

import (
	"alertly/internal/dbtypes"
	"fmt"
)

func SaveScore(dbExec DBExecutor, accountID int64, score uint8) error {
	// 1. Actualizar score en la cuenta
	updateQuery := `UPDATE account SET score = score + $1 WHERE account_id = $2`
	_, err := dbExec.Exec(updateQuery, score, accountID)
	if err != nil {
		return fmt.Errorf("failed to update score: %w", err)
	}

	// 2. Crear notificación citizen score (solo in-app)
	err = saveScoreNotification(dbExec, accountID, score)
	if err != nil {
		// Log pero no fallas - el score ya se guardó exitosamente
		fmt.Printf("Warning: failed to create score notification: %v\n", err)
	}

	return nil
}

// saveScoreNotification crea una notificación in-app para citizen score
func saveScoreNotification(dbExec DBExecutor, accountID int64, score uint8) error {
	// Crear título personalizado con puntos ganados
	title := fmt.Sprintf("Congratulations! You've Earned %d Citizen Points.", score)
	message := "Keep contributing to your community!"

	// Insertar notificación con RETURNING para PostgreSQL (LastInsertId no funciona en pq)
	// Convertir bool a int (0/1) para SMALLINT columns
	notiQuery := `INSERT INTO notifications(owner_account_id, title, message, type, link, must_send_as_notification_push, must_send_as_notification, must_be_processed, reference_id, created_at)
				  VALUES($1, $2, $3, $4, $5, $6, $7, $8, $9, NOW()) RETURNING noti_id`

	var notiID int64
	err := dbExec.(interface {
		QueryRow(query string, args ...interface{}) interface{ Scan(dest ...interface{}) error }
	}).QueryRow(notiQuery, accountID, title, message, "earn_citizen_score", "ProfileScreen", dbtypes.BoolToInt(false), dbtypes.BoolToInt(true), dbtypes.BoolToInt(false), score).Scan(&notiID)
	if err != nil {
		return fmt.Errorf("failed to create notification: %w", err)
	}

	// Crear delivery inmediatamente para el usuario (para que aparezca en /notifications)
	deliveryQuery := `INSERT INTO notification_deliveries (noti_id, to_account_id, title, message, created_at)
					  VALUES ($1, $2, $3, $4, NOW())`
	_, err = dbExec.Exec(deliveryQuery, notiID, accountID, title, message)
	if err != nil {
		return fmt.Errorf("failed to create notification delivery: %w", err)
	}

	return nil
}
