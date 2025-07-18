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
	UpdateClusterAsTrue(inclId int64, accountID int64, latitude, longitude float32) (sql.Result, error)
	UpdateClusterAsFalse(inclId int64, accountID int64, latitude, longitude float32) (sql.Result, error)
	// SaveAsUpdate(incident IncidentReport) error
	HasAccountVoted(inclID, accountID int64) (bool, bool, error)
	UpdateClusterLocation(inclId int64, latitude, longitude float32) (sql.Result, error)
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

/*
Guarda un incidente del cluster. Basicamente es una actualizacion del seguimiento del cluster de una persona que ya ha votado o haya creado el cluster.
*/
func (r *mysqlRepository) Save(incident IncidentReport) (int64, error) {

	query := "INSERT INTO incident_reports(account_id, insu_id, incl_id, description, event_type, address, city, province, postal_code, latitude, longitude, subcategory_name, is_anonymous, media_url, subcategory_code, category_code, created_at) VALUES( ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, NOW())"
	result, err := r.db.Exec(query,
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
		return 0, fmt.Errorf("failed to insert incident report: %w", err)
	}

	id, err := result.LastInsertId()
	if err != nil {
		return 0, fmt.Errorf("failed to get last insert id: %w", err)
	}

	// ----------------------CITIZEN SCORE----------------------- //
	// -- Save to DB
	err = common.SaveScore(r.db, incident.AccountId, 20)
	if err != nil {
		fmt.Println("error saving score") // It's not necesary to stop the server
	}
	// ----------------------NOTIFICATION----------------------- //
	// -- Save to DB
	err = common.SaveNotification(r.db, "new_cluster", incident.AccountId, id)

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

	query := `
	INSERT INTO incident_clusters (
		created_at,
		start_time,
		end_time,
		media_url,
		center_latitude,
		center_longitude,
		insu_id,
		media_type,
		event_type,
		description,
		address,
		city,
		province,
		postal_code,
		subcategory_name,
		category_code,
		subcategory_code,
		score_true,
		score_false,
		credibility
	  )
	  VALUES (
		?,  ?,  ?,  ?,  ?,  ?,  ?,  ?,  ?,  ?,  ?,  ?,  ?,  ?,  ?,  ?,  ?,
		(SELECT a.credibility      FROM account a WHERE a.account_id = ?),
		(10 - (SELECT a.credibility FROM account a WHERE a.account_id = ?)),
		(SELECT a.credibility      FROM account a WHERE a.account_id = ?)
	  );`
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
		accountID,
		accountID,
		accountID,
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

	fmt.Printf("incidente creado con ID: %d\n", id)
	return id, nil
}

// -- Actualiza la localizacion del cluster cuando se crea un incidente nuevo del cluster. Ya que el incidente nuevo no necesariamente tiene que esta ubicada en las coordenadas exactas del cluster. Por eso actualiza.
func (r *mysqlRepository) UpdateClusterAsTrue(inclId int64, accountID int64, latitude, longitude float32) (sql.Result, error) {
	query := `
	UPDATE incident_clusters ic
	JOIN account a ON a.account_id = ?
	SET
	ic.center_latitude      = (ic.center_latitude + ?) / 2,
	ic.center_longitude     = (ic.center_longitude + ?) / 2,
	ic.counter_total_votes  = ic.counter_total_votes + 1,
	ic.score_true           = ic.score_true + a.credibility,
	ic.score_false          = ic.score_false + (10 - a.credibility),
	ic.credibility          = ic.score_true
								/ GREATEST(ic.score_true + ic.score_false, 1)
								* 10
	WHERE ic.incl_id = ?;
	`

	result, err := r.db.Exec(query, accountID, latitude, longitude, inclId)
	return result, err
}

func (r *mysqlRepository) UpdateClusterAsFalse(inclId int64, accountID int64, latitude, longitude float32) (sql.Result, error) {
	query := `
	UPDATE incident_clusters ic
	JOIN account a ON a.account_id = ?
	SET
	ic.center_latitude      = (ic.center_latitude + ?) / 2,
	ic.center_longitude     = (ic.center_longitude + ?) / 2,
	ic.counter_total_votes  = ic.counter_total_votes + 1,
	ic.score_true           = ic.score_true + (10 - a.credibility),
	ic.score_false          = ic.score_false + a.credibility,
	ic.credibility          = ic.score_true
								/ GREATEST(ic.score_true + ic.score_false, 1)
								* 10
	WHERE ic.incl_id = ?;
	`

	result, err := r.db.Exec(query, accountID, latitude, longitude, inclId)
	return result, err
}

