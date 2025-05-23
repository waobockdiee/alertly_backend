package saveclusteraccount

type Service interface {
	ToggleSaveClusterAccount(accountID, inclID int64) error
	GetMyList(accountID int64) ([]MyList, error)
	DeleteFollowIncident(acsID, accountID int64) error
}

type service struct {
	repo Repository
}

func NewService(repo Repository) Service {
	return &service{repo: repo}
}

func (s *service) ToggleSaveClusterAccount(accountID, inclID int64) error {
	var err error

	err = s.repo.ToggleSaveClusterAccount(accountID, inclID)
	return err
}

func (s *service) GetMyList(accountID int64) ([]MyList, error) {
	var myList []MyList
	var err error

	myList, err = s.repo.GetMyList(accountID)
	return myList, err
}

func (s *service) DeleteFollowIncident(acsID, accountID int64) error {
	err := s.repo.DeleteFollowIncident(acsID, accountID)
	return err
}
