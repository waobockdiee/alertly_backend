package feedback

type Service interface {
	SendFeedback(feedback Feedback) error
}

type service struct {
	repo Repository
}

func NewService(repo Repository) Service {
	return &service{repo: repo}
}

func (s *service) SendFeedback(feedback Feedback) error {
	return s.repo.SendFeedback(feedback)
}
