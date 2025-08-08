package tutorial

type Service interface {
	FinishTutorial(accountID int64) error
}

type service struct {
	repo Repository
}

func NewService(repo Repository) Service {
	return &service{repo: repo}
}

func (s *service) FinishTutorial(accountID int64) error {
	return s.repo.MarkTutorialAsFinished(accountID)
}
