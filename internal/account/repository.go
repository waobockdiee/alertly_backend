package account

import (
	"database/sql"
	"log"
)

type Repository interface {
	GetMyInfo(accountID int64) (MyInfo, error)
	GetHistory(accountID int64) ([]History, error)
	ClearHistory(accountID int64) error
	DeleteAccount(accountID int64) error
	GetCounterHistories(accountID int64) (Counter, error)
	SaveLastRequest(AccountID int64, ip string) error
	SetHasFinishedTutorial(accountID int64) error
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

	if err != nil {
		return nil, err
	}

	return histories, nil
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
