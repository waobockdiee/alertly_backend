package account

import (
	"alertly/internal/database"
	"alertly/internal/myplaces"
	"log"
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
	SetHasFinishedTutorial(accountID int64, latitude, longitude *float32) error
	UpdatePremiumStatus(accountID int64, isPremium bool, subscriptionType string, expirationDate *time.Time, platform string) error
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

func (s *service) SetHasFinishedTutorial(accountID int64, latitude, longitude *float32) error {
	// 1. Mark tutorial as finished (main operation)
	if err := s.repo.SetHasFinishedTutorial(accountID); err != nil {
		return err
	}

	// 2. Create initial place if coordinates are provided (optional - don't fail tutorial if this fails)
	if latitude != nil && longitude != nil && *latitude != 0 && *longitude != 0 {
		if err := s.createInitialPlace(accountID, *latitude, *longitude); err != nil {
			log.Printf("Could not create initial place for user %d: %v", accountID, err)
			// Continue - tutorial completion is more important than place creation
		}
	}

	return nil
}

func (s *service) createInitialPlace(accountID int64, latitude, longitude float32) error {
	myPlace := myplaces.MyPlaces{
		AccountId:                 accountID,
		Title:                     "My Place",
		Latitude:                  latitude,
		Longitude:                 longitude,
		City:                      "", // Will be empty - skip reverse geocoding for simplicity
		Province:                  "",
		PostalCode:                "",
		Status:                    true,
		Crime:                     true,
		TrafficAccident:           true,
		MedicalEmergency:          true,
		FireIncident:              true,
		Vandalism:                 true,
		SuspiciousActivity:        true,
		InfrastructureIssues:      true,
		ExtremeWeather:            true,
		CommunityEvents:           true,
		DangerousWildlifeSighting: true,
		PositiveActions:           true,
		LostPet:                   true,
		Radius:                    5000, // 5km fixed radius
	}

	// Create myplaces repository and add the place
	myplacesRepo := myplaces.NewRepository(database.DB)
	_, err := myplacesRepo.Add(myPlace)
	return err
}

func (s *service) UpdatePremiumStatus(accountID int64, isPremium bool, subscriptionType string, expirationDate *time.Time, platform string) error {
	return s.repo.UpdatePremiumStatus(accountID, isPremium, subscriptionType, expirationDate, platform)
}
