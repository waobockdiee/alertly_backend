package achievements

import (
	"database/sql"
)

type Repository interface {
	ShowByAccountID(accountID int64) ([]Achievement, error)
	Update(acacID, accountID int64) error
}
type mysqlRepository struct {
	db *sql.DB
}

func NewRepository(db *sql.DB) Repository {
	return &mysqlRepository{db: db}
}

func (r *mysqlRepository) ShowByAccountID(accountID int64) ([]Achievement, error) {
	query := `SELECT 
	acac_id, account_id, name, description, created, show_in_modal, type, text_to_show
	FROM account_achievements
	WHERE account_id = ? AND show = 1`
	rows, err := r.db.Query(query, accountID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var stcs []Achievement
	for rows.Next() {
		var c Achievement
		if err := rows.Scan(
			&c.AcacID,
			&c.AccountID,
			&c.Name,
			&c.Description,
			&c.Created,
			&c.ShowInModal,
			&c.Type,
			&c.TextToShow,
		); err != nil {
			return nil, err
		}
		stcs = append(stcs, c)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}
	return stcs, nil
}

func (r *mysqlRepository) Update(acacID, accountID int64) error {
	query := `UPDATE account_achievements SET show = 0 WHERE acac_id = ? AND account_id = ?`
}

func (r *mysqlRepository) Save(achievement Achievement) error {
	query := `INSERT INTO account_achievements(account_id, name, description, type, TextToShow, icon_url) VALUES()`
}
