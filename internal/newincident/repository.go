package newincident

import (
	"alertly/internal/common"
	"database/sql"
	"fmt"
)

type Repository interface {
	CheckAndGetIfClusterExist(incident IncidentReport) (Cluster, error)
	Save(incident IncidentReport) (int64, error)
	SaveCluster(cluster Cluster, accountId int64) (int64, error)
	UpdateCluster(inclId int64, incident IncidentReport) (sql.Result, error)
}

type mysqlRepository struct {
	db *sql.DB
}

func NewRepository(db *sql.DB) Repository {
	return &mysqlRepository{db: db}
}

func (r *mysqlRepository) CheckAndGetIfClusterExist(incident IncidentReport) (Cluster, error) {
	query := `SELECT incl_id FROM incident_clusters WHERE insu_id = ? 
	  AND ST_Distance_Sphere(
		POINT(center_longitude, center_latitude),
		POINT(?, ?)
	  ) < ?
	  AND created_at >= DATE_SUB(NOW(), INTERVAL 24 HOUR);`

	row := r.db.QueryRow(query, incident.InsuId, incident.Longitude, incident.Latitude, incident.DefaultCircleRange)
	var cluster Cluster
	err := row.Scan(&cluster.InclId)
	return cluster, err
}

// -- guarda un incidente del cluster. Basicamente es una actualizacion del seguimiento del cluster
func (r *mysqlRepository) Save(incident IncidentReport) (int64, error) {

	tx, err := r.db.Begin()

	if err != nil {
		return 0, fmt.Errorf("failed to begin transaction: %w", err)
	}

	defer func() {
		if err != nil {
			_ = tx.Rollback()
		} else {
			_ = tx.Commit()
		}
	}()

	query := "INSERT INTO incident_reports(account_id, insu_id, incl_id, description, event_type, address, city, province, postal_code, latitude, longitude, subcategory_name, is_anonymous, media_url, subcategory_code, category_code, created_at) VALUES( ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, NOW())"
	result, err := tx.Exec(query,
		incident.AccountId,
		incident.InsuId,
		incident.InclId,
		incident.Description,
		incident.EventType,
		incident.Address,
		incident.City,
		incident.Province,
		incident.PostalCode,
		incident.Latitude,
		incident.Longitude,
		incident.SubCategoryName,
		incident.IsAnonymous,
		incident.Media.Uri,
		incident.SubcategoryCode,
		incident.CategoryCode,
	)

	if err != nil {
		return 0, fmt.Errorf("failed to insert comment: %w", err)
	}

	id, err := result.LastInsertId()
	if err != nil {
		return 0, fmt.Errorf("failed to get last insert id: %w", err)
	}

	// ----------------------CITIZEN SCORE----------------------- //
	// -- Save to DB
	err = common.SaveScore(tx, incident.AccountId, 20)
	if err != nil {
		fmt.Println("error saving score") // It's not necesary to stop the server
	}
	// ----------------------NOTIFICATION----------------------- //
	// -- Save to DB
	err = common.SaveNotification(tx, "new_cluster", incident.AccountId, id)

	if err != nil {
		fmt.Println("error saving notification on createincident event")
	}

	return id, nil
}

func (r *mysqlRepository) SaveCluster(cluster Cluster, accountID int64) (int64, error) {
	tx, err := r.db.Begin()
	if err != nil {
		return 0, fmt.Errorf("failed to begin transaction: %w", err)
	}

	defer func() {
		if err != nil {
			_ = tx.Rollback()
		} else {
			_ = tx.Commit()
		}
	}()

	query := "INSERT INTO incident_clusters(created_at, start_time, end_time, media_url, center_latitude, center_longitude, insu_id, media_type, event_type, description, address, city, province, postal_code, subcategory_name, category_code, subcategory_code) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)"
	result, err := r.db.Exec(query,
		cluster.CreatedAt,
		cluster.StartTime,
		cluster.EndTime,
		cluster.MediaUrl,
		cluster.CenterLatitude,
		cluster.CenterLongitude,
		cluster.InsuId,
		cluster.MediaType,
		cluster.EventType,
		cluster.Description,
		cluster.Address,
		cluster.City,
		cluster.Province,
		cluster.PostalCode,
		cluster.SubcategoryName,
		cluster.CategoryCode,
		cluster.SubcategoryCode,
	)

	if err != nil {
		return 0, err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}

	// ----------------------CITIZEN SCORE----------------------- //
	// -- Save to DB
	err = common.SaveScore(tx, accountID, 20)
	if err != nil {
		fmt.Println("error saving score") // It's not necesary to stop the server
	}
	// ----------------------NOTIFICATION----------------------- //
	// -- Save to DB
	err = common.SaveNotification(tx, "new_cluster", accountID, id)

	if err != nil {
		fmt.Println("error saving notification on createincident event")
	}

	fmt.Printf("Usuario insertado con ID: %d\n", id)
	return id, nil
}

// -- Actualiza la localizacion del cluster cuando se crea un incidente nuevo del cluster. Ya que el incidente nuevo no necesariamente tiene que esta ubicada en las coordenadas exactas del cluster. Por eso actualiza.
func (r *mysqlRepository) UpdateCluster(inclId int64, incident IncidentReport) (sql.Result, error) {
	query := "UPDATE incident_clusters SET center_latitude = (center_latitude + ?) / 2, center_longitude = (center_longitude + ?) / 2, incident_count = incident_count+1 WHERE incl_id = ?"

	result, err := r.db.Exec(query, incident.Latitude, incident.Longitude, inclId)
	if err != nil {
		return result, fmt.Errorf("error updating cluster  %w", err)
	}
	return result, nil
}
