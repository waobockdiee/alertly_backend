package auth

import (
	"database/sql"
)

type Repository interface {
	GetUserByEmail(email string) (User, error)
}

type mysqlRepository struct {
	db *sql.DB
}

func NewRepository(db *sql.DB) Repository {
	return &mysqlRepository{db: db}
}

func (repo *mysqlRepository) GetUserByEmail(email string) (User, error) {
	query := `
		SELECT account_id, email, password, phone_number, first_name, last_name, status, is_premium
		FROM account WHERE email = ?
	`
	row := repo.db.QueryRow(query, email)
	var user User
	err := row.Scan(&user.AccountID, &user.Email, &user.Password, &user.PhoneNumber, &user.FirstName, &user.LastName, &user.Status, &user.IsPremium)
	return user, err
}
