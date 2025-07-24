package getincidentsasreels

import (
	"alertly/internal/getclusterby"
	"database/sql"
	"encoding/json"
	"fmt"
	"strings"
)

const maxDistanceMeters = 500

type Repository interface {
	GetReel(inputs Inputs, accountID int64) ([]getclusterby.Cluster, error)
}

type mysqlRepository struct {
	db *sql.DB
}

func NewRepository(db *sql.DB) Repository {
	return &mysqlRepository{db: db}
}

func (r *mysqlRepository) GetReel(inputs Inputs, accountID int64) ([]getclusterby.Cluster, error) {

	idQuery := `
    SELECT c.incl_id
    FROM incident_clusters c
    WHERE
        (c.center_latitude  BETWEEN ? AND ? AND c.center_longitude BETWEEN ? AND ?)
        OR EXISTS (
          SELECT 1
          FROM account_favorite_locations f
          WHERE f.account_id = ?
            AND ST_Distance_Sphere(
                  POINT(c.center_longitude, c.center_latitude),
                  POINT(f.longitude, f.latitude)
                ) <= ?
        )
	AND c.is_active = 1
    ORDER BY RAND()
    LIMIT 20
    `
	rows, err := r.db.Query(idQuery,
		inputs.MinLatitude, inputs.MaxLatitude, inputs.MinLongitude, inputs.MaxLongitude,
		accountID, maxDistanceMeters,
	)
	if err != nil {
		return nil, fmt.Errorf("fetch random ids: %w", err)
	}
	defer rows.Close()

	var ids []int64
	for rows.Next() {
		var id int64
		if err := rows.Scan(&id); err != nil {
			return nil, fmt.Errorf("scan id: %w", err)
		}
		ids = append(ids, id)
	}
	if len(ids) == 0 {
		return nil, nil
	}

	// Paso 2: reconstruir el detalle solo para esos 20 IDs
	placeholders := make([]string, len(ids))
	args := make([]interface{}, len(ids))
	for i, id := range ids {
		placeholders[i] = "?"
		args[i] = id
	}
	inClause := strings.Join(placeholders, ",")

	detailQuery := fmt.Sprintf(`
    SELECT
      c.incl_id,
      c.created_at,
      c.start_time,
      c.end_time,
      c.insu_id,
      c.media_url,
      c.center_latitude,
      c.center_longitude,
      c.is_active,
      c.media_type,
      c.event_type,
      c.description,
      c.address,
      c.city,
      c.province,
      c.postal_code,
      c.subcategory_name,
      c.category_code,
      c.subcategory_code,
      c.incident_count,
      c.counter_total_comments,
      c.counter_total_votes,
      c.counter_total_views,
      c.counter_total_flags,
      c.counter_total_votes_true,
      c.counter_total_votes_false,
      COALESCE(c.credibility,0) AS credibility,
      IFNULL(
        (
          SELECT JSON_ARRAYAGG(
            JSON_OBJECT(
              'inre_id',        r.inre_id,
              'media_url',      r.media_url,
              'description',    r.description,
              'event_type',     r.event_type,
              'is_anonymous',   r.is_anonymous,
              'subcategory_name', r.subcategory_name,
              'account_id',     a.account_id,
              'nickname',       IF(r.is_anonymous, '', a.nickname),
              'first_name',     IF(r.is_anonymous, '', a.first_name),
              'last_name',      IF(r.is_anonymous, '', a.last_name),
              'is_private_profile', a.is_private_profile,
              'thumbnail_url',  IF(r.is_anonymous, '', a.thumbnail_url),
              'created_at',     r.created_at
            )
          )
          FROM incident_reports r
          INNER JOIN account a ON r.account_id = a.account_id
          WHERE r.incl_id = c.incl_id
        ),
        JSON_ARRAY()
      ) AS incidents_json
    FROM incident_clusters c
    WHERE c.incl_id IN (%s)
    `, inClause)

	rows2, err := r.db.Query(detailQuery, args...)
	if err != nil {
		return nil, fmt.Errorf("fetch details: %w", err)
	}
	defer rows2.Close()

	var results []getclusterby.Cluster
	for rows2.Next() {
		var cl getclusterby.Cluster
		var rawIncidents sql.NullString

		// Escanea en los campos de tu Cluster y en rawIncidents
		if err := rows2.Scan(
			&cl.InclId,
			&cl.CreatedAt,
			&cl.StartTime,
			&cl.EndTime,
			&cl.InsuId,
			&cl.MediaUrl,
			&cl.CenterLatitude,
			&cl.CenterLongitude,
			&cl.IsActive,
			&cl.MediaType,
			&cl.EventType,
			&cl.Description,
			&cl.Address,
			&cl.City,
			&cl.Province,
			&cl.PostalCode,
			&cl.SubcategoryName,
			&cl.CategoryCode,
			&cl.SubcategoryCode,
			&cl.IncidentCount,
			&cl.CounterTotalComments,
			&cl.CounterTotalVotes,
			&cl.CounterTotalViews,
			&cl.CounterTotalFlags,
			&cl.CounterTotalVotesTrue,
			&cl.CounterTotalVotesFalse,
			&cl.Credibility,
			&rawIncidents,
		); err != nil {
			return nil, fmt.Errorf("scan detail row: %w", err)
		}

		// Deserializa el JSON de incidentes
		cl.Incidents = make([]getclusterby.Incident, 0)
		if rawIncidents.Valid {
			if err := json.Unmarshal([]byte(rawIncidents.String), &cl.Incidents); err != nil {
				return nil, fmt.Errorf("unmarshal incidents JSON: %w", err)
			}
		}

		results = append(results, cl)
	}

	return results, nil
}
