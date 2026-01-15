package auth

import (
	"database/sql"
	"errors"
	"fmt"
)

type Repository interface {
	GetUserByEmail(email string) (User, error)
	GetUserById(accountID int64) (PasswordMatch, error)
}

type pgRepository struct {
	db *sql.DB
}

func NewRepository(db *sql.DB) Repository {
	return &pgRepository{db: db}
}

func (r *pgRepository) GetUserByEmail(email string) (User, error) {
	query := `
		SELECT account_id, email, password, phone_number, first_name, last_name, status, is_premium, has_finished_tutorial
		FROM account WHERE email = $1
	`
	row := r.db.QueryRow(query, email)
	var user User
	err := row.Scan(&user.AccountID, &user.Email, &user.Password, &user.PhoneNumber, &user.FirstName, &user.LastName, &user.Status, &user.IsPremium, &user.HasFinishedTutorial)

	// Normalizar error para no exponer detalles de implementación SQL
	if err == sql.ErrNoRows {
		return User{}, errors.New("invalid credentials")
	}

	return user, err
}

func (r *pgRepository) GetUserById(accountID int64) (PasswordMatch, error) {
	fmt.Println("ID:", accountID)
	query := `SELECT email, password FROM account WHERE account_id = $1`
	row := r.db.QueryRow(query, accountID)

	var pm PasswordMatch
	err := row.Scan(&pm.Email, &pm.Password)

	return pm, err
}
