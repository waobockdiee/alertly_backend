package auth

import (
	"database/sql"
	"fmt"
)

type Repository interface {
	GetUserByEmail(email string) (User, error)
	GetUserById(accountID int64) (PasswordMatch, error)
}

type mysqlRepository struct {
	db *sql.DB
}

func NewRepository(db *sql.DB) Repository {
	return &mysqlRepository{db: db}
}

func (r *mysqlRepository) GetUserByEmail(email string) (User, error) {
	query := `
		SELECT account_id, email, password, phone_number, first_name, last_name, status, is_premium
		FROM account WHERE email = ?
	`
	row := r.db.QueryRow(query, email)
	var user User
	err := row.Scan(&user.AccountID, &user.Email, &user.Password, &user.PhoneNumber, &user.FirstName, &user.LastName, &user.Status, &user.IsPremium)
	return user, err
}

func (r *mysqlRepository) GetUserById(accountID int64) (PasswordMatch, error) {
	fmt.Println("ID:", accountID)
	query := `SELECT email, password FROM account WHERE account_id = ?`
	row := r.db.QueryRow(query, accountID)

	var pm PasswordMatch
	err := row.Scan(&pm.Email, &pm.Password)

	return pm, err
}
