package getsubcategoriesbycategoryid

type Service interface {
	GetSubcategoriesByCategoryId(id int) ([]Subcategory, error)
}

type service struct {
	repo Repository
}

func NewService(repo Repository) Service {
	return &service{repo: repo}
}

func (s *service) GetSubcategoriesByCategoryId(id int) ([]Subcategory, error) {
	result, err := s.repo.GetSubcategoriesByCategoryId(id)
	if err != nil {
		return []Subcategory{}, err
	}
	return result, nil
}
