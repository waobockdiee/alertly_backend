package cjnewcluster

import (
	"alertly/internal/cronjobs/shared"
	"database/sql"
	"log"
)

// SubscribedUser representa un usuario que debe ser notificado
type SubscribedUser struct {
	DeviceToken     string
	AccountID       int64
	LocationTitle   string
	SubcategoryName string
}

// Repository encapsula acceso a BD
type Repository struct {
	db *sql.DB
}

func NewRepository(db *sql.DB) *Repository {
	return &Repository{db: db}
}

// FetchPending obtiene notificaciones no procesadas en batch
func (r *Repository) FetchPending(limit int64) ([]Notification, error) {
	query := `SELECT noti_id, reference_id, created_at FROM notifications
         WHERE type = 'new_cluster' AND must_be_processed = 1
         ORDER BY created_at
         LIMIT $1`

	// DEBUG: Log the query
	log.Printf("🔍 Executing query: %s (limit=%d)", query, limit)

	rows, err := r.db.Query(query, limit)
	if err != nil {
		log.Printf("❌ Query error: %v", err)
		return nil, err
	}
	defer rows.Close()

	var list []Notification
	for rows.Next() {
		var n Notification
		if err := rows.Scan(&n.ID, &n.ClusterID, &n.CreatedAt); err != nil {
			log.Printf("❌ Scan error: %v", err)
			return nil, err
		}
		log.Printf("✅ Found notification: noti_id=%d, cluster_id=%d", n.ID, n.ClusterID)
		list = append(list, n)
	}
	log.Printf("📋 Total notifications found: %d", len(list))
	return list, nil
}

// FindSubscribedUsersForCluster encuentra usuarios que deben ser notificados sobre un cluster
func (r *Repository) FindSubscribedUsersForCluster(clusterID int64) ([]SubscribedUser, error) {
	query := `
        SELECT
            dt.device_token,
            a.account_id,
            afl.title AS location_title,
            ic.subcategory_name
        FROM
            incident_clusters ic
        JOIN
            account_favorite_locations afl ON
            -- Fórmula de Haversine para calcular la distancia en KM usando PostgreSQL
            ST_DistanceSphere(
                ST_MakePoint(ic.center_longitude, ic.center_latitude),
                ST_MakePoint(afl.longitude, afl.latitude)
            ) <= afl.radius
        JOIN
            account a ON afl.account_id = a.account_id
        JOIN
            device_tokens dt ON a.account_id = dt.account_id
        WHERE
            ic.incl_id = $1
            AND a.status = 'active'
            AND TRIM(a.is_premium) = '1'
            AND TRIM(a.receive_notifications) = '1'
            AND TRIM(afl.status) = '1'
            AND CASE
                WHEN ic.category_code = 'crime' THEN TRIM(afl.crime) = '1'
                WHEN ic.category_code = 'traffic_accident' THEN TRIM(afl.traffic_accident) = '1'
                WHEN ic.category_code = 'medical_emergency' THEN TRIM(afl.medical_emergency) = '1'
                WHEN ic.category_code = 'fire_incident' THEN TRIM(afl.fire_incident) = '1'
                WHEN ic.category_code = 'vandalism' THEN TRIM(afl.vandalism) = '1'
                WHEN ic.category_code = 'suspicious_activity' THEN TRIM(afl.suspicious_activity) = '1'
                WHEN ic.category_code = 'infrastructure_issues' THEN TRIM(afl.infrastructure_issues) = '1'
                WHEN ic.category_code = 'extreme_weather' THEN TRIM(afl.extreme_weather) = '1'
                WHEN ic.category_code = 'community_events' THEN TRIM(afl.community_events) = '1'
                WHEN ic.category_code = 'dangerous_wildlife_sighting' THEN TRIM(afl.dangerous_wildlife_sighting) = '1'
                WHEN ic.category_code = 'positive_actions' THEN TRIM(afl.positive_actions) = '1'
                WHEN ic.category_code = 'lost_pet' THEN TRIM(afl.lost_pet) = '1'
                ELSE false
            END
    `

	rows, err := r.db.Query(query, clusterID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []SubscribedUser
	for rows.Next() {
		var u SubscribedUser
		if err := rows.Scan(&u.DeviceToken, &u.AccountID, &u.LocationTitle, &u.SubcategoryName); err != nil {
			return nil, err
		}
		users = append(users, u)
	}
	return users, nil
}

// MarkProcessed actualiza processed=true
func (r *Repository) MarkProcessed(ids []int64) error {
	return shared.MarkItemsAsProcessed(r.db, "notifications", "noti_id", ids)
}

// InsertDeliveries inserta en batch registros de envío
func (r *Repository) InsertDeliveries(deliveries []shared.Delivery) error {
	return shared.InsertDeliveries(r.db, deliveries)
}
