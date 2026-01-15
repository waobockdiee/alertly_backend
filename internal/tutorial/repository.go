package tutorial

import "database/sql"

type Repository interface {
	MarkTutorialAsFinished(accountID int64) error
}

type pgRepository struct {
	db *sql.DB
}

func NewRepository(db *sql.DB) Repository {
	return &pgRepository{db: db}
}

func (r *pgRepository) MarkTutorialAsFinished(accountID int64) error {
	query := `UPDATE account SET has_finished_tutorial = 1 WHERE account_id = $1`
	_, err := r.db.Exec(query, accountID)
	return err
}
