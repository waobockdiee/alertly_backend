package account

import (
	"time"
)

type Service interface {
	GetMyInfo(accountID int64) (MyInfo, error)
	GetHistory(accountID int64) ([]History, error)
	GetViewedIncidentIds(accountID int64) ([]int64, error)
	ClearHistory(accountID int64) error
	DeleteAccount(accountID int64) error
	GetCounterHistories(accountID int64) (Counter, error)
	SaveLastRequest(accountID int64, ip string) error
	SetHasFinishedTutorial(accountID int64) error
	UpdatePremiumStatus(accountID int64, isPremium bool, subscriptionType string, purchaseDate *time.Time, platform string) error
}

type service struct {
	repo Repository
}

func NewService(repo Repository) Service {
	return &service{repo: repo}
}

func (s *service) GetMyInfo(accountID int64) (MyInfo, error) {
	return s.repo.GetMyInfo(accountID)
}

func (s *service) GetHistory(accountID int64) ([]History, error) {
	return s.repo.GetHistory(accountID)
}

func (s *service) GetViewedIncidentIds(accountID int64) ([]int64, error) {
	return s.repo.GetViewedIncidentIds(accountID)
}

func (s *service) ClearHistory(accountID int64) error {
	return s.repo.ClearHistory(accountID)
}

func (s *service) DeleteAccount(accountID int64) error {
	return nil
}

func (s *service) GetCounterHistories(accountID int64) (Counter, error) {
	return s.repo.GetCounterHistories(accountID)
}

func (s *service) SaveLastRequest(accountID int64, ip string) error {
	return s.repo.SaveLastRequest(accountID, ip)
}

func (s *service) SetHasFinishedTutorial(accountID int64) error {
	return s.repo.SetHasFinishedTutorial(accountID)
}

func (s *service) UpdatePremiumStatus(accountID int64, isPremium bool, subscriptionType string, purchaseDate *time.Time, platform string) error {
	return s.repo.UpdatePremiumStatus(accountID, isPremium, subscriptionType, purchaseDate, platform)
}
