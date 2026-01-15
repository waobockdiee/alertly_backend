package alerts

import (
	"database/sql"
	"time"
)

type Repository interface {
	GetAlerts(accountID int64) ([]Alert, error)
	GetNewAlertsCount(accountID int64) (int64, error)
}

type pgRepository struct {
	db *sql.DB
}

func NewRepository(db *sql.DB) Repository {
	return &pgRepository{db: db}
}

func (r *pgRepository) GetAlerts(accountID int64) ([]Alert, error) {
	tx, err := r.db.Begin()
	var alerts []Alert

	if err != nil {
		return alerts, err
	}

	defer func() {
		if err != nil {
			_ = tx.Rollback()
		} else {
			_ = tx.Commit()
		}
	}()

	queryGetLastNotificationsDate := `SELECT last_notifications_view FROM account WHERE account_id = $1`

	var lastDate time.Time
	err = tx.QueryRow(queryGetLastNotificationsDate, accountID).Scan(&lastDate)

	if err != nil {
		return alerts, err
	}

	query := `SELECT *
	FROM account_notifications t1 INNER JOIN account
	LIMIT 20`
	rows, err := tx.Query(query, accountID, &lastDate)

	if err != nil {
		return alerts, err
	}

	defer rows.Close()

	for rows.Next() {
		var alert Alert

		if err := rows.Scan(); err != nil {
			return alerts, err
		}
		alerts = append(alerts, alert)
	}

	updateQuery := `UPDATE account SET last_notifications_view = NOW() WHERE account_id = $1`
	_, err = tx.Exec(updateQuery, accountID)

	if err != nil {
		return alerts, err
	}

	return alerts, nil

}

func (r *pgRepository) GetNewAlertsCount(accountID int64) (int64, error) {
	tx, err := r.db.Begin()

	if err != nil {
		return 0, err
	}

	defer func() {
		if err != nil {
			_ = tx.Rollback()
		} else {
			_ = tx.Commit()
		}
	}()
	queryGetLastNotificationsDate := `SELECT last_notifications_view FROM account WHERE account_id = $1`

	var lastDate time.Time
	err = tx.QueryRow(queryGetLastNotificationsDate, accountID).Scan(&lastDate)

	if err != nil {
		return 0, err
	}

	query := `SELECT COUNT(*) FROM account_notifications t1 WHERE created_at > $1`

	var counter int64

	err = tx.QueryRow(query, lastDate).Scan(&counter)

	if err != nil {
		return 0, err
	}

	return counter, nil
}
