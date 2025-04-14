package myplaces

import "alertly/internal/common"

type Service interface {
	Get(accountId int) ([]MyPlaces, error)
	Add(myPlace MyPlaces) (int64, error)
	Update(myPlace MyPlaces) error
	GetByAccountId(accountId int64) ([]MyPlaces, error)
	GetById(accountId, aflId int64) (MyPlaces, error)
	FullUpdate(myPlace MyPlaces) error
	Delete(accountID, aflID int64) error
}

type service struct {
	repo Repository
}

func NewService(repo Repository) Service {
	return &service{repo: repo}
}

func (s *service) Get(accountId int) ([]MyPlaces, error) {
	result, err := s.repo.Get(accountId)
	if err != nil {
		return []MyPlaces{}, err
	}
	return result, nil
}

func (s *service) Add(myPlace MyPlaces) (int64, error) {

	_, city, province, postalCode, errGeo := common.ReverseGeocode(myPlace.Latitude, myPlace.Longitude)

	if errGeo != nil {
		return 0, errGeo
	}

	myPlace.City = city
	myPlace.Province = province
	myPlace.PostalCode = postalCode

	result, err := s.repo.Add(myPlace)

	if err != nil {
		return 0, err
	}
	return result, err
}

func (s *service) Update(myPlace MyPlaces) error {
	err := s.repo.Update(myPlace)
	return err
}

func (s *service) FullUpdate(myPlace MyPlaces) error {
	err := s.repo.FullUpdate(myPlace)
	return err
}

func (s *service) GetByAccountId(accountId int64) ([]MyPlaces, error) {
	data, err := s.repo.GetByAccountId(accountId)
	return data, err
}

func (s *service) GetById(accountId, aflId int64) (MyPlaces, error) {
	data, err := s.repo.GetById(accountId, aflId)
	return data, err
}

func (s *service) Delete(accountID, aflID int64) error {
	err := s.repo.Delete(accountID, aflID)
	return err
}
