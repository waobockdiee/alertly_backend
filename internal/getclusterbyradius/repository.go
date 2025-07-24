package getclusterbyradius

import (
	"database/sql"
	"strings"
)

type Repository interface {
	GetClustersByRadius(inputs Inputs) ([]Cluster, error)
}

type mysqlRepository struct {
	db *sql.DB
}

func NewRepository(db *sql.DB) Repository {
	return &mysqlRepository{db: db}
}

func (r *mysqlRepository) GetClustersByRadius(inputs Inputs) ([]Cluster, error) {
	query := `
		SELECT
			t1.incl_id, t1.center_latitude, t1.center_longitude, t1.insu_id, t1.category_code, t1.subcategory_code
		FROM incident_clusters t1
		WHERE ST_Distance_Sphere(
			point(t1.center_longitude, t1.center_latitude),
			point(?, ?)
		) <= ?
		AND DATE(t1.start_time) <= ?
		AND DATE(t1.end_time) >= ?
		AND (? = 0 OR t1.insu_id = ?)
		AND tt1.is_active = 1
	`

	params := []interface{}{
		inputs.Longitude, inputs.Latitude, inputs.Radius,
		inputs.ToDate,
		inputs.FromDate,
		inputs.InsuID, inputs.InsuID,
	}

	if inputs.Categories != "" {
		cats := strings.Split(inputs.Categories, ",")
		placeholders := make([]string, len(cats))
		for i := range cats {
			placeholders[i] = "?"
			params = append(params, strings.TrimSpace(cats[i]))
		}
		query += " AND t1.category_code IN (" + strings.Join(placeholders, ",") + ")"
	}

	var clusters []Cluster
	rows, err := r.db.Query(query, params...)
	if err != nil {
		return clusters, err
	}
	defer rows.Close()

	for rows.Next() {
		var cluster Cluster
		if err := rows.Scan(&cluster.InclId, &cluster.Latitude, &cluster.Longitude, &cluster.InsuId, &cluster.CategoryCode, &cluster.SubcategoryCode); err != nil {
			return clusters, err
		}
		clusters = append(clusters, cluster)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return clusters, nil
}
