package editprofile

import (
	"database/sql"
	"errors"
	"fmt"
	"log"

	"alertly/internal/dbtypes"
)

type Repository interface {
	GetAccountByID(accountID int64) (Account, error)
	SaveCodeBeforeUpdateEmail(code string, accountID int64) error
	ValidateUpdateEmailCode(accountID int64, code string) (bool, error)
	UpdateThumbnail(accountID int64, mediaUrl string) error
	UpdateEmail(accountID int64, email string) error
	UpdatePassword(accountID int64, password string) error
	UpdateNickname(accountID int64, nickname string) error
	UpdatePhoneNumber(accountID int64, phoneNumber string) error
	UpdateFullName(accountID int64, firstName, lastName string) error
	UpdateIsPrivateProfile(accountID int64, isPrivateProfile bool) error

	UpdateBirthDate(accountID int64, year, month, day string) error
	UpdateReceiveNotifications(accountID int64) error
	DesactivateAccount(account Account) error
}

type pgRepository struct {
	db *sql.DB
}

func NewRepository(db *sql.DB) Repository {
	return &pgRepository{db: db}
}

func (r *pgRepository) GetAccountByID(accountID int64) (Account, error) {
	query := `SELECT
		account_id,
		email,
		first_name,
		last_name,
		nickname,
		password,
		can_update_nickname,
		can_update_fullname,
		can_update_birthdate,
		COALESCE(birth_year, '') as birth_year,
		COALESCE(birth_month, '') as birth_month,
		COALESCE(birth_day, '') as birth_day,
		can_update_email,
		COALESCE(thumbnail_url, '') as thumbnail_url,
		receive_notifications,
		status
	FROM account
	WHERE account_id = $1`

	row := r.db.QueryRow(query, accountID)

	var account Account

	// Usar NullBool para campos booleanos que pueden ser SMALLINT/CHAR/BOOLEAN
	var canUpdateNickname, canUpdateFullName, canUpdateBirthDate, canUpdateEmail, receiveNotifications dbtypes.NullBool

	err := row.Scan(
		&account.AccountID,
		&account.Email,
		&account.FirstName,
		&account.LastName,
		&account.NickName,
		&account.Password,
		&canUpdateNickname,
		&canUpdateFullName,
		&canUpdateBirthDate,
		&account.BirthYear,
		&account.BirthMonth,
		&account.BirthDay,
		&canUpdateEmail,
		&account.ThumbnailURL,
		&receiveNotifications,
		&account.Status,
	)

	if err != nil {
		log.Printf("Error scanning account data for ID %d: %v", accountID, err)
		return account, err
	}

	// Convertir NullBool a bool
	account.CanUpdateNickname = canUpdateNickname.Valid && canUpdateNickname.Bool
	account.CanUpdateFullName = canUpdateFullName.Valid && canUpdateFullName.Bool
	account.CanUpdateBirthDate = canUpdateBirthDate.Valid && canUpdateBirthDate.Bool
	account.CanUpdateEmail = canUpdateEmail.Valid && canUpdateEmail.Bool
	account.ReceiveNotifications = receiveNotifications.Valid && receiveNotifications.Bool

	return account, nil
}

func (r *pgRepository) SaveCodeBeforeUpdateEmail(code string, accountID int64) error {
	query := `UPDATE account SET update_email_code = $1 WHERE account_id = $2`
	_, err := r.db.Exec(query, code, accountID)
	return err
}

func (r *pgRepository) UpdateReceiveNotifications(accountID int64) error {
	query := `UPDATE account SET receive_notifications = NOT receive_notifications WHERE account_id = $1`
	_, err := r.db.Exec(query, accountID)
	return err
}

func (r *pgRepository) ValidateUpdateEmailCode(accountID int64, code string) (bool, error) {
	var match bool
	query := `SELECT EXISTS(
		SELECT 1
		FROM account
		WHERE account_id = $1 AND update_email_code = $2
	)
	`
	err := r.db.QueryRow(query, accountID, code).Scan(&match)

	if err != nil {
		return match, err
	}

	return match, nil
}

func (r *pgRepository) UpdateEmail(accountID int64, email string) error {

	log.Printf("account_id: %v", accountID)
	log.Printf("email: %v", email)
	query := `UPDATE account SET email = $1, can_update_email = $2 WHERE account_id = $3`
	res, err := r.db.Exec(query, email, dbtypes.BoolToInt(false), accountID)

	if err != nil {
		log.Printf("Error: %v", err)
		return err
	}

	n, err := res.RowsAffected()

	log.Printf("n: %v", n)

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

func (r *pgRepository) UpdateThumbnail(accountID int64, mediaUrl string) error {
	query := `UPDATE account SET thumbnail_url = $1 WHERE account_id = $2`
	_, err := r.db.Exec(query, mediaUrl, accountID)

	if err != nil {
		log.Printf("Error: %v", err)
	}

	return err
}

func (r *pgRepository) UpdatePassword(accountID int64, password string) error {
	query := `UPDATE account SET password = $1 WHERE account_id = $2`
	_, err := r.db.Exec(query, password, accountID)

	if err != nil {
		log.Printf("Error: %v", err)
	}

	return err
}

func (r *pgRepository) UpdateNickname(accountID int64, nickname string) error {
	query := `UPDATE account SET nickname = $1, can_update_nickname = $2 WHERE account_id = $3 AND can_update_nickname = $4`
	_, err := r.db.Exec(query, nickname, dbtypes.BoolToInt(false), accountID, dbtypes.BoolToInt(true))

	if err != nil {
		log.Printf("Error: %v", err)
	}

	return err
}

func (r *pgRepository) UpdatePhoneNumber(accountID int64, phoneNumber string) error {
	query := `UPDATE account SET phone_number = $1 WHERE account_id = $2`
	_, err := r.db.Exec(query, phoneNumber, accountID)

	if err != nil {
		log.Printf("Error: %v", err)
	}

	return err
}

func (r *pgRepository) UpdateFullName(accountID int64, firstName, lastName string) error {
	fmt.Println("account_id", accountID)
	fmt.Println("first_name", firstName)
	fmt.Println("first_name", lastName)
	query := `UPDATE account SET first_name = $1, last_name = $2, can_update_fullname = $3 WHERE account_id = $4 AND can_update_fullname = $5`
	_, err := r.db.Exec(query, firstName, lastName, dbtypes.BoolToInt(false), accountID, dbtypes.BoolToInt(true))

	if err != nil {
		log.Printf("Error: %v", err)
	}

	return err
}

func (r *pgRepository) UpdateIsPrivateProfile(accountID int64, isPrivateProfile bool) error {
	query := `UPDATE account SET is_private_profile = $1 WHERE account_id = $2`
	_, err := r.db.Exec(query, isPrivateProfile, accountID)

	if err != nil {
		log.Printf("Error: %v", err)
	}

	return err
}



func (r *pgRepository) UpdateBirthDate(accountID int64, year, month, day string) error {
	query := `UPDATE account SET birth_year = $1, birth_month = $2, birth_day = $3, can_update_birthdate = $4 WHERE account_id = $5 AND can_update_birthdate = $6`
	_, err := r.db.Exec(query, year, month, day, dbtypes.BoolToInt(false), accountID, dbtypes.BoolToInt(true))

	if err != nil {
		log.Printf("Error: %v", err)
	}

	return err
}

func (r *pgRepository) DesactivateAccount(account Account) error {

	query := `UPDATE account SET status = $1 WHERE account_id = $2`
	_, err := r.db.Exec(query, account.Status, account.AccountID)

	if err != nil {
		log.Printf("Error: %v", err)
		return err
	}

	return nil
}
