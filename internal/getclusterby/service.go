package getclusterby

import (
	"alertly/internal/comments"
	"alertly/internal/common"
	"alertly/internal/database"
	"math"
)

type Service interface {
	GetIncidentBy(inclId, accountID int64) (Cluster, error)
}

type service struct {
	repo Repository
}

func NewService(repo Repository) Service {
	return &service{repo: repo}
}

func (s *service) GetIncidentBy(inclId, accountID int64) (Cluster, error) {
	result, err := s.repo.GetIncidentBy(inclId)
	if err != nil {
		return Cluster{}, err
	}
	result.CredibilityPercent = calculateCredibilityPercent(result.CounterTotalVotesTrue, result.CounterTotalVotesFalse)
	result.GetAccountAlreadyVoted, _ = s.repo.GetAccountAlreadyVoted(result.InclId, accountID)
	result.GetAccountAlreadySaved, _ = s.repo.GetAccountAlreadySaved(result.InclId, accountID)

	repo := comments.NewRepository(database.DB)
	cs := comments.NewService(repo)
	result.Comments, err = cs.GetClusterCommentsByID(result.InclId)

	for i := range result.Incidents {
		result.Incidents[i].TimeDiff = common.TimeAgo(result.Incidents[i].CreatedAt.Time)
	}

	if err != nil {
		return Cluster{}, err
	}
	return result, nil
}

func calculateCredibilityPercent(counterTotalVotesTrue, counterTotalVotesFake int) float64 {
	totalVotes := float64(counterTotalVotesTrue + counterTotalVotesFake)
	if totalVotes == 0 {
		return 0
	}
	credibilityPercent := (float64(counterTotalVotesTrue) / totalVotes) * 100
	credibilityPercent = math.Round(credibilityPercent*10) / 10
	return credibilityPercent
}
