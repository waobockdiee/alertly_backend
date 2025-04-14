package activate

import (
	"database/sql"
	"fmt"
)

type Repository interface {
	ActivateAccount(user ActivateAccountRequest) (sql.Result, error)
}

type mysqlRepository struct {
	db *sql.DB
}

func NewRepository(db *sql.DB) Repository {
	return &mysqlRepository{db: db}
}

func (r *mysqlRepository) ActivateAccount(user ActivateAccountRequest) (sql.Result, error) {
	query := `UPDATE account SET status = 'active' WHERE email = ? AND activation_code = ?`
	result, err := r.db.Exec(query, user.Email, user.ActivationCode)
	if err != nil {
		return nil, fmt.Errorf("error activando cuenta: %w", err)
	}
	return result, nil
}
