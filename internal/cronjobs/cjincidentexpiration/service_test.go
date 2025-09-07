package cjincidentexpiration

import (
	"database/sql"
	"errors"
	"testing"
)

// mockRepository is a mock implementation of the Repository interface for testing.
type mockRepository struct {
	clustersToReturn      []ExpiredCluster
	votesToReturn         map[int64][]VoteRecord
	statsUpdateCalls      []UserStatsArgs
	processedClusters     []int64
	notifications         []NotificationArgs
	getExpiredClustersErr error
	getVotesErr           error
}

// UserStatsArgs captures the arguments for UpdateUserStats calls.
type UserStatsArgs struct {
	AccountID         int64
	ScoreChange       float64
	CredibilityChange float64
}

// NotificationArgs captures the arguments for CreateNotification calls.
type NotificationArgs struct {
	AccountID int64
	ClusterID int64
	Message   string
}

func (m *mockRepository) GetExpiredClusters() ([]ExpiredCluster, error) {
	return m.clustersToReturn, m.getExpiredClustersErr
}

func (m *mockRepository) GetVotesForCluster(clusterID int64) ([]VoteRecord, error) {
	if m.getVotesErr != nil {
		return nil, m.getVotesErr
	}
	votes, ok := m.votesToReturn[clusterID]
	if !ok {
		return nil, errors.New("no votes found for this cluster in mock")
	}
	return votes, nil
}

func (m *mockRepository) UpdateUserStats(accountID int64, scoreChange float64, credibilityChange float64) error {
	m.statsUpdateCalls = append(m.statsUpdateCalls, UserStatsArgs{accountID, scoreChange, credibilityChange})
	return nil
}

func (m *mockRepository) MarkClusterProcessed(clusterID int64) error {
	m.processedClusters = append(m.processedClusters, clusterID)
	return nil
}

func (m *mockRepository) CreateNotification(accountID int64, clusterID int64, message string) error {
	m.notifications = append(m.notifications, NotificationArgs{accountID, clusterID, message})
	return nil
}

func TestService_Run_HappyPath(t *testing.T) {
	// 1. Setup
	mockRepo := &mockRepository{
		clustersToReturn: []ExpiredCluster{
			{ID: 101, Credibility: sql.NullFloat64{Float64: 8.5, Valid: true}}, // This cluster will have a "True" outcome
		},
		votesToReturn: map[int64][]VoteRecord{
			101: {
				{AccountID: 1, Vote: true},  // This user voted correctly
				{AccountID: 2, Vote: false}, // This user voted incorrectly
			},
		},
	}

	service := NewService(mockRepo)

	// 2. Execute
	service.Run()

	// 3. Assert
	// Check if user stats were updated correctly
	if len(mockRepo.statsUpdateCalls) != 2 {
		t.Fatalf("expected 2 calls to UpdateUserStats, but got %d", len(mockRepo.statsUpdateCalls))
	}

	winnerCall := mockRepo.statsUpdateCalls[0]
	if winnerCall.AccountID != 1 || winnerCall.ScoreChange != scoreWin || winnerCall.CredibilityChange != credibilityWin {
		t.Errorf("incorrect stats for winner. Got: %+v, Want: AccountID=1, Score=%.1f, Credibility=%.1f", winnerCall, scoreWin, credibilityWin)
	}

	loserCall := mockRepo.statsUpdateCalls[1]
	if loserCall.AccountID != 2 || loserCall.ScoreChange != scoreLoss || loserCall.CredibilityChange != credibilityLoss {
		t.Errorf("incorrect stats for loser. Got: %+v, Want: AccountID=2, Score=%.1f, Credibility=%.1f", loserCall, scoreLoss, credibilityLoss)
	}

	// Check if notifications were created
	if len(mockRepo.notifications) != 2 {
		t.Fatalf("expected 2 notifications to be created, but got %d", len(mockRepo.notifications))
	}

	// Check if cluster was marked as processed
	if len(mockRepo.processedClusters) != 1 || mockRepo.processedClusters[0] != 101 {
		t.Errorf("expected cluster 101 to be marked as processed, but got %v", mockRepo.processedClusters)
	}
}

func TestService_Run_NoClusters(t *testing.T) {
	// 1. Setup
	mockRepo := &mockRepository{
		clustersToReturn: []ExpiredCluster{}, // No clusters to process
	}

	service := NewService(mockRepo)

	// 2. Execute
	service.Run()

	// 3. Assert
	if len(mockRepo.statsUpdateCalls) != 0 {
		t.Errorf("expected 0 calls to UpdateUserStats, but got %d", len(mockRepo.statsUpdateCalls))
	}
	if len(mockRepo.processedClusters) != 0 {
		t.Errorf("expected 0 clusters to be processed, but got %d", len(mockRepo.processedClusters))
	}
}
