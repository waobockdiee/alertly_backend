package saveclusteraccount

import "database/sql"

type Repository interface {
	ToggleSaveClusterAccount(accountID, inclID int64) error
	GetMyList(accountID int64) ([]MyList, error)
}

type mysqlRepository struct {
	db *sql.DB
}

func NewRepository(db *sql.DB) Repository {
	return &mysqlRepository{db: db}
}

func (r *mysqlRepository) ToggleSaveClusterAccount(accountID, inclID int64) (err error) {
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
	queryCheck := `SELECT COUNT(*) AS total FROM account_cluster_saved WHERE account_id = ? AND incl_id = ?`
	err = tx.QueryRow(queryCheck, accountID, inclID).Scan(&total)
	if err != nil {
		return err
	}

	if total > 0 {
		queryDelete := `DELETE FROM account_cluster_saved WHERE account_id = ? AND incl_id = ?`
		_, err = tx.Exec(queryDelete, accountID, inclID)
		if err != nil {
			return err
		}
	} else {
		queryInsert := `INSERT INTO account_cluster_saved (account_id, incl_id) VALUES (?, ?)`
		_, err = tx.Exec(queryInsert, accountID, inclID)
		if err != nil {
			return err
		}
	}

	return tx.Commit()
}

func (r *mysqlRepository) GetMyList(accountID int64) ([]MyList, error) {

	var list []MyList

	query := `SELECT 
	t1.acs_id, t1.account_id, t1.incl_id, t2.media_url
	FROM account_cluster_saved t1 INNER JOIN incident_clusters t2 ON t1.incl_id = t2.incl_id
	WHERE t1.account_id = ?`

	rows, err := r.db.Query(query, accountID)

	if err != nil {
		return list, err
	}

	defer rows.Close()

	for rows.Next() {
		var myList MyList

		if err := rows.Scan(&myList.AcsID, &myList.AccountID, &myList.InclID, &myList.MediaUrl); err != nil {
			return list, err
		}

		list = append(list, myList)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return list, nil
}
