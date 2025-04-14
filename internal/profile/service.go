package profile

type Service interface {
	GetById(accountID int64) (Profile, error)
	UpdateTotalIncidents(accountID int64) error
}

type service struct {
	repo Repository
}

func NewService(repo Repository) Service {
	return &service{repo: repo}
}

func (s *service) GetById(accountID int64) (Profile, error) {
	result, err := s.repo.GetById(accountID)
	result.Range = GetUserRange(result.Score)

	if err != nil {
		return Profile{}, err
	}

	return result, nil
}

func GetUserRange(score int) Range {
	switch {
	case score >= 0 && score <= 500:
		return Range{"New Neighbor", "#575dc6", "#575dc6", ""}
	case score > 500 && score <= 1500:
		return Range{"Community Champion", "#ffcc80", "#333333", ""}
	case score > 1500 && score <= 3000:
		return Range{"Neighborhood Legend", "#575dc6", "#e8e8f9", ""}
	case score > 3000 && score <= 6000:
		return Range{"Urban Guardian", "#575dc6", "#e8e8f9", ""}
	case score > 6000 && score <= 10000:
		return Range{"Civic Hero", "#575dc6", "#e8e8f9", ""}
	case score > 10000 && score <= 20000:
		return Range{"Maple Leaf Icon", "#575dc6", "#e8e8f9", ""}
	default:
		return Range{"Unknow", "#575dc6", "#e8e8f9", ""}
	}
}

func (s *service) UpdateTotalIncidents(accountID int64) error {
	return s.repo.UpdateTotalIncidents(accountID)
}
