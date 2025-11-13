package signup

import (
	"database/sql"
	"fmt"
)

type Repository interface {
	InsertUser(user User) (int64, error)
	GetUserByID(id int64) (User, error)
}

type mysqlRepository struct {
	db *sql.DB
}

func NewRepository(db *sql.DB) Repository {
	return &mysqlRepository{db: db}
}

func (repo *mysqlRepository) InsertUser(user User) (int64, error) {
	// Calcular días de trial según si tiene código de referral
	trialDays := 7 // Por defecto: 7 días
	if user.ReferralCode != "" {
		trialDays = 14 // Con código de referral: 14 días
	}

	query := `
		INSERT INTO account (email, first_name, last_name, password, activation_code, nickname,
		                     is_premium, premium_expired_date)
		VALUES (?, ?, ?, ?, ?, ?, 1, DATE_ADD(NOW(), INTERVAL ? DAY))
	`
	result, err := repo.db.Exec(query,
		user.Email,
		user.FirstName,
		user.LastName,
		user.Password,
		user.ActivationCode,
		user.Nickname,
		trialDays,
	)
	if err != nil {
		return 0, err
	}
	id, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}
	if user.ReferralCode != "" {
		fmt.Printf("✅ Usuario insertado con ID: %d (con código de referral: %s, trial: %d días)\n", id, user.ReferralCode, trialDays)
	} else {
		fmt.Printf("Usuario insertado con ID: %d (sin código de referral, trial: %d días)\n", id, trialDays)
	}
	return id, nil
}

func (repo *mysqlRepository) GetUserByID(id int64) (User, error) {
	query := `
		SELECT account_id, email, password, activation_code, first_name
		FROM account WHERE account_id = ?
	`
	row := repo.db.QueryRow(query, id)
	var user User
	err := row.Scan(&user.ID, &user.Email, &user.Password, &user.ActivationCode, &user.FirstName)
	return user, err
}
