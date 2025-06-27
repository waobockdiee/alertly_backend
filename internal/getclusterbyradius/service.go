package getclusterbyradius

type Service interface {
	GetClustersByRadius(inputs Inputs) ([]Cluster, error)
}

type service struct {
	repo Repository
}

func NewService(repo Repository) Service {
	return &service{repo: repo}
}

func (s *service) GetClustersByRadius(inputs Inputs) ([]Cluster, error) {
	result, err := s.repo.GetClustersByRadius(inputs)
	if err != nil {
		return []Cluster{}, err
	}
	return result, nil
}
