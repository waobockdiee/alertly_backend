package saveclusteraccount

import (
	"database/sql"
)

type Repository interface {
	ToggleSaveClusterAccount(accountID, inclID int64) error
	GetMyList(accountID int64) ([]MyList, error)
	DeleteFollowIncident(acsID, accountID int64) error
}

type pgRepository struct {
	db *sql.DB
}

func NewRepository(db *sql.DB) Repository {
	return &pgRepository{db: db}
}

func (r *pgRepository) ToggleSaveClusterAccount(accountID, inclID int64) (err error) {
	tx, err := r.db.Begin()
	if err != nil {
		return err
	}
	defer func() {
		if err != nil {
			tx.Rollback()
		}
	}()

	var total int
	queryCheck := `SELECT COUNT(*) AS total FROM account_cluster_saved WHERE account_id = $1 AND incl_id = $2`
	err = tx.QueryRow(queryCheck, accountID, inclID).Scan(&total)
	if err != nil {
		return err
	}

	if total > 0 {
		queryDelete := `DELETE FROM account_cluster_saved WHERE account_id = $1 AND incl_id = $2`
		_, err = tx.Exec(queryDelete, accountID, inclID)
		if err != nil {
			return err
		}
	} else {
		queryInsert := `INSERT INTO account_cluster_saved (account_id, incl_id) VALUES ($1, $2)`
		_, err = tx.Exec(queryInsert, accountID, inclID)
		if err != nil {
			return err
		}
	}

	return tx.Commit()
}

func (r *pgRepository) GetMyList(accountID int64) ([]MyList, error) {

	var list []MyList

	query := `SELECT
	t1.acs_id, t1.account_id, t1.incl_id, t2.media_url, t2.credibility
	FROM account_cluster_saved t1 INNER JOIN incident_clusters t2 ON t1.incl_id = t2.incl_id
	WHERE t1.account_id = $1`

	rows, err := r.db.Query(query, accountID)

	if err != nil {
		return list, err
	}

	defer rows.Close()

	for rows.Next() {
		var myList MyList

		if err := rows.Scan(&myList.AcsID, &myList.AccountID, &myList.InclID, &myList.MediaUrl, &myList.Credibility); err != nil {
			return list, err
		}

		list = append(list, myList)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return list, nil
}

func (r *pgRepository) DeleteFollowIncident(acsID, accountID int64) error {
	query := `DELETE FROM account_cluster_saved WHERE acs_id = $1 AND account_id = $2`
	_, err := r.db.Exec(query, acsID, accountID)
	return err
}
