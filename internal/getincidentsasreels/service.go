// internal/getincidentsasreels/service.go

package getincidentsasreels

import (
	"alertly/internal/comments"
	"alertly/internal/common"
	"alertly/internal/database"
	"alertly/internal/getclusterby"
	"math"
)

type Service interface {
	GetReel(inputs Inputs, accountID int64) ([]getclusterby.Cluster, error)
}

type service struct {
	repo Repository
}

func NewService(repo Repository) Service {
	return &service{repo: repo}
}

func (s *service) GetReel(inputs Inputs, accountID int64) ([]getclusterby.Cluster, error) {
	// 1) obtenemos los clusters base
	clusters, err := s.repo.GetReel(inputs, accountID)
	if err != nil {
		return nil, err
	}

	// pre-inicializamos repositorios/servicios auxiliares
	cbRepo := getclusterby.NewRepository(database.DB)
	commentsRepo := comments.NewRepository(database.DB)
	commentsSvc := comments.NewService(commentsRepo)

	// 2) para cada uno, aplicamos la lógica extra
	for i := range clusters {
		c := &clusters[i]

		// 2.1 CredibilityPercent
		c.CredibilityPercent = calculateCredibilityPercent(
			c.CounterTotalVotesTrue,
			c.CounterTotalVotesFalse,
		)

		// 2.2 Votado y guardado por este usuario
		c.GetAccountAlreadyVoted, _ = cbRepo.GetAccountAlreadyVoted(c.InclId, accountID)
		c.GetAccountAlreadySaved, _ = cbRepo.GetAccountAlreadySaved(c.InclId, accountID)

		// 2.3 Comentarios
		c.Comments, _ = commentsSvc.GetClusterCommentsByID(c.InclId)

		// 2.4 TimeDiff para cada incidente
		for j := range c.Incidents {
			c.Incidents[j].TimeDiff = common.TimeAgo(
				c.Incidents[j].CreatedAt.Time,
			)
		}
	}

	return clusters, nil
}

// Misma función que tenías en getclusterby
func calculateCredibilityPercent(trueCount, falseCount int) float64 {
	total := float64(trueCount + falseCount)
	if total == 0 {
		return 0
	}
	percent := (float64(trueCount) / total) * 100
	return math.Round(percent*10) / 10
}
