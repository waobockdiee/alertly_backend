package activate

import "errors"

type Service interface {
	ActivateAccount(user ActivateAccountRequest) error
}

type service struct {
	repo Repository
}

func NewService(repo Repository) Service {
	return &service{repo: repo}
}

func (s *service) ActivateAccount(user ActivateAccountRequest) error {
	result, err := s.repo.ActivateAccount(user)

	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()

	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		return errors.New("no se encontró usuario, código incorrecto o ya has activado tu cuenta")
	}
	return nil
}
