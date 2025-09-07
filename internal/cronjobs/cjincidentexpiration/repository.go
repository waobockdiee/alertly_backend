package cjincidentexpiration

import (
	"alertly/internal/common"
	"database/sql"
	"fmt"
)

// Repository defines the interface for database operations related to incident expiration.
type Repository interface {
	GetExpiredClusters() ([]ExpiredCluster, error)
	GetVotesForCluster(clusterID int64) ([]VoteRecord, error)
	UpdateUserStats(accountID int64, scoreChange float64, credibilityChange float64) error
	MarkClusterProcessed(clusterID int64) error
	SaveWinNotification(accountID int64, clusterID int64, message string) error
	SaveLossNotification(accountID int64, clusterID int64, message string) error
}

type mysqlRepository struct {
	db *sql.DB
}

// NewRepository creates a new instance of the repository.
func NewRepository(db *sql.DB) Repository {
	return &mysqlRepository{db: db}
}

// ExpiredCluster holds the necessary information for a cluster that has expired.
type ExpiredCluster struct {
	ID          int64
	Credibility sql.NullFloat64 // Can be NULL in the database
}

// VoteRecord holds information about a single user's vote on an incident.
type VoteRecord struct {
	AccountID int64
	Vote      bool // true for 'true', false for 'false'
}

// GetExpiredClusters fetches all active clusters that have passed their expiration time.
func (r *mysqlRepository) GetExpiredClusters() ([]ExpiredCluster, error) {
	query := `
		SELECT
			ic.incl_id,
			ic.credibility
		FROM
			incident_clusters AS ic
		JOIN
			incident_subcategories AS isu ON ic.insu_id = isu.insu_id
		WHERE
			ic.is_active = '1'
			AND NOW() >= TIMESTAMPADD(HOUR, isu.default_duration_hours, ic.created_at);
	`
	rows, err := r.db.Query(query)
	if err != nil {
		return nil, fmt.Errorf("failed to query expired clusters: %w", err)
	}
	defer rows.Close()

	var clusters []ExpiredCluster
	for rows.Next() {
		var cluster ExpiredCluster
		if err := rows.Scan(&cluster.ID, &cluster.Credibility); err != nil {
			return nil, fmt.Errorf("failed to scan expired cluster: %w", err)
		}
		clusters = append(clusters, cluster)
	}

	return clusters, nil
}

// GetVotesForCluster retrieves all votes for a given cluster ID.
func (r *mysqlRepository) GetVotesForCluster(clusterID int64) ([]VoteRecord, error) {
	query := `
		SELECT
			account_id,
			vote
		FROM
			incident_reports
		WHERE
			incl_id = ? AND vote IS NOT NULL;
	`
	rows, err := r.db.Query(query, clusterID)
	if err != nil {
		return nil, fmt.Errorf("failed to query votes for cluster %d: %w", clusterID, err)
	}
	defer rows.Close()

	var votes []VoteRecord
	for rows.Next() {
		var vote VoteRecord
		if err := rows.Scan(&vote.AccountID, &vote.Vote); err != nil {
			return nil, fmt.Errorf("failed to scan vote for cluster %d: %w", clusterID, err)
		}
		votes = append(votes, vote)
	}
	return votes, nil
}

// UpdateUserStats updates the score and credibility for a given user.
func (r *mysqlRepository) UpdateUserStats(accountID int64, scoreChange float64, credibilityChange float64) error {
	query := `
		UPDATE account
		SET
			score = score + ?,
			credibility = credibility + ?
		WHERE
			account_id = ?;
	`
	_, err := r.db.Exec(query, scoreChange, credibilityChange, accountID)
	if err != nil {
		return fmt.Errorf("failed to update stats for user %d: %w", accountID, err)
	}
	return nil
}

// MarkClusterProcessed marks a cluster as inactive.
func (r *mysqlRepository) MarkClusterProcessed(clusterID int64) error {
	query := `
		UPDATE incident_clusters
		SET is_active = '0'
		WHERE incl_id = ?;
	`
	_, err := r.db.Exec(query, clusterID)
	if err != nil {
		return fmt.Errorf("failed to mark cluster %d as processed: %w", clusterID, err)
	}
	return nil
}

// SaveWinNotification saves a win notification for a user.
func (r *mysqlRepository) SaveWinNotification(accountID int64, clusterID int64, message string) error {
	return common.SaveNotification(r.db, "incident_result_win", accountID, clusterID, "¡Incidente resuelto!", message)
}

// SaveLossNotification saves a loss notification for a user.
func (r *mysqlRepository) SaveLossNotification(accountID int64, clusterID int64, message string) error {
	return common.SaveNotification(r.db, "incident_result_loss", accountID, clusterID, "¡Incidente resuelto!", message)
}
