package feedback

import (
	"database/sql"
)

type Repository interface {
	SendFeedback(feedback Feedback) error
}

type mysqlRepository struct {
	db *sql.DB
}

func NewRepository(db *sql.DB) Repository {
	return &mysqlRepository{db: db}
}

func (r *mysqlRepository) SendFeedback(feedback Feedback) error {
	query := `INSERT INTO feedback(account_id, subject, description) VALUES(?, ?, ?)`
	_, err := r.db.Exec(query, feedback.AccountID, feedback.Subject, feedback.Description)
	return err
}
