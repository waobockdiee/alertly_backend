package auth

import (
	"database/sql"
	"errors"
	"fmt"

	"alertly/internal/dbtypes"
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

	// Usar NullBool para escanear columnas booleanas que pueden ser SMALLINT/CHAR/BOOLEAN
	var isPremium, hasFinishedTutorial dbtypes.NullBool

	err := row.Scan(
		&user.AccountID,
		&user.Email,
		&user.Password,
		&user.PhoneNumber,
		&user.FirstName,
		&user.LastName,
		&user.Status,
		&isPremium,
		&hasFinishedTutorial,
	)

	// Normalizar error para no exponer detalles de implementación SQL
	if err == sql.ErrNoRows {
		fmt.Printf("❌ [AUTH-REPO] User not found in database: %s\n", email)
		return User{}, errors.New("invalid credentials")
	}
	if err != nil {
		fmt.Printf("❌ [AUTH-REPO] Database error for %s: %v\n", email, err)
		return User{}, err
	}

	// Convertir NullBool a bool (si es NULL, usar false por defecto)
	user.IsPremium = isPremium.Valid && isPremium.Bool
	user.HasFinishedTutorial = hasFinishedTutorial.Valid && hasFinishedTutorial.Bool

	fmt.Printf("✅ [AUTH-REPO] User found: %s (id: %d, status: %s, password_hash_len: %d)\n", email, user.AccountID, user.Status, len(user.Password))
	return user, nil
}

func (r *pgRepository) GetUserById(accountID int64) (PasswordMatch, error) {
	fmt.Println("ID:", accountID)
	query := `SELECT email, password FROM account WHERE account_id = $1`
	row := r.db.QueryRow(query, accountID)

	var pm PasswordMatch
	err := row.Scan(&pm.Email, &pm.Password)

	return pm, err
}
