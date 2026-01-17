package account

import (
	"database/sql"
	"log"
	"time"

	"alertly/internal/dbtypes"
)

type Repository interface {
	GetMyInfo(accountID int64) (MyInfo, error)
	GetHistory(accountID int64) ([]History, error)
	GetViewedIncidentIds(accountID int64) ([]int64, error)
	ClearHistory(accountID int64) error
	DeleteAccount(accountID int64) error
	GetAccountPassword(accountID int64) (string, error)
	GetCounterHistories(accountID int64) (Counter, error)
	SaveLastRequest(AccountID int64, ip string) error
	SetHasFinishedTutorial(accountID int64) error
	UpdatePremiumStatus(accountID int64, isPremium bool, subscriptionType string, expirationDate *time.Time, platform string) error
}

type pgRepository struct {
	db *sql.DB
}

func NewRepository(db *sql.DB) Repository {
	return &pgRepository{db: db}
}

func (r *pgRepository) GetMyInfo(accountID int64) (MyInfo, error) {
	var myInfo MyInfo

	// Usar NullBool para campos booleanos que pueden ser SMALLINT/CHAR/BOOLEAN
	var isPremium, hasFinishedTutorial dbtypes.NullBool

	query := `SELECT account_id, email, is_premium, status, has_finished_tutorial FROM account WHERE account_id = $1`
	err := r.db.QueryRow(query, accountID).Scan(
		&myInfo.AccountID,
		&myInfo.Email,
		&isPremium,
		&myInfo.Status,
		&hasFinishedTutorial,
	)

	if err != nil {
		log.Printf("Error fetching MyInfo for account ID %d: %v", accountID, err)
		return myInfo, err
	}

	// Convertir NullBool a bool
	myInfo.IsPremium = isPremium.Valid && isPremium.Bool

	// HasFinishedTutorial es string en el modelo (legacy)
	if hasFinishedTutorial.Valid && hasFinishedTutorial.Bool {
		myInfo.HasFinishedTutorial = "1"
	} else {
		myInfo.HasFinishedTutorial = "0"
	}

	return myInfo, nil
}

