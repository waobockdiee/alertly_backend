package alerts

type Service interface {
	GetAlerts(accountID int64) ([]Alert, error)
	GetNewAlertsCount(accountID int64) (int64, error)
}

type service struct {
	repo Repository
}

func NewService(repo Repository) Service {
	return &service{repo: repo}
}

func (s *service) GetAlerts(accountID int64) ([]Alert, error) {
	var alerts []Alert
	return alerts, nil
}

func (s *service) GetNewAlertsCount(accountID int64) (int64, error) {
	return s.repo.GetNewAlertsCount(accountID)
}
