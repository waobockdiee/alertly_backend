package feedback

import (
	"database/sql"
)

type Repository interface {
	SendFeedback(feedback Feedback) error
}

type pgRepository struct {
	db *sql.DB
}

func NewRepository(db *sql.DB) Repository {
	return &pgRepository{db: db}
}

func (r *pgRepository) SendFeedback(feedback Feedback) error {
	query := `INSERT INTO feedback(account_id, subject, description) VALUES($1, $2, $3)`
	_, err := r.db.Exec(query, feedback.AccountID, feedback.Subject, feedback.Description)
	return err
}
