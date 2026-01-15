package activate

import (
	"database/sql"
	"fmt"
)

type Repository interface {
	ActivateAccount(user ActivateAccountRequest) (sql.Result, error)
}

type pgRepository struct {
	db *sql.DB
}

func NewRepository(db *sql.DB) Repository {
	return &pgRepository{db: db}
}

func (r *pgRepository) ActivateAccount(user ActivateAccountRequest) (sql.Result, error) {
	query := `UPDATE account SET status = 'active' WHERE email = $1 AND activation_code = $2`
	result, err := r.db.Exec(query, user.Email, user.ActivationCode)
	if err != nil {
		return nil, fmt.Errorf("we couldn't activate your account. Please make sure your email and code are correct, then try again: %w", err)
	}
	return result, nil
}
