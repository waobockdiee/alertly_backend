package tutorial

import (
	"alertly/internal/database"
	"alertly/internal/myplaces"
	"log"
)

type Service interface {
	FinishTutorial(accountID int64, latitude, longitude *float32) error
}

type service struct {
	repo Repository
}

func NewService(repo Repository) Service {
	return &service{repo: repo}
}

func (s *service) FinishTutorial(accountID int64, latitude, longitude *float32) error {
	// 1. Mark tutorial as finished (main operation)
	if err := s.repo.MarkTutorialAsFinished(accountID); err != nil {
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
