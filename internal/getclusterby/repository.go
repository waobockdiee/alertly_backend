package getclusterby

import (
	"alertly/internal/common"
	"alertly/internal/dbtypes"
	"database/sql"
	"fmt"
	"log"
	"time"
)

type Repository interface {
	GetIncidentBy(inclId int64) (Cluster, error)
	GetIncidentByPublic(inclId int64) (Cluster, error)
	GetAccountAlreadyVoted(inclID, AccountID int64) (bool, error)
	GetAccountAlreadySaved(inclID, AccountID int64) (bool, error)
	GetUserVote(inclID, AccountID int64) (int, error)
	SaveAccountHistory(accountID, inclID int64) error
}

type pgRepository struct {
	db *sql.DB
}

func NewRepository(db *sql.DB) Repository {
	return &pgRepository{db: db}
}

func (r *pgRepository) GetIncidentBy(inclId int64) (Cluster, error) {
	return r.getIncidentByWithActiveFilter(inclId, true)
}

// GetIncidentByPublic obtiene un incidente sin filtrar por is_active (para endpoint público)
func (r *pgRepository) GetIncidentByPublic(inclId int64) (Cluster, error) {
	return r.getIncidentByWithActiveFilter(inclId, false)
}

// getIncidentByWithActiveFilter es el método base que permite filtrar opcionalmente por is_active
func (r *pgRepository) getIncidentByWithActiveFilter(inclId int64, activeOnly bool) (Cluster, error) {
	// ✅ SEGURIDAD: Primero verificar si existe cluster o solo incident_report
	// 1. Query principal del cluster (más eficiente)
	activeFilter := ""
	if activeOnly {
		activeFilter = "AND c.is_active = '1'"
	}

	clusterQuery := fmt.Sprintf(`
    SELECT
    c.incl_id,
    COALESCE(c.address, ''),
    c.center_latitude,
    c.center_longitude,
    COALESCE(c.city, ''),
    c.counter_total_comments,
    c.counter_total_flags,
    c.counter_total_views,
    c.counter_total_votes,
    c.counter_total_votes_true,
    c.counter_total_votes_false,
    c.created_at,
    COALESCE(c.description, ''),
    c.end_time,
    COALESCE(c.event_type, ''),
    c.incident_count,
    c.is_active,
    c.insu_id,
    COALESCE(c.media_type, ''),
    COALESCE(c.media_url, ''),
    COALESCE(c.postal_code, ''),
    COALESCE(c.province, ''),
    c.start_time,
    COALESCE(c.subcategory_name, ''),
    COALESCE(c.category_code, ''),
    COALESCE(c.subcategory_code, ''),
    c.credibility,
    c.account_id
  FROM incident_clusters c
  WHERE c.incl_id = $1 %s;
  `, activeFilter)

	var cluster Cluster
	var isActive dbtypes.NullBool
	err := r.db.QueryRow(clusterQuery, inclId).Scan(
		&cluster.InclId, &cluster.Address, &cluster.CenterLatitude, &cluster.CenterLongitude,
		&cluster.City, &cluster.CounterTotalComments, &cluster.CounterTotalFlags, &cluster.CounterTotalViews,
		&cluster.CounterTotalVotes, &cluster.CounterTotalVotesTrue, &cluster.CounterTotalVotesFalse,
		&cluster.CreatedAt, &cluster.Description, &cluster.EndTime, &cluster.EventType, &cluster.IncidentCount,
		&isActive, &cluster.InsuId, &cluster.MediaType, &cluster.MediaUrl, &cluster.PostalCode,
		&cluster.Province, &cluster.StartTime, &cluster.SubcategoryName, &cluster.CategoryCode,
		&cluster.SubcategoryCode, &cluster.Credibility, &cluster.AccountId,
	)
	cluster.IsActive = isActive.Valid && isActive.Bool

	if err != nil {
		// ✅ FALLBACK: Si no hay cluster, intentar crear uno temporal desde incident_report
		if err == sql.ErrNoRows {
			log.Printf("No cluster found for incl_id %d, attempting fallback to individual incident", inclId)
			return r.createClusterFromIndividualIncident(inclId, activeOnly)
		}
		return cluster, fmt.Errorf("error scanning cluster: %w", err)
	}

	// 2. Query separada para incidentes (más eficiente)
	incidentsQuery := `
        SELECT
            r.inre_id,
            COALESCE(r.media_url, ''),
            COALESCE(r.description, ''),
            COALESCE(r.event_type, ''),
            COALESCE(r.is_anonymous, '0'),
            COALESCE(r.subcategory_name, ''),
            a.account_id,
            CASE WHEN TRIM(COALESCE(r.is_anonymous, '0')) = '1' THEN '' ELSE COALESCE(a.nickname, '') END as nickname,
            CASE WHEN TRIM(COALESCE(r.is_anonymous, '0')) = '1' THEN '' ELSE COALESCE(a.first_name, '') END as first_name,
            CASE WHEN TRIM(COALESCE(r.is_anonymous, '0')) = '1' THEN '' ELSE COALESCE(a.last_name, '') END as last_name,
            COALESCE(a.is_private_profile, 0) as is_private_profile,
            CASE WHEN TRIM(COALESCE(r.is_anonymous, '0')) = '1' THEN '' ELSE COALESCE(a.thumbnail_url, '') END as thumbnail_url,
            CASE WHEN TRIM(COALESCE(r.is_anonymous, '0')) = '1' THEN 0 ELSE COALESCE(a.score, 0) END as score,
            r.created_at,
            r.incl_id,
            COALESCE(r.status, '')
        FROM incident_reports r
        INNER JOIN account a ON r.account_id = a.account_id
        WHERE r.incl_id = $1 AND COALESCE(r.is_active, '0') = '1'
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
	query := `UPDATE incident_clusters SET counter_total_views = counter_total_views + 1 WHERE incl_id = $1`
	_, err = r.db.Exec(query, inclId)

	if err != nil {
		log.Printf("Error updating view count for cluster %d: %v", inclId, err)
	}

	return cluster, nil
}

// para saber si el usuario voto. Es necesario saber que la forma de verificar es que al un usuario votar. Basicamente esta creando un incidente nuevo. Y este se asocia al cluster.
func (r *pgRepository) GetAccountAlreadyVoted(inclID, AccountID int64) (bool, error) {
	query := `SELECT COUNT(*) AS total FROM incident_reports WHERE incl_id = $1 AND account_id = $2 AND vote IS NOT NULL`
	var total int
	err := r.db.QueryRow(query, inclID, AccountID).Scan(&total)
	if err != nil {
		return false, err
	}
	return total > 0, nil
}

func (r *pgRepository) GetAccountAlreadySaved(inclID, AccountID int64) (bool, error) {
	query := `SELECT COUNT(*) AS total FROM account_cluster_saved WHERE incl_id = $1 AND account_id = $2`
	var total int
	err := r.db.QueryRow(query, inclID, AccountID).Scan(&total)
	if err != nil {
		return false, err
	}
	return total > 0, nil
}

func (r *pgRepository) GetUserVote(inclID, AccountID int64) (int, error) {
	query := `SELECT vote FROM incident_reports WHERE incl_id = $1 AND account_id = $2 AND vote IS NOT NULL LIMIT 1`
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

func (r *pgRepository) SaveAccountHistory(accountID, inclID int64) error {
	query := `INSERT INTO account_history(account_id, incl_id) VALUES($1, $2)`
	_, err := r.db.Exec(query, accountID, inclID)
	return err
}

// ✅ FALLBACK: Crear cluster temporal desde incident_report individual
func (r *pgRepository) createClusterFromIndividualIncident(inclId int64, activeOnly bool) (Cluster, error) {
	log.Printf("Creating temporary cluster from individual incident %d (activeOnly: %v)", inclId, activeOnly)

	// Query para obtener datos del incident_report individual
	statusFilter := ""
	if activeOnly {
		statusFilter = "AND r.status = 'active'"
	}

	individualQuery := fmt.Sprintf(`
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
		WHERE r.incl_id = $1 %s
		LIMIT 1
	`, statusFilter)

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
	incidentsQuery := fmt.Sprintf(`
        SELECT
            r.inre_id,
            COALESCE(r.media_url, ''),
            COALESCE(r.description, ''),
            COALESCE(r.event_type, ''),
            COALESCE(r.is_anonymous, '0'),
            COALESCE(r.subcategory_name, ''),
            COALESCE(a.account_id, 0),
            CASE WHEN TRIM(COALESCE(r.is_anonymous, '0')) = '1' THEN '' ELSE COALESCE(a.nickname, '') END as nickname,
            CASE WHEN TRIM(COALESCE(r.is_anonymous, '0')) = '1' THEN '' ELSE COALESCE(a.first_name, '') END as first_name,
            CASE WHEN TRIM(COALESCE(r.is_anonymous, '0')) = '1' THEN '' ELSE COALESCE(a.last_name, '') END as last_name,
            COALESCE(a.is_private_profile, 0) as is_private_profile,
            CASE WHEN TRIM(COALESCE(r.is_anonymous, '0')) = '1' THEN '' ELSE COALESCE(a.thumbnail_url, '') END as thumbnail_url,
            CASE WHEN TRIM(COALESCE(r.is_anonymous, '0')) = '1' THEN 0 ELSE COALESCE(a.score, 0) END as score,
            r.created_at,
            r.incl_id,
            COALESCE(r.status, '')
        FROM incident_reports r
        LEFT JOIN account a ON r.account_id = a.account_id
        WHERE r.incl_id = $1 %s
        ORDER BY r.created_at DESC
    `, statusFilter)

	rows, err := r.db.Query(incidentsQuery, inclId)
	if err != nil {
		return cluster, fmt.Errorf("error querying incidents for fallback cluster: %w", err)
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
			return cluster, fmt.Errorf("error scanning incident for fallback: %w", err)
		}

		// ✅ Convertir sql.NullTime a common.CustomTime
		if createdAt.Valid {
			incident.CreatedAt = common.CustomTime{Time: createdAt.Time}
		} else {
			incident.CreatedAt = common.CustomTime{Time: time.Time{}}
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
