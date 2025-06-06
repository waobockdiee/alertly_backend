package editprofile

import (
	"alertly/internal/common"
	"alertly/internal/emails"
	"errors"
	"log"

	"golang.org/x/crypto/bcrypt"
)

type Service interface {
	GetAccountByID(accountID int64) (Account, error)
	GenerateCodeUpdateEmail(accountID int64) error
	ValidateUpdateEmailCode(accountID int64, code string) (bool, error)
	UpdateEmail(accountID int64, email string) error
	UpdatePassword(accountID int64, newPassword string) error
	UpdateNickname(accountID int64, nickname string) error
	UpdatePhoneNumber(accountID int64, phoneNumber string) error
	UpdateFullName(accountID int64, firstName, lastName string) error
	UpdateIsPrivateProfile(accountID int64, isPrivateProfile bool) error
	UpdateIsPremium(accountID int64, isPremium bool) error
	UpdateBirthDate(accountID int64, year, month, day string) error
	CheckPasswordMatch(password, newPassword string, accountID int64) error
	// UpdateThumbnail(accountID int64) error
}

type service struct {
	repo Repository
}

func NewService(repo Repository) Service {
	return &service{repo: repo}
}

func (s *service) GetAccountByID(accountID int64) (Account, error) {
	return s.repo.GetAccountByID(accountID)
}

func (s *service) GenerateCodeUpdateEmail(accountID int64) error {

	var account Account
	var err error

	account, err = s.repo.GetAccountByID(accountID)

	if err != nil {
		log.Printf("Error in editprofile/service.go GetAccountByID: %v", err)
		return err
	}

	account.Code, err = common.GenerateCode()

	if err != nil {
		log.Printf("Error in editprofile/service.go GenerateCode: %v", err)
		return err
	}

	err = s.repo.SaveCodeBeforeUpdateEmail(account.Code, accountID)

	if err != nil {
		log.Printf("Error in editprofile/service.go SaveCodeBeforeUpdateEmail: %v", err)
		return err
	}

	body :=
		`
		<div style="margin-bottom: 20px;"}>Hi` + account.FirstName + `,</div>
		<div style="margin-bottom: 20px">We received a request to change the email address for your Alertly account.  
		To confirm this change and save your new email, please enter the verification code below in the app:
		</div>
		<div style="display: flex; gap: 20px;">
			<div>Vertification code:</div>
			<div style="background-color: #1e1c3b; color: #fff; font-size: 18px; padding: 15px 30px; border-radius: 10px;">` + account.Code + `</div>
		</div>
		`

	emails.Send(account.Email, "Code verification", body)
	return nil
}

func (s *service) ValidateUpdateEmailCode(accountID int64, code string) (bool, error) {
	return s.repo.ValidateUpdateEmailCode(accountID, code)
}

func (s *service) UpdateEmail(accountID int64, email string) error {
	return s.repo.UpdateEmail(accountID, email)
}

func (s *service) UpdateNickname(accountID int64, nickname string) error {
	return s.repo.UpdateNickname(accountID, nickname)
}

func (s *service) UpdatePhoneNumber(accountID int64, phoneNumber string) error {
	return s.repo.UpdatePhoneNumber(accountID, phoneNumber)
}

func (s *service) UpdateFullName(accountID int64, firstName, lastName string) error {
	return s.repo.UpdateFullName(accountID, firstName, lastName)
}

func (s *service) UpdateIsPrivateProfile(accountID int64, isPrivateProfile bool) error {
	return s.repo.UpdateIsPrivateProfile(accountID, isPrivateProfile)
}

func (s *service) UpdateIsPremium(accountID int64, isPremium bool) error {
	return s.repo.UpdateIsPremium(accountID, isPremium)
}

func (s *service) UpdateBirthDate(accountID int64, year, month, day string) error {
	return s.repo.UpdateBirthDate(accountID, year, month, day)
}

func (s *service) CheckPasswordMatch(password, newPassword string, accountID int64) error {
	var account Account
	var err error

	account, err = s.repo.GetAccountByID(accountID)

	if err != nil {
		log.Printf("Error in CheckPasswordMatch editprofile/service.go GetAccountByID: %v", err)
		return err
	}

	if err := bcrypt.CompareHashAndPassword([]byte(account.Password), []byte(password)); err != nil {
		return errors.New("wrong password")
	}

	return nil
}

func (s *service) UpdatePassword(accountID int64, newPassword string) error {

	var password string

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(newPassword), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	password = string(hashedPassword)

	return s.repo.UpdatePassword(accountID, password)
}

// func (s *service) UpdateThumbnail(accountID int64) error {
// 	return s.repo.UpdateThumbnail(accountID)
// }
