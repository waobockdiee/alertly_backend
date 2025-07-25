package invitefriend

type Service interface {
	Save(invitation Invitation) error
}

type service struct {
	repo Repository
}

func NewService(repo Repository) Service {
	return &service{repo: repo}
}

func (s *service) Save(invitation Invitation) error {
	return s.repo.Save(invitation)
}
