package getclusterby

import (
	"alertly/internal/common"
	"database/sql"
	"fmt"
	"log"
	"time"
)

type Repository interface {
	GetIncidentBy(inclId int64) (Cluster, error)
	GetAccountAlreadyVoted(inclID, AccountID int64) (bool, error)
	GetAccountAlreadySaved(inclID, AccountID int64) (bool, error)
	GetUserVote(inclID, AccountID int64) (int, error)
	SaveAccountHistory(accountID, inclID int64) error
}

type mysqlRepository struct {
	db *sql.DB
}

func NewRepository(db *sql.DB) Repository {
	return &mysqlRepository{db: db}
}

func (r *mysqlRepository) GetIncidentBy(inclId int64) (Cluster, error) {
	// ✅ SEGURIDAD: Primero verificar si existe cluster o solo incident_report
	// 1. Query principal del cluster (más eficiente)
	clusterQuery := `
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
    c.account_id
  FROM incident_clusters c
  WHERE c.incl_id = ? AND c.is_active = 1;
  `

	var cluster Cluster
	err := r.db.QueryRow(clusterQuery, inclId).Scan(
		&cluster.InclId, &cluster.Address, &cluster.CenterLatitude, &cluster.CenterLongitude,
		&cluster.City, &cluster.CounterTotalComments, &cluster.CounterTotalFlags, &cluster.CounterTotalViews,
		&cluster.CounterTotalVotes, &cluster.CounterTotalVotesTrue, &cluster.CounterTotalVotesFalse,
		&cluster.CreatedAt, &cluster.Description, &cluster.EndTime, &cluster.EventType, &cluster.IncidentCount,
		&cluster.IsActive, &cluster.InsuId, &cluster.MediaType, &cluster.MediaUrl, &cluster.PostalCode,
		&cluster.Province, &cluster.StartTime, &cluster.SubcategoryName, &cluster.CategoryCode,
		&cluster.SubcategoryCode, &cluster.Credibility, &cluster.AccountId,
	)

	if err != nil {
		// ✅ FALLBACK: Si no hay cluster, intentar crear uno temporal desde incident_report
		if err == sql.ErrNoRows {
			log.Printf("No cluster found for incl_id %d, attempting fallback to individual incident", inclId)
			return r.createClusterFromIndividualIncident(inclId)
		}
		return cluster, fmt.Errorf("error scanning cluster: %w", err)
	}

	// 2. Query separada para incidentes (más eficiente)
	incidentsQuery := `
        SELECT 
            r.inre_id,
            r.media_url,
            r.description,
            r.event_type,
            r.is_anonymous,
            r.subcategory_name,
            a.account_id,
            IF(r.is_anonymous, '', a.nickname) as nickname,
            IF(r.is_anonymous, '', a.first_name) as first_name,
            IF(r.is_anonymous, '', a.last_name) as last_name,
            a.is_private_profile,
            IF(r.is_anonymous, '', COALESCE(a.thumbnail_url, '')) as thumbnail_url,
            IF(r.is_anonymous, 0, COALESCE(a.score, 0)) as score,
            r.created_at,
            r.incl_id,
            r.status
        FROM incident_reports r
        INNER JOIN account a ON r.account_id = a.account_id
        WHERE r.incl_id = ? AND r.is_active = 1
        ORDER BY r.created_at DESC
        LIMIT 50
    `

	rows, err := r.db.Query(incidentsQuery, inclId)
	if err != nil {
		return cluster, fmt.Errorf("error querying incidents: %w", err)
	}
	defer rows.Close()

	var incidents []Incident
	for rows.Next() {
		var incident Incident
		var createdAt sql.NullTime
		err := rows.Scan(
			&incident.InreId, &incident.MediaUrl, &incident.Description, &incident.EventType,
			&incident.IsAnonymous, &incident.SubcategortyName, &incident.AccountId, &incident.Nickname,
			&incident.FirstName, &incident.LastName, &incident.IsPrivateProfile, &incident.ThumbnailUrl,
			&incident.Score, &createdAt, &incident.InclID, &incident.Status,
		)
		if err != nil {
			return cluster, fmt.Errorf("error scanning incident: %w", err)
		}

		// ✅ Convertir sql.NullTime a common.CustomTime
		if createdAt.Valid {
			incident.CreatedAt = common.CustomTime{Time: createdAt.Time}
		} else {
			incident.CreatedAt = common.CustomTime{Time: time.Time{}}
		}

		incidents = append(incidents, incident)
	}

	if err := rows.Err(); err != nil {
		return cluster, fmt.Errorf("error iterating incidents: %w", err)
	}

	cluster.Incidents = incidents

	// ✅ Actualizar contador de vistas
	query := `UPDATE incident_clusters SET counter_total_views = counter_total_views + 1 WHERE incl_id = ?`
	_, err = r.db.Exec(query, inclId)

	if err != nil {
		log.Printf("Error updating view count for cluster %d: %v", inclId, err)
	}

	return cluster, nil
}

