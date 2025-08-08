package cjbadgeearn

import (
	"database/sql"
	"fmt"
	"time"
)

// Repository encapsula el acceso a la base de datos para el cronjob de insignias.
type Repository struct {
	db *sql.DB
}

// NewRepository crea una nueva instancia de Repository.
func NewRepository(db *sql.DB) *Repository {
	return &Repository{db: db}
}

// FetchUsersActivity obtiene los contadores de actividad de todos los usuarios activos.
func (r *Repository) FetchUsersActivity() ([]UserActivity, error) {
	query := `
        SELECT
            account_id,
            counter_total_incidents_created,
            incident_as_update,
            crime,
            traffic_accident,
            medical_emergency,
            fire_incident,
            vandalism,
            suspicious_activity,
            infrastructure_issues,
            extreme_weather,
            community_events,
            dangerous_wildlife_sighting,
            positive_actions,
            lost_pet
        FROM
            account
        WHERE
            status = 'active' AND receive_notifications = 1
    `

	rows, err := r.db.Query(query)
	if err != nil {
		return nil, fmt.Errorf("FetchUsersActivity: %w", err)
	}
	defer rows.Close()

	var usersActivity []UserActivity
	for rows.Next() {
		var ua UserActivity
		if err := rows.Scan(
			&ua.AccountID,
			&ua.CounterTotalIncidentsCreated,
			&ua.IncidentAsUpdate,
			&ua.Crime,
			&ua.TrafficAccident,
			&ua.MedicalEmergency,
			&ua.FireIncident,
			&ua.Vandalism,
			&ua.SuspiciousActivity,
			&ua.InfrastructureIssues,
			&ua.ExtremeWeather,
			&ua.CommunityEvents,
			&ua.DangerousWildlifeSighting,
			&ua.PositiveActions,
			&ua.LostPet,
		); err != nil {
			return nil, fmt.Errorf("scanning user activity: %w", err)
		}
		usersActivity = append(usersActivity, ua)
	}

	return usersActivity, nil
}

// GetEarnedBadgesForUser obtiene todas las insignias que un usuario ya ha ganado.
func (r *Repository) GetEarnedBadgesForUser(accountID int64) ([]EarnedBadge, error) {
	query := `
        SELECT
            account_id, name, type, badge_threshold, created
        FROM
            account_achievements
        WHERE
            account_id = ?
    `

	rows, err := r.db.Query(query, accountID)
	if err != nil {
		return nil, fmt.Errorf("GetEarnedBadgesForUser: %w", err)
	}
	defer rows.Close()

	var earnedBadges []EarnedBadge
	for rows.Next() {
		var eb EarnedBadge
		if err := rows.Scan(&eb.AccountID, &eb.Name, &eb.Type, &eb.BadgeThreshold, &eb.CreatedAt); err != nil {
			return nil, fmt.Errorf("scanning earned badge: %w", err)
		}
		earnedBadges = append(earnedBadges, eb)
	}

	return earnedBadges, nil
}

// InsertEarnedBadge inserta un nuevo registro en account_achievements.
func (r *Repository) InsertEarnedBadge(accountID int64, badge Badge) error {
	query := `
        INSERT INTO account_achievements
            (account_id, name, description, type, icon_url, badge_threshold, created)
        VALUES
            (?, ?, ?, ?, ?, ?, ?)
    `
	_, err := r.db.Exec(
		query,
		accountID,
		badge.Title,
		badge.Description,
		badge.Code, // Usamos Code como el 'type' en la tabla
		badge.Image, // Asumiendo que 'Image' del JSON es 'icon_url'
		badge.Number,
		time.Now(),
	)
	if err != nil {
		return fmt.Errorf("InsertEarnedBadge: %w", err)
	}
	return nil
}

// InsertNotification inserta una nueva notificación de tipo 'badge_earned'.
func (r *Repository) InsertNotification(accountID int64, title string, message string) error {
	query := `
        INSERT INTO notifications
            (owner_account_id, title, message, type, created_at, must_be_processed)
        VALUES
            (?, ?, ?, ?, ?, ?)
    `
	_, err := r.db.Exec(
		query,
		accountID,
		title,
		message,
		"badge_earned",
		time.Now(),
		1, // must_be_processed = 1 para que el cronjob de notificaciones la procese
	)
	if err != nil {
		return fmt.Errorf("InsertNotification: %w", err)
	}
	return nil
}
