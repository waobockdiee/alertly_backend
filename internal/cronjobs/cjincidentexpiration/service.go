package cjincidentexpiration

import (
	"fmt"
	"log"
)

// Service defines the interface for the incident expiration cronjob service.
type Service interface {
	Run()
}

type service struct {
	repo Repository
}

// NewService creates a new instance of the service.
func NewService(repo Repository) Service {
	return &service{repo: repo}
}

const (
	credibilityThreshold = 5.0
	scoreWin           = 10.0
	scoreLoss          = -5.0
	credibilityWin     = 0.5
	credibilityLoss    = -0.2
)

// Run executes the main logic of the cronjob.
func (s *service) Run() {
	clusters, err := s.repo.GetExpiredClusters()
	if err != nil {
		log.Printf("Error getting expired clusters: %v", err)
		return
	}

	if len(clusters) == 0 {
		log.Println("No expired clusters to process.")
		return
	}

	log.Printf("Found %d expired clusters to process.", len(clusters))

	for _, cluster := range clusters {
		s.processCluster(cluster)
	}
}

func (s *service) processCluster(cluster ExpiredCluster) {
	// If credibility is NULL, we can't process it. Log and skip.
	if !cluster.Credibility.Valid {
		log.Printf("Skipping cluster ID: %d because its credibility is NULL.", cluster.ID)
		// Mark the cluster as processed even if credibility is null to avoid reprocessing
		if err := s.repo.MarkClusterProcessed(cluster.ID); err != nil {
			log.Printf("Error marking cluster %d with NULL credibility as processed: %v", cluster.ID, err)
		}
		return
	}

	finalCredibility := cluster.Credibility.Float64
	log.Printf("Processing cluster ID: %d with final credibility: %.1f", cluster.ID, finalCredibility)

	votes, err := s.repo.GetVotesForCluster(cluster.ID)
	if err != nil {
		log.Printf("Error getting votes for cluster %d: %v", cluster.ID, err)
		return
	}

	outcomeIsTrue := finalCredibility >= credibilityThreshold

	for _, vote := range votes {
		userVotedTrue := vote.Vote

		if userVotedTrue == outcomeIsTrue {
			// User was correct
			err := s.repo.UpdateUserStats(vote.AccountID, scoreWin, credibilityWin)
			if err != nil {
				log.Printf("Error updating stats for winning user %d: %v", vote.AccountID, err)
				continue // Skip to next vote
			}
			winMessage := fmt.Sprintf("Congratulations! Your vote on incident #%d was correct. You've earned +%.0f score points!", cluster.ID, scoreWin)
			if err := s.repo.SaveWinNotification(vote.AccountID, cluster.ID, winMessage); err != nil {
				log.Printf("Error creating win notification for user %d: %v", vote.AccountID, err)
			}
		} else {
			// User was incorrect
			err := s.repo.UpdateUserStats(vote.AccountID, scoreLoss, credibilityLoss)
			if err != nil {
				log.Printf("Error updating stats for losing user %d: %v", vote.AccountID, err)
				continue // Skip to next vote
			}
			lossMessage := fmt.Sprintf("Thanks for your input on incident #%d. This time it didn't match the final outcome, and your score has been updated.", cluster.ID)
			if err := s.repo.SaveLossNotification(vote.AccountID, cluster.ID, lossMessage); err != nil {
				log.Printf("Error creating loss notification for user %d: %v", vote.AccountID, err)
			}
		}
	}

	// Mark the cluster as processed to avoid re-processing
	if err := s.repo.MarkClusterProcessed(cluster.ID); err != nil {
		log.Printf("Error marking cluster %d as processed: %v", cluster.ID, err)
	}

	log.Printf("Finished processing cluster ID: %d", cluster.ID)
}
