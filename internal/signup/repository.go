package signup

import (
	"database/sql"
	"fmt"
)

type Repository interface {
	InsertUser(user User) (int64, error)
	GetUserByID(id int64) (User, error)
}

type pgRepository struct {
	db *sql.DB
}

func NewRepository(db *sql.DB) Repository {
	return &pgRepository{db: db}
}

func (repo *pgRepository) InsertUser(user User) (int64, error) {
	// Calcular días de trial según si tiene código de referral
	trialDays := 7 // Por defecto: 7 días
	if user.ReferralCode != "" {
		trialDays = 14 // Con código de referral: 14 días
	}

	query := `
		INSERT INTO account (email, first_name, last_name, password, activation_code, nickname,
		                     is_premium, premium_expired_date)
		VALUES ($1, $2, $3, $4, $5, $6, 1, NOW() + INTERVAL '1 day' * $7)
		RETURNING account_id
	`
	var id int64
	err := repo.db.QueryRow(query,
		user.Email,
		user.FirstName,
		user.LastName,
		user.Password,
		user.ActivationCode,
		user.Nickname,
		trialDays,
	).Scan(&id)
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

func (repo *pgRepository) GetUserByID(id int64) (User, error) {
	query := `
		SELECT account_id, email, password, activation_code, first_name
		FROM account WHERE account_id = $1
	`
	row := repo.db.QueryRow(query, id)
	var user User
	err := row.Scan(&user.ID, &user.Email, &user.Password, &user.ActivationCode, &user.FirstName)
	return user, err
}
