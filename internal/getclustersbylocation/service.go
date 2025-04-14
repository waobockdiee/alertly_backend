package getclustersbylocation

type Service interface {
	GetClustersByLocation(inputs Inputs) ([]Cluster, error)
}

type service struct {
	repo Repository
}

func NewService(repo Repository) Service {
	return &service{repo: repo}
}

func (s *service) GetClustersByLocation(inputs Inputs) ([]Cluster, error) {
	result, err := s.repo.GetClustersByLocation(inputs)
	if err != nil {
		return []Cluster{}, err
	}
	return result, nil
}
