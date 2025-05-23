package profile

type Service interface {
	GetById(accountID int64) (Profile, error)
	UpdateTotalIncidents(accountID int64) error
	ReportAccount(report ReportAccountInput) error
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
		return Range{"New Neighbor", "new_neighbor", "#CFD8DC", "#546E7A", ""}
	case score > 500 && score <= 1500:
		return Range{"Community Champion", "community_champion", "#AED581", "#33691E", ""}
	case score > 1500 && score <= 3000:
		return Range{"Neighborhood Legend", "neighborhood_legend", "#FFF176", "#F9A825", ""}
	case score > 3000 && score <= 6000:
		return Range{"Urban Guardian", "urban_guardian", "#4FC3F7", "#01579B", ""}
	case score > 6000 && score <= 10000:
		return Range{"Civic Hero", "civic_hero", "#FF8A65", "#BF360C", ""}
	case score > 10000 && score <= 20000:
		return Range{"Maple Leaf Icon", "maple_leaf_icon", "#E57373", "#B71C1C", ""}
	default:
		return Range{"Unknown", "unknown", "#CFD8DC", "#546E7A", ""}
	}
}

func (s *service) UpdateTotalIncidents(accountID int64) error {
	return s.repo.UpdateTotalIncidents(accountID)
}

func (s *service) ReportAccount(report ReportAccountInput) error {
	err := s.repo.ReportAccount(report)
	return err
}
