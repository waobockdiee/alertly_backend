package getcategories

type Service interface {
	GetCategories() ([]Category, error)
}

type service struct {
	repo Repository
}

func NewService(repo Repository) Service {
	return &service{repo: repo}
}

func (s *service) GetCategories() ([]Category, error) {
	result, err := s.repo.GetCategories()
	if err != nil {
		return []Category{}, err
	}
	return result, nil
}