func (r *pgRepository) GetHistory(accountID int64) ([]History, error) {
	query := `SELECT
	t1.his_id, t1.account_id, t1.incl_id, t1.created_at, t2.address, t2.description
	FROM account_history t1 INNER JOIN incident_clusters t2 ON t1.incl_id = t2.incl_id
	WHERE t1.account_id = $1 ORDER BY t1.his_id DESC LIMIT 1000`
	rows, err := r.db.Query(query, accountID)

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var histories []History

	for rows.Next() {
		var h History

		if err := rows.Scan(
			&h.HisID,
			&h.AccountID,
			&h.InclID,
			&h.CreatedAt,
			&h.Address,
			&h.Description,
		); err != nil {
			return nil, err
		}
		histories = append(histories, h)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return histories, nil
}

func (r *pgRepository) GetViewedIncidentIds(accountID int64) ([]int64, error) {
	query := `SELECT incl_id FROM account_history WHERE account_id = $1 ORDER BY created_at DESC`
	rows, err := r.db.Query(query, accountID)

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var incidentIds []int64

	for rows.Next() {
		var inclId int64
		if err := rows.Scan(&inclId); err != nil {
			return nil, err
		}
		incidentIds = append(incidentIds, inclId)
	}

	return incidentIds, nil
}

func (r *pgRepository) GetCounterHistories(accountID int64) (Counter, error) {
	var counter Counter
	query := "SELECT COUNT(*) AS counter FROM account_history WHERE account_id = $1"
	err := r.db.QueryRow(query, accountID).Scan(&counter.Counter)

	if err != nil {
		return counter, err
	}
	return counter, err
}

func (r *pgRepository) ClearHistory(accountID int64) error {
	query := `DELETE FROM account_history WHERE account_id = $1`
	_, err := r.db.Exec(query, accountID)

	if err != nil {
		return err
	}
	return nil
}

func (r *pgRepository) GetAccountPassword(accountID int64) (string, error) {
	var password string
	query := "SELECT password FROM account WHERE account_id = $1"
	err := r.db.QueryRow(query, accountID).Scan(&password)

	if err != nil {
		log.Printf("Error fetching password for account %d: %v", accountID, err)
		return "", err
	}

	return password, nil
}

func (r *pgRepository) DeleteAccount(accountID int64) error {
	// Start a transaction to ensure all deletions succeed or fail together
	tx, err := r.db.Begin()
	if err != nil {
		log.Printf("Error starting transaction for account deletion: %v", err)
		return err
	}
	defer tx.Rollback()

	// Get user's thumbnail URL before deletion (to delete from S3 later)
	var thumbnailURL sql.NullString
	queryThumbnail := "SELECT thumbnail_url FROM account WHERE account_id = $1"
	err = tx.QueryRow(queryThumbnail, accountID).Scan(&thumbnailURL)
	if err != nil && err != sql.ErrNoRows {
		log.Printf("Error fetching thumbnail for account %d: %v", accountID, err)
		return err
	}

	// 1. Delete account history
	_, err = tx.Exec("DELETE FROM account_history WHERE account_id = $1", accountID)
	if err != nil {
		log.Printf("Error deleting account_history for account %d: %v", accountID, err)
		return err
	}

	// 2. Delete saved clusters
	_, err = tx.Exec("DELETE FROM saved_clusters_account WHERE account_id = $1", accountID)
	if err != nil {
		log.Printf("Error deleting saved_clusters_account for account %d: %v", accountID, err)
		return err
	}

	// 3. Delete notification deliveries
	_, err = tx.Exec("DELETE FROM notification_deliveries WHERE account_id = $1", accountID)
	if err != nil {
		log.Printf("Error deleting notification_deliveries for account %d: %v", accountID, err)
		return err
	}

	// 4. Delete notifications
	_, err = tx.Exec("DELETE FROM notifications WHERE account_id = $1", accountID)
	if err != nil {
		log.Printf("Error deleting notifications for account %d: %v", accountID, err)
		return err
	}

	// 5. Delete device tokens
	_, err = tx.Exec("DELETE FROM account_device_tokens WHERE account_id = $1", accountID)
	if err != nil {
		log.Printf("Error deleting account_device_tokens for account %d: %v", accountID, err)
		return err
	}

	// 6. Delete session history
	_, err = tx.Exec("DELETE FROM account_session_history WHERE account_id = $1", accountID)
	if err != nil {
		log.Printf("Error deleting account_session_history for account %d: %v", accountID, err)
		return err
	}

	// 7. Delete premium payment history
	_, err = tx.Exec("DELETE FROM account_premium_payment_history WHERE account_id = $1", accountID)
	if err != nil {
		log.Printf("Error deleting account_premium_payment_history for account %d: %v", accountID, err)
		return err
	}

	// 8. Delete favorite locations (my places)
	_, err = tx.Exec("DELETE FROM account_favorite_locations WHERE account_id = $1", accountID)
	if err != nil {
		log.Printf("Error deleting account_favorite_locations for account %d: %v", accountID, err)
		return err
	}

	// 9. Delete incident reports created by this user
	// NOTE: We're deleting the user's incident_reports, but NOT the incident_clusters
	// because clusters may contain reports from multiple users
	_, err = tx.Exec("DELETE FROM incident_reports WHERE account_id = $1", accountID)
	if err != nil {
		log.Printf("Error deleting incident_reports for account %d: %v", accountID, err)
		return err
	}

	// 10. Delete achievement progress
	_, err = tx.Exec("DELETE FROM achievement_progress WHERE account_id = $1", accountID)
	if err != nil {
		log.Printf("Error deleting achievement_progress for account %d: %v", accountID, err)
		return err
	}

	// 11. Finally, delete the account itself
	_, err = tx.Exec("DELETE FROM account WHERE account_id = $1", accountID)
	if err != nil {
		log.Printf("Error deleting account %d: %v", accountID, err)
		return err
	}

	// Commit the transaction
	err = tx.Commit()
	if err != nil {
		log.Printf("Error committing account deletion transaction for account %d: %v", accountID, err)
		return err
	}

	log.Printf("✅ Successfully deleted account %d and all related data", accountID)

	// TODO: Delete user's images from S3 (thumbnail and incident photos)
	// This should be done asynchronously after the transaction commits
	// if thumbnailURL.Valid && thumbnailURL.String != "" {
	//     go deleteFromS3(thumbnailURL.String)
	// }

	return nil
}

func (r *pgRepository) SaveLastRequest(AccountID int64, ip string) error {
	query := `INSERT INTO account_session_history (account_id, ip) VALUES($1, $2)`
	_, err := r.db.Exec(query, AccountID, ip)

	return err
}
func (r *pgRepository) SetHasFinishedTutorial(accountID int64) error {
	// Usar dbtypes.BoolToInt para insertar en columnas SMALLINT
	query := "UPDATE account SET has_finished_tutorial = $1 WHERE account_id = $2"
	_, err := r.db.Exec(query, dbtypes.BoolToInt(true), accountID)
	if err != nil {
		log.Printf("Error actualizando has_finished_tutorial (ID: %d): %v", accountID, err)
	}
	return err
}

// UpdatePremiumStatus updates the user's premium status and logs the payment history
func (r *pgRepository) UpdatePremiumStatus(accountID int64, isPremium bool, subscriptionType string, expirationDate *time.Time, platform string) error {
	// Start a transaction to ensure both operations succeed or fail together
	tx, err := r.db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	// 1. Update the account's is_premium status and expiration date
	var updateAccountQuery string
	var args []interface{}

	if isPremium && expirationDate != nil {
		updateAccountQuery = "UPDATE account SET is_premium = $1, premium_expired_date = $2 WHERE account_id = $3"
		args = []interface{}{dbtypes.BoolToInt(isPremium), expirationDate, accountID}
	} else if !isPremium {
		// When cancelling or expiring, set is_premium to false and clear expiration date
		updateAccountQuery = "UPDATE account SET is_premium = $1, premium_expired_date = NULL WHERE account_id = $2"
		args = []interface{}{dbtypes.BoolToInt(isPremium), accountID}
	} else {
		// Fallback for safety, though should not be reached in normal flow
		updateAccountQuery = "UPDATE account SET is_premium = $1 WHERE account_id = $2"
		args = []interface{}{dbtypes.BoolToInt(isPremium), accountID}
	}

	_, err = tx.Exec(updateAccountQuery, args...)
	if err != nil {
		log.Printf("Error updating account premium status for account %d: %v", accountID, err)
		return err
	}

	// 2. Map subscription type to your enum values for history
	var typePlan string
	if isPremium {
		// We use the subscriptionType which should now be the product_id from Apple
		typePlan = subscriptionType
	} else {
		typePlan = "free"
	}

	// 3. Create description based on the action
	var description string
	if isPremium {
		description = "Premium subscription activated or updated via " + platform + ". Product ID: " + subscriptionType
		if expirationDate != nil {
			description += ". Expires on: " + expirationDate.Format("2006-01-02 15:04:05")
		}
	} else {
		description = "Premium subscription cancelled or expired"
	}

	// 4. Insert payment history record
	insertHistoryQuery := `
		INSERT INTO account_premium_payment_history
		(account_id, type_plan, description, created)
		VALUES ($1, $2, $3, NOW())
	`
	_, err = tx.Exec(insertHistoryQuery, accountID, typePlan, description)
	if err != nil {
		log.Printf("Error inserting premium payment history for account %d: %v", accountID, err)
		return err
	}

	// Commit the transaction
	err = tx.Commit()
	if err != nil {
		log.Printf("Error committing premium status update transaction for account %d: %v", accountID, err)
		return err
	}

	log.Printf("✅ Successfully updated premium status for account %d: isPremium=%v, typePlan=%s",
		accountID, isPremium, typePlan)

	return nil
}
