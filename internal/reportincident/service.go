package reportincident

type Service interface {
	ReportIncident(report Report) error
}

type service struct {
	repo Repository
}

func NewService(repo Repository) Service {
	return &service{repo: repo}
}

func (s *service) ReportIncident(report Report) error {
	return s.repo.ReportIncident(report)
}
