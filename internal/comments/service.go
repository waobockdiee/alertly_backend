package comments

type Service interface {
	Save(comment InComment) (int64, error)
	GetClusterCommentsByID(inclID int64) ([]Comment, error)
	GetCommentById(incoID int64) (Comment, error)
}

type service struct {
	repo Repository
}

func NewService(repo Repository) Service {
	return &service{repo: repo}
}

func (s *service) Save(comment InComment) (int64, error) {
	var err error
	var id int64

	id, err = s.repo.Save(comment)
	return id, err
}

func (s *service) GetClusterCommentsByID(inclID int64) ([]Comment, error) {
	var err error
	var comments []Comment

	comments, err = s.repo.GetClusterCommentsByID(inclID)
	return comments, err
}

func (s *service) GetCommentById(incoID int64) (Comment, error) {
	var err error
	var comment Comment

	comment, err = s.repo.GetCommentById(incoID)
	return comment, err
}
