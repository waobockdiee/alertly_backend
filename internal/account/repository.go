package account

import (
	"database/sql"
)

type Repository interface {
	GetHistory(accountID int64) ([]History, error)
	ClearHistory(accountID int64) error
	DeleteAccount(accountID int64) error
}

type mysqlRepository struct {
	db *sql.DB
}

func NewRepository(db *sql.DB) Repository {
	return &mysqlRepository{db: db}
}

func (r *mysqlRepository) GetHistory(accountID int64) ([]History, error) {
	query := `SELECT
	his_id, account_id, incl_id, created_at
	FROM account_history WHERE account_id = ?`
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

func (r *mysqlRepository) ClearHistory(accountID int64) error {
	query := `DELETE account_history WHERE account_id = ?`
	_, err := r.db.Exec(query, accountID)

	if err != nil {
		return err
	}
	return nil
}

func (r *mysqlRepository) DeleteAccount(accountID int64) error {
	return nil
}