// para saber si el usuario voto. Es necesario saber que la forma de verificar es que al un usuario votar. Basicamente esta creando un incidente nuevo. Y este se asocia al cluster.
func (r *mysqlRepository) GetAccountAlreadyVoted(inclID, AccountID int64) (bool, error) {
	query := `SELECT COUNT(*) AS total FROM incident_reports WHERE incl_id = ? AND account_id = ? AND vote IS NOT NULL`
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

func (r *mysqlRepository) GetUserVote(inclID, AccountID int64) (int, error) {
	query := `SELECT vote FROM incident_reports WHERE incl_id = ? AND account_id = ? AND vote IS NOT NULL LIMIT 1`
	var vote sql.NullInt64
	err := r.db.QueryRow(query, inclID, AccountID).Scan(&vote)
	if err != nil {
		if err == sql.ErrNoRows {
			return -1, nil // -1 significa que no ha votado
		}
		return -1, err
	}
	if !vote.Valid {
		return -1, nil // No hay voto válido
	}
	// Convertir el valor de la base de datos a nuestro formato
	// 1 = TRUE, 0 = FALSE
	return int(vote.Int64), nil
}

func (r *mysqlRepository) SaveAccountHistory(accountID, inclID int64) error {
	query := `INSERT INTO account_history(account_id, incl_id) VALUES(?, ?)`
	_, err := r.db.Exec(query, accountID, inclID)
	return err
}

// ✅ FALLBACK: Crear cluster temporal desde incident_report individual
func (r *mysqlRepository) createClusterFromIndividualIncident(inclId int64) (Cluster, error) {
	log.Printf("Creating temporary cluster from individual incident %d", inclId)

	// Query para obtener datos del incident_report individual
	individualQuery := `
		SELECT 
			r.incl_id,
			COALESCE(r.address, 'Unknown Location') as address,
			COALESCE(r.latitude, 0) as latitude,
			COALESCE(r.longitude, 0) as longitude,
			COALESCE(r.city, '') as city,
			r.description,
			r.event_type,
			r.subcategory_name,
			r.created_at,
			r.media_url,
			r.media_type,
			sc.category_code,
			r.subcategory_code,
			r.insu_id
		FROM incident_reports r
		LEFT JOIN subcategories sc ON r.subcategory_code = sc.subcategory_code
		WHERE r.incl_id = ? AND r.status = 'active'
		LIMIT 1
	`

	var cluster Cluster
	err := r.db.QueryRow(individualQuery, inclId).Scan(
		&cluster.InclId,
		&cluster.Address,
		&cluster.CenterLatitude,
		&cluster.CenterLongitude,
		&cluster.City,
		&cluster.Description,
		&cluster.EventType,
		&cluster.SubcategoryName,
		&cluster.CreatedAt,
		&cluster.MediaUrl,
		&cluster.MediaType,
		&cluster.CategoryCode,
		&cluster.SubcategoryCode,
		&cluster.InsuId,
	)

	if err != nil {
		return cluster, fmt.Errorf("error creating cluster from individual incident: %w", err)
	}

	// ✅ Valores por defecto para cluster temporal
	cluster.IncidentCount = 1
	cluster.IsActive = true // ✅ FIX: bool, no int
	cluster.CounterTotalVotes = 0
	cluster.CounterTotalVotesTrue = 0
	cluster.CounterTotalVotesFalse = 0
	cluster.CounterTotalComments = 0
	cluster.CounterTotalFlags = 0
	cluster.CounterTotalViews = 0
	cluster.Credibility = 0.0
	cluster.StartTime = cluster.CreatedAt
	cluster.EndTime = cluster.CreatedAt

	log.Printf("✅ Temporary cluster created for incident %d", inclId)

	// Ahora obtener los incident_reports (será solo uno en este caso)
	incidentsQuery := `
        SELECT 
            r.inre_id,
            r.media_url,
            r.description,
            r.event_type,
            r.is_anonymous,
            r.subcategory_name,
            a.account_id,
            IF(r.is_anonymous, '', a.nickname) as nickname,
            IF(r.is_anonymous, '', a.first_name) as first_name,
            IF(r.is_anonymous, '', a.last_name) as last_name,
            a.is_private_profile,
            IF(r.is_anonymous, '', COALESCE(a.thumbnail_url, '')) as thumbnail_url,
            IF(r.is_anonymous, 0, COALESCE(a.score, 0)) as score,
            r.created_at,
            r.incl_id,
            r.status
        FROM incident_reports r
        LEFT JOIN account a ON r.account_id = a.account_id
        WHERE r.incl_id = ? AND r.status = 'active'
        ORDER BY r.created_at DESC
    `

	rows, err := r.db.Query(incidentsQuery, inclId)
	if err != nil {
		return cluster, fmt.Errorf("error querying incidents for fallback cluster: %w", err)
	}
	defer rows.Close()

	var incidents []Incident
	for rows.Next() {
		var incident Incident

		err := rows.Scan(
			&incident.InreId, &incident.MediaUrl, &incident.Description, &incident.EventType,
			&incident.IsAnonymous, &incident.SubcategortyName, &incident.AccountId, &incident.Nickname,
			&incident.FirstName, &incident.LastName, &incident.IsPrivateProfile, &incident.ThumbnailUrl,
			&incident.Score, &incident.CreatedAt, &incident.InclID, &incident.Status,
		)
		if err != nil {
			return cluster, fmt.Errorf("error scanning incident for fallback: %w", err)
		}

		// ✅ FIX: Calcular time_diff manualmente (función común no existe)
		incident.TimeDiff = calculateTimeDifference(incident.CreatedAt.Time)
		incidents = append(incidents, incident)
	}

	if err = rows.Err(); err != nil {
		return cluster, fmt.Errorf("error iterating incidents for fallback: %w", err)
	}

	cluster.Incidents = incidents
	log.Printf("✅ Fallback cluster completed with %d incidents", len(incidents))

	return cluster, nil
}

// ✅ Helper function para calcular diferencia de tiempo
func calculateTimeDifference(createdAt time.Time) string {
	now := time.Now()
	diff := now.Sub(createdAt)

	if diff < time.Minute {
		return "just now"
	} else if diff < time.Hour {
		minutes := int(diff.Minutes())
		if minutes == 1 {
			return "1 minute ago"
		}
		return fmt.Sprintf("%d minutes ago", minutes)
	} else if diff < 24*time.Hour {
		hours := int(diff.Hours())
		if hours == 1 {
			return "1 hour ago"
		}
		return fmt.Sprintf("%d hours ago", hours)
	} else {
		days := int(diff.Hours() / 24)
		if days == 1 {
			return "1 day ago"
		}
		return fmt.Sprintf("%d days ago", days)
	}
}
