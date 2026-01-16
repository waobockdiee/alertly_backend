package profile

import (
	"database/sql"
	"encoding/json"
	"fmt"

	"alertly/internal/dbtypes"
)

type Repository interface {
	GetById(accountID int64) (Profile, error)
	UpdateTotalIncidents(accountID int64) error
	ReportAccount(report ReportAccountInput) error
}

type pgRepository struct {
	db *sql.DB
}

func NewRepository(db *sql.DB) Repository {
	return &pgRepository{db: db}
}

func (r *pgRepository) GetById(accountID int64) (Profile, error) {
	query := `
		SELECT
			a.account_id,
			a.nickname,
			a.first_name,
			a.last_name,
			a.phone_number,
			a.status,
			a.credibility,
			a.is_private_profile,
			a.score,
			a.is_premium,
			a.counter_total_incidents_created,
			a.counter_total_votes_made,
			a.counter_total_comments_made,
			a.counter_total_locations,
			a.counter_total_flags,
			a.counter_total_medals,
			COALESCE(a.birth_year, '') as birth_year,
			COALESCE(a.birth_month, '') as birth_month,
			COALESCE(a.birth_day, '') as birth_day,
			a.has_finished_tutorial,
			a.has_watch_new_incident_tutorial,
			COALESCE(a.thumbnail_url, '') as thumbnail_url,
			a.crime,
			a.traffic_accident,
			a.medical_emergency,
			a.fire_incident,
			a.vandalism,
			a.suspicious_activity,
			a.infrastructure_issues,
			a.extreme_weather,
			a.community_events,
			a.dangerous_wildlife_sighting,
			a.positive_actions,
			a.lost_pet,
			a.incident_as_update,
			COALESCE(
				(
				SELECT JSON_AGG(
					JSON_BUILD_OBJECT(
					'inre_id', i.inre_id,
					'media_url', i.media_url,
					'description', i.description,
					'event_type', i.event_type,
					'subcategory_name', i.subcategory_name,
					'credibility', ic.credibility,
					'incl_id', i.incl_id,
					'is_anonymous', i.is_anonymous,
					'created_at', i.created_at
					)
				)
				FROM incident_reports i
				INNER JOIN incident_clusters ic
					ON i.incl_id = ic.incl_id
				WHERE i.account_id = a.account_id
				ORDER BY i.created_at DESC
				),
				'[]'::json
			) AS incidents
		FROM account a
		WHERE a.account_id = $1
		`
	var stc Profile
	var rawIncidents string

	// Usar NullBool para campos booleanos que pueden ser SMALLINT/CHAR/BOOLEAN
	var isPrivateProfile, isPremium, hasFinishedTutorial, hasWatchNewIncidentTutorial dbtypes.NullBool

	err := r.db.QueryRow(query, accountID).Scan(
		&stc.AccountID,
		&stc.Nickname,
		&stc.FirstName,
		&stc.LastName,
		&stc.PhoneNumber,
		&stc.Status,
		&stc.Credibility,
		&isPrivateProfile,
		&stc.Score,
		&isPremium,
		&stc.CounterTotalIncidentsCreated,
		&stc.CounterTotalVotesMade,
		&stc.CounterTotalCommentsMade,
		&stc.CounterTotalLocations,
		&stc.CounterTotalFlags,
		&stc.CounterTotalMedals,
		&stc.BirthYear,
		&stc.BirthMonth,
		&stc.BirthDay,
		&hasFinishedTutorial,
		&hasWatchNewIncidentTutorial,
		&stc.ThumbnailUrl,
		&stc.Crime,
		&stc.TrafficAccident,
		&stc.MedicalEmergency,
		&stc.FireIncident,
		&stc.Vandalism,
		&stc.SuspiciousActivity,
		&stc.InfrastructureIssues,
		&stc.ExtremeWeather,
		&stc.CommunityEvents,
		&stc.DangerousWildlifeSighting,
		&stc.PositiveActions,
		&stc.LostPet,
		&stc.IncidentAsUpdate,
		&rawIncidents)

	if err != nil {
		return Profile{}, fmt.Errorf("error scanning row: %w", err)
	}

	// Convertir NullBool a bool
	stc.IsPrivateProfile = isPrivateProfile.Valid && isPrivateProfile.Bool
	stc.IsPremium = isPremium.Valid && isPremium.Bool
	stc.HasFinishedTutorial = hasFinishedTutorial.Valid && hasFinishedTutorial.Bool
	stc.HasWatchNewIncidentTutorial = hasWatchNewIncidentTutorial.Valid && hasWatchNewIncidentTutorial.Bool

	if rawIncidents == "" {
		stc.Incidents = []Incident{}
	} else {
		var incidents []Incident
		if err := json.Unmarshal([]byte(rawIncidents), &incidents); err != nil {
			return Profile{}, fmt.Errorf("error unmarshalling incidents: %w", err)
		}
		stc.Incidents = incidents
	}

	return stc, nil
}

func (r *pgRepository) UpdateTotalIncidents(accountID int64) error {
	query := `UPDATE account SET counter_total_incidents_created = counter_total_incidents_created + 1 WHERE account_id = $1`
	_, err := r.db.Exec(query, accountID)
	if err != nil {
		return fmt.Errorf("error updating total incidents %w", err)
	}

	return nil
}

func (r *pgRepository) ReportAccount(report ReportAccountInput) error {
	query := `INSERT INTO account_reports(account_id_whos_reporting, account_id, message) VALUES($1, $2 ,$3)`
	_, err := r.db.Exec(query, report.AccountIDWhosReporting, report.AccountID, report.Message)
	return err
}
