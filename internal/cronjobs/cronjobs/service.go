package cronjobs

import "fmt"

type Service interface {
	SetClusterToInactiveAndSetAccountScore()
}

type service struct {
	repo Repository
}

func NewService(repo Repository) Service {
	return &service{repo: repo}
}

func (s *service) NewCluster() {
	
}

func (s *service) SetClusterToInactiveAndSetAccountScore() {
	err := s.repo.SetClusterToInactiveAndSetAccountScore()

	if err != nil {
		fmt.Printf("error in SetClusterToInactiveAndSetAccountScore: %v", err)
	}
}
