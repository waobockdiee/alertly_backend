package achievements

type Service interface {
	GetPendingByAccountID(accountID int64) ([]Achievement, error)
	MarkAsShown(acacID, accountID int64) error
}

type service struct {
	repo Repository
}

func NewService(repo Repository) Service {
	return &service{repo: repo}
}

func (s *service) GetPendingByAccountID(accountID int64) ([]Achievement, error) {
	return s.repo.ShowByAccountID(accountID)
}

func (s *service) MarkAsShown(acacID, accountID int64) error {
	return s.repo.Update(acacID, accountID)
}
