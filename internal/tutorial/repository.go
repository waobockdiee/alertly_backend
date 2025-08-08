package tutorial

import "database/sql"

type Repository interface {
	MarkTutorialAsFinished(accountID int64) error
}

type mysqlRepository struct {
	db *sql.DB
}

func NewRepository(db *sql.DB) Repository {
	return &mysqlRepository{db: db}
}

func (r *mysqlRepository) MarkTutorialAsFinished(accountID int64) error {
	query := `UPDATE account SET has_finished_tutorial = 1 WHERE account_id = ?`
	_, err := r.db.Exec(query, accountID)
	return err
}
