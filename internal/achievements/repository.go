package achievements

import (
	"database/sql"
)

type Repository interface {
	ShowByAccountID(accountID int64) ([]Achievement, error)
	Update(acacID, accountID int64) error
}
type pgRepository struct {
	db *sql.DB
}

func NewRepository(db *sql.DB) Repository {
	return &pgRepository{db: db}
}

func (r *pgRepository) ShowByAccountID(accountID int64) ([]Achievement, error) {
	query := `SELECT
	acac_id, account_id, name, description, created, show_in_modal, type, text_to_show, icon_url, badge_threshold
	FROM account_achievements
	WHERE account_id = $1 AND show_in_modal = 1
	ORDER BY created DESC`

	rows, err := r.db.Query(query, accountID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var achievements []Achievement
	for rows.Next() {
		var a Achievement
		if err := rows.Scan(
			&a.AcacID,
			&a.AccountID,
			&a.Name,
			&a.Description,
			&a.Created,
			&a.ShowInModal,
			&a.Type,
			&a.TextToShow,
			&a.IconUrl,
			&a.BadgeThreshold,
		); err != nil {
			return nil, err
		}
		achievements = append(achievements, a)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}
	return achievements, nil
}

func (r *pgRepository) Update(acacID, accountID int64) error {
	query := `UPDATE account_achievements SET show_in_modal = 0 WHERE acac_id = $1 AND account_id = $2`
	_, err := r.db.Exec(query, acacID, accountID)
	return err
}
