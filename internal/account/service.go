package account

type Service interface {
	GetMyInfo(accountID int64) (MyInfo, error)
	GetHistory(accountID int64) ([]History, error)
	ClearHistory(accountID int64) error
	DeleteAccount(accountID int64) error
	GetCounterHistories(accountID int64) (Counter, error)
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

func (s *service) ClearHistory(accountID int64) error {
	return s.repo.ClearHistory(accountID)
}

func (s *service) DeleteAccount(accountID int64) error {
	return nil
}

func (s *service) GetCounterHistories(accountID int64) (Counter, error) {
	return s.repo.GetCounterHistories(accountID)
}
