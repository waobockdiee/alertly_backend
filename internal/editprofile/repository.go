package editprofile

import (
	"database/sql"
	"errors"
	"fmt"
	"log"
)

type Repository interface {
	GetAccountByID(accountID int64) (Account, error)
	SaveCodeBeforeUpdateEmail(code string, accountID int64) error
	ValidateUpdateEmailCode(accountID int64, code string) (bool, error)
	UpdateThumbnail(accountID int64, media Media) error
	UpdateEmail(accountID int64, email string) error
	UpdatePassword(accountID int64, password string) error
	UpdateNickname(accountID int64, nickname string) error
	UpdatePhoneNumber(accountID int64, phoneNumber string) error
	UpdateFullName(accountID int64, firstName, lastName string) error
	UpdateIsPrivateProfile(accountID int64, isPrivateProfile bool) error
	UpdateIsPremium(accountID int64, isPremium bool) error
	UpdateBirthDate(accountID int64, year, month, day string) error
}

type mysqlRepository struct {
	db *sql.DB
}

func NewRepository(db *sql.DB) Repository {
	return &mysqlRepository{db: db}
}

func (r *mysqlRepository) GetAccountByID(accountID int64) (Account, error) {
	query := `SELECT account_id, email, first_name, last_name, password, can_update_nickname, can_update_fullname, can_update_birthdate FROM account WHERE account_id = ?`
	row := r.db.QueryRow(query, accountID)

	var account Account

	err := row.Scan(&account.AccountID, &account.Email, &account.FirstName, &account.LastName, &account.Password, &account.CanUpdateNickname, &account.CanUpdateFullName, &account.CanUpdateBirthDate)
	return account, err
}

func (r *mysqlRepository) SaveCodeBeforeUpdateEmail(code string, accountID int64) error {
	query := `UPDATE account SET update_email_code = ? WHERE account_id = ?`
	_, err := r.db.Exec(query, accountID)
	return err
}

func (r *mysqlRepository) ValidateUpdateEmailCode(accountID int64, code string) (bool, error) {
	var match bool
	query := `SELECT EXISTS(
		SELECT 1
		FROM account
		WHERE account_id = ? AND update_email_code = ?
	)
	`
	err := r.db.QueryRow(query, accountID, code).Scan(&match)

	if err != nil {
		return match, err
	}

	return match, nil
}

func (r *mysqlRepository) UpdateEmail(accountID int64, email string) error {
	query := `UPDATE account SET email = ? WHERE account_id = ?`
	res, err := r.db.Exec(query, email, accountID)

	if err != nil {
		log.Printf("Error: %v", err)
		return err
	}

	n, err := res.RowsAffected()

	if err != nil {
		log.Printf("Error: %v", err)
		return err
	}

	if n == 0 {
		err := errors.New("wrong code")
		return err
	}

	return nil
}

func (r *mysqlRepository) UpdateThumbnail(accountID int64, media Media) error {
	query := `UPDATE account SET thumbnail_url = ? WHERE account_id = ?`
	_, err := r.db.Exec(query, media.Uri, accountID)

	if err != nil {
		log.Printf("Error: %v", err)
	}

	return err
}

func (r *mysqlRepository) UpdatePassword(accountID int64, password string) error {
	query := `UPDATE account SET password = ? WHERE account_id = ?`
	_, err := r.db.Exec(query, password, accountID)

	if err != nil {
		log.Printf("Error: %v", err)
	}

	return err
}

func (r *mysqlRepository) UpdateNickname(accountID int64, nickname string) error {
	query := `UPDATE account SET nickname = ?, can_update_nickname = 0 WHERE account_id = ? AND can_update_nickname = 1`
	_, err := r.db.Exec(query, nickname, accountID)

	if err != nil {
		log.Printf("Error: %v", err)
	}

	return err
}

func (r *mysqlRepository) UpdatePhoneNumber(accountID int64, phoneNumber string) error {
	query := `UPDATE account SET phone_number = ? WHERE account_id = ?`
	_, err := r.db.Exec(query, phoneNumber, accountID)

	if err != nil {
		log.Printf("Error: %v", err)
	}

	return err
}

func (r *mysqlRepository) UpdateFullName(accountID int64, firstName, lastName string) error {
	fmt.Println("account_id", accountID)
	fmt.Println("first_name", firstName)
	fmt.Println("first_name", lastName)
	query := `UPDATE account SET first_name = ?, last_name = ?, can_update_fullname = 0 WHERE account_id = ? AND can_update_fullname = 1`
	_, err := r.db.Exec(query, firstName, lastName, accountID)

	if err != nil {
		log.Printf("Error: %v", err)
	}

	return err
}

func (r *mysqlRepository) UpdateIsPrivateProfile(accountID int64, isPrivateProfile bool) error {
	query := `UPDATE account SET is_private_profile = ? WHERE account_id = ?`
	_, err := r.db.Exec(query, isPrivateProfile, accountID)

	if err != nil {
		log.Printf("Error: %v", err)
	}

	return err
}

func (r *mysqlRepository) UpdateIsPremium(accountID int64, isPremium bool) error {
	query := `UPDATE account SET is_premium = ? WHERE account_id = ?`
	_, err := r.db.Exec(query, isPremium, accountID)

	if err != nil {
		log.Printf("Error: %v", err)
	}

	return err
}

func (r *mysqlRepository) UpdateBirthDate(accountID int64, year, month, day string) error {
	query := `UPDATE account SET birth_year = ?, birth_month = ?, birth_day = ?, can_update_birthdate = 0 WHERE account_id = ? AND can_update_birthdate = 1`
	_, err := r.db.Exec(query, year, month, day, accountID)

	if err != nil {
		log.Printf("Error: %v", err)
	}

	return err
}