// func (r *mysqlRepository) SaveAsUpdate(incident IncidentReport) error {
// 	tx, err := r.db.Begin()

// 	if err != nil {
// 		fmt.Printf("error repository SaveAsUpdate: %v", err)
// 		return err
// 	}

// 	defer func() {
// 		if err != nil {
// 			_ = tx.Rollback()
// 		} else {
// 			_ = tx.Commit()
// 		}
// 	}()

// 	query := "INSERT INTO incident_reports(account_id, insu_id, incl_id, description, event_type, address, city, province, postal_code, latitude, longitude, subcategory_name, is_anonymous, media_url, subcategory_code, category_code, created_at) VALUES( ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, NOW())"
// 	result, err := tx.Exec(query,
// 		incident.AccountId,
// 		incident.InsuId,
// 		incident.InclId,
// 		incident.Description,
// 		incident.EventType,
// 		incident.Address,
// 		incident.City,
// 		incident.Province,
// 		incident.PostalCode,
// 		incident.Latitude,
// 		incident.Longitude,
// 		incident.SubCategoryName,
// 		incident.IsAnonymous,
// 		incident.Media.Uri,
// 		incident.SubcategoryCode,
// 		incident.CategoryCode,
// 	)

// 	if err != nil {
// 		return err
// 	}

// 	id, err := result.LastInsertId()
// 	if err != nil {
// 		return err
// 	}

// 	// ----------------------CITIZEN SCORE----------------------- //
// 	// -- Save to DB
// 	err = common.SaveScore(tx, incident.AccountId, 20)
// 	if err != nil {
// 		fmt.Println("error saving score") // It's not necesary to stop the server
// 	}
// 	// ----------------------NOTIFICATION----------------------- //
// 	// -- Save to DB
// 	err = common.SaveNotification(tx, "new_cluster", incident.AccountId, id)

// 	if err != nil {
// 		fmt.Println("error saving notification on createincident event")
// 		return err
// 	}

// 	err = UpdateClusterOnNewIncidentCluster(tx, id, *incident.Vote)
// 	if err != nil {
// 		fmt.Println("error saving notification on createincident event")
// 		return err
// 	}

// 	err = UpdateCounterIncidentsAccount(tx, incident.AccountId)
// 	if err != nil {
// 		fmt.Println("error saving notification on createincident event")
// 		return err
// 	}

// 	return nil
// }

func UpdateClusterOnNewIncidentCluster(tx *sql.Tx, InclID int64, vote bool) error {

	query := `UPDATE incident_clusters 
	SET counter_total_votes = counter_total_incidents_created + 1, counter_total_votes_true = counter_total_votes_true + 1 
	WHERE incl_id = ?`

	if vote != true {
		query = `UPDATE incident_clusters 
		SET counter_total_votes = counter_total_incidents_created + 1, counter_total_votes_false = counter_total_votes_false + 1 
		WHERE incl_id = ?`
	}
	_, err := tx.Exec(query, InclID)
	if err != nil {
		return fmt.Errorf("failed to update score: %w", err)
	}
	return nil
}

func UpdateCounterIncidentsAccount(tx *sql.Tx, accountID int64) error {
	query := `UPDATE account SET counter_total_incidents_created = counter_total_incidents_created+1  WHERE account_id = ?`
	_, err := tx.Exec(query, accountID)
	if err != nil {
		return fmt.Errorf("failed to update score: %w", err)
	}
	return nil
}

func (r *mysqlRepository) HasAccountVoted(inclID, accountID int64) (bool, bool, error) {
	var votedVal sql.NullBool
	err := r.db.QueryRow(
		`SELECT is_true FROM incident_reports WHERE incl_id = ? AND account_id = ? LIMIT 1`,
		inclID, accountID,
	).Scan(&votedVal)
	if err == sql.ErrNoRows {
		return false, false, nil
	}
	if err != nil {
		return false, false, err
	}
	return true, votedVal.Bool, nil
}

func (r *mysqlRepository) UpdateClusterLocation(inclId int64, latitude, longitude float32) (sql.Result, error) {
	query := `
    UPDATE incident_clusters 
    SET 
      center_latitude  = (center_latitude  + ?) / 2,
      center_longitude = (center_longitude + ?) / 2
    WHERE incl_id = ?;`
	return r.db.Exec(query, latitude, longitude, inclId)
}
