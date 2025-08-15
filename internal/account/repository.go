package account

import (
	"database/sql"
	"log"
	"time"
)

type Repository interface {
	GetMyInfo(accountID int64) (MyInfo, error)
	GetHistory(accountID int64) ([]History, error)
	GetViewedIncidentIds(accountID int64) ([]int64, error)
	ClearHistory(accountID int64) error
	DeleteAccount(accountID int64) error
	GetCounterHistories(accountID int64) (Counter, error)
	SaveLastRequest(AccountID int64, ip string) error
	SetHasFinishedTutorial(accountID int64) error
	UpdatePremiumStatus(accountID int64, isPremium bool, subscriptionType string, purchaseDate *time.Time, platform string) error
}

type mysqlRepository struct {
	db *sql.DB
}

func NewRepository(db *sql.DB) Repository {
	return &mysqlRepository{db: db}
}

func (r *mysqlRepository) GetMyInfo(accountID int64) (MyInfo, error) {
	var myInfo MyInfo

	query := `SELECT account_id, email, is_premium, status, has_finished_tutorial FROM account WHERE account_id = ?`
	err := r.db.QueryRow(query, accountID).Scan(&myInfo.AccountID, &myInfo.Email, &myInfo.IsPremium, &myInfo.Status, &myInfo.HasFinishedTutorial)

	if err != nil {
		log.Printf("Error fetching MyInfo for account ID %d: %v", accountID, err)
		return myInfo, err
	}

	return myInfo, nil
}

func (r *mysqlRepository) GetHistory(accountID int64) ([]History, error) {
	query := `SELECT
	t1.his_id, t1.account_id, t1.incl_id, t1.created_at, t2.address, t2.description
	FROM account_history t1 INNER JOIN incident_clusters t2 ON t1.incl_id = t2.incl_id
	WHERE t1.account_id = ? ORDER BY t1.his_id DESC LIMIT 1000`
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

func (r *mysqlRepository) GetViewedIncidentIds(accountID int64) ([]int64, error) {
	query := `SELECT incl_id FROM account_history WHERE account_id = ? ORDER BY created_at DESC`
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

func (r *mysqlRepository) GetCounterHistories(accountID int64) (Counter, error) {
	var counter Counter
	query := "SELECT COUNT(*) AS counter FROM account_history WHERE account_id = ?"
	err := r.db.QueryRow(query, accountID).Scan(&counter.Counter)

	if err != nil {
		return counter, err
	}
	return counter, err
}

func (r *mysqlRepository) ClearHistory(accountID int64) error {
	query := `DELETE FROM account_history WHERE account_id = ?`
	_, err := r.db.Exec(query, accountID)

	if err != nil {
		return err
	}
	return nil
}

func (r *mysqlRepository) DeleteAccount(accountID int64) error {
	return nil
}

func (r *mysqlRepository) SaveLastRequest(AccountID int64, ip string) error {
	query := `INSERT INTO account_session_history (account_id, ip) VALUES(?, ?)`
	_, err := r.db.Exec(query, AccountID, ip)

	return err
}
func (r *mysqlRepository) SetHasFinishedTutorial(accountID int64) error {
	query := "UPDATE account SET has_finished_tutorial = 1 WHERE account_id = ?"
	_, err := r.db.Exec(query, accountID)
	if err != nil {
		log.Printf("Error actualizando notificación (ID: %d) como procesada: %v", accountID, err)
	}
	return err
}

// UpdatePremiumStatus updates the user's premium status and logs the payment history
func (r *mysqlRepository) UpdatePremiumStatus(accountID int64, isPremium bool, subscriptionType string, purchaseDate *time.Time, platform string) error {
	// Start a transaction to ensure both operations succeed or fail together
	tx, err := r.db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	// 1. Calculate expiration date based on subscription type
	var expirationDate *time.Time
	if isPremium && purchaseDate != nil {
		expDate := *purchaseDate
		switch subscriptionType {
		case "monthly":
			expDate = expDate.AddDate(0, 1, 0) // Add 1 month
		case "yearly":
			expDate = expDate.AddDate(1, 0, 0) // Add 1 year
		default:
			expDate = expDate.AddDate(0, 1, 0) // Default to 1 month
		}
		expirationDate = &expDate
	}

	// 2. Update the account's is_premium status and expiration date
	var updateAccountQuery string
	var args []interface{}

	if isPremium && expirationDate != nil {
		updateAccountQuery = "UPDATE account SET is_premium = ?, premium_expired_date = ? WHERE account_id = ?"
		args = []interface{}{isPremium, expirationDate, accountID}
	} else if !isPremium {
		// When cancelling, set is_premium to false and clear expiration date
		updateAccountQuery = "UPDATE account SET is_premium = ?, premium_expired_date = NULL WHERE account_id = ?"
		args = []interface{}{isPremium, accountID}
	} else {
		// Fallback: just update is_premium
		updateAccountQuery = "UPDATE account SET is_premium = ? WHERE account_id = ?"
		args = []interface{}{isPremium, accountID}
	}

	_, err = tx.Exec(updateAccountQuery, args...)
	if err != nil {
		log.Printf("Error updating account premium status for account %d: %v", accountID, err)
		return err
	}

	// 2. Map subscription type to your enum values
	var typePlan string
	switch subscriptionType {
	case "monthly":
		typePlan = "1 month"
	case "yearly":
		typePlan = "12 months"
	case "restored":
		typePlan = "1 month" // Default for restored purchases
	default:
		if isPremium {
			typePlan = "1 month" // Default to monthly if unknown
		} else {
			typePlan = "free"
		}
	}

	// 3. Create description based on the action
	var description string
	if isPremium {
		description = "Premium subscription activated via " + platform + " platform. Subscription type: " + subscriptionType
		if purchaseDate != nil {
			description += ". Purchase date: " + purchaseDate.Format("2006-01-02 15:04:05")
		}
	} else {
		description = "Premium subscription cancelled or expired"
		typePlan = "free"
	}

	// 4. Insert payment history record
	insertHistoryQuery := `
		INSERT INTO account_premium_payment_history 
		(account_id, type_plan, description, created) 
		VALUES (?, ?, ?, NOW())
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
