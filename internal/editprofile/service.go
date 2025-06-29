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
	UpdateThumbnail(accountID int64, mediaUrl string) error
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

	emails.SendTemplate(account.Email, "Code verification", "update_email_verification_code", map[string]string{
		"Name": account.FirstName,
		"Code": account.Code,
	})
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
		log.Printf("Error CompareHashAndPassword editprofile service.go: %v", err)
		return errors.New("incorrect password. Please try again")
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

func (s *service) UpdateThumbnail(accountID int64, mediaUrl string) error {
	return s.repo.UpdateThumbnail(accountID, mediaUrl)
}
