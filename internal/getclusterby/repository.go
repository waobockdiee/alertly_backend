package getclusterby

import (
	"database/sql"
	"encoding/json"
	"fmt"
)

type Repository interface {
	GetIncidentBy(inclId int64) (Cluster, error)
	GetAccountAlreadyVoted(inclID, AccountID int64) (bool, error)
	GetAccountAlreadySaved(inclID, AccountID int64) (bool, error)
	SaveAccountHistory(accountID, inclID int64) error
}

type mysqlRepository struct {
	db *sql.DB
}

func NewRepository(db *sql.DB) Repository {
	return &mysqlRepository{db: db}
}

func (r *mysqlRepository) GetIncidentBy(inclId int64) (Cluster, error) {
	query := `
  SELECT 
  c.incl_id,
  c.address,
  c.center_latitude,
  c.center_longitude,
  c.city,
  c.counter_total_comments,
  c.counter_total_flags,
  c.counter_total_views,
  c.counter_total_votes,
  c.counter_total_votes_true,
  c.counter_total_votes_false,
  c.created_at,
  c.description,
  c.end_time,
  c.event_type,
  c.incident_count,
  c.is_active,
  c.insu_id,
  c.media_type,
  c.media_url,
  c.postal_code,
  c.province,
  c.start_time,
  c.subcategory_name,
  c.category_code,
  c.subcategory_code,
  c.credibility,
  IFNULL(
    (
      SELECT JSON_ARRAYAGG(
        JSON_OBJECT(
          'inre_id', r.inre_id,
          'media_url', r.media_url,
          'description', r.description,
          'event_type', r.event_type,
          'is_anonymous', r.is_anonymous,
          'subcategory_name', r.subcategory_name,
          'account_id', a.account_id,
          'nickname', IF(r.is_anonymous, '', a.nickname),
          'first_name', IF(r.is_anonymous, '', a.first_name),
          'last_name', IF(r.is_anonymous, '', a.last_name),
          'is_private_profile', a.is_private_profile,
          'thumbnail_url', IF(r.is_anonymous, '', a.thumbnail_url),
          'created_at', r.created_at
        )
      )
      FROM incident_reports r
      INNER JOIN account a ON r.account_id = a.account_id
      WHERE r.incl_id = c.incl_id
    ),
    JSON_ARRAY()
  ) AS incidents
FROM incident_clusters c
WHERE c.incl_id = ?;
`

	var cluster Cluster
	var rawIncidents string

	err := r.db.QueryRow(query, inclId).Scan(&cluster.InclId, &cluster.Address, &cluster.CenterLatitude, &cluster.CenterLongitude, &cluster.City, &cluster.CounterTotalComments, &cluster.CounterTotalFlags, &cluster.CounterTotalViews, &cluster.CounterTotalVotes, &cluster.CounterTotalVotesTrue, &cluster.CounterTotalVotesFalse, &cluster.CreatedAt, &cluster.Description, &cluster.EndTime, &cluster.EventType, &cluster.IncidentCount, &cluster.IsActive, &cluster.InsuId, &cluster.MediaType, &cluster.MediaUrl, &cluster.PostalCode, &cluster.Province, &cluster.StartTime, &cluster.SubcategoryName, &cluster.CategoryCode, &cluster.SubcategoryCode, &cluster.Credibility, &rawIncidents)

	if err != nil {
		return cluster, fmt.Errorf("error scanning row: %w", err)
	}

	if err := json.Unmarshal([]byte(rawIncidents), &cluster.Incidents); err != nil {
		return cluster, fmt.Errorf("error unmarshalling incidents: %w", err)
	}

	// Deserializar los comentarios
	// if err := json.Unmarshal([]byte(rawComments), &cluster.Comments); err != nil {
	// 	return cluster, fmt.Errorf("error unmarshalling comments: %w", err)
	// }

	return cluster, nil
}

// para saber si el usuario voto. Es necesario saber que la forma de verificar es que al un usuario votar. Basicamente esta creando un incidente nuevo. Y este se asocia al cluster.
func (r *mysqlRepository) GetAccountAlreadyVoted(inclID, AccountID int64) (bool, error) {
	query := `SELECT COUNT(*) AS total FROM incident_reports WHERE incl_id = ? AND account_id = ?`
	var total int
	err := r.db.QueryRow(query, inclID, AccountID).Scan(&total)
	if err != nil {
		return false, err
	}
	return total > 0, nil
}

func (r *mysqlRepository) GetAccountAlreadySaved(inclID, AccountID int64) (bool, error) {
	query := `SELECT COUNT(*) AS total FROM account_cluster_saved WHERE incl_id = ? AND account_id = ?`
	var total int
	err := r.db.QueryRow(query, inclID, AccountID).Scan(&total)
	if err != nil {
		return false, err
	}
	return total > 0, nil
}

func (r *mysqlRepository) SaveAccountHistory(accountID, inclID int64) error {
	query := `INSERT INTO account_history(account_id, incl_id) VALUES(?, ?)`
	_, err := r.db.Exec(query, accountID, inclID)
	return err
}
