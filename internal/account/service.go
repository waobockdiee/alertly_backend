package account

type Service interface {
	GetHistory(accountID int64) ([]History, error)
	ClearHistory(accountID int64) error
	DeleteAccount(accountID int64) error
}

type service struct {
	repo Repository
}

func NewService(repo Repository) Service {
	return &service{repo: repo}
}

func (s *service) GetHistory(accountID int64) ([]History, error) {
	return s.repo.GetHistory(accountID)
}

func (s *service) ClearHistory(accountID int64) error {
	return s.repo.ClearHistory(accountID)
}

func (s *service) DeleteAccount(accountID int64) error {
	return nil
}
