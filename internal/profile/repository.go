package profile

import (
	"database/sql"
	"encoding/json"
	"fmt"
)

type Repository interface {
	GetById(accountID int64) (Profile, error)
	UpdateTotalIncidents(accountID int64) error
}

type mysqlRepository struct {
	db *sql.DB
}

func NewRepository(db *sql.DB) Repository {
	return &mysqlRepository{db: db}
}

func (r *mysqlRepository) GetById(accountID int64) (Profile, error) {
	query := `SELECT 
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
	a.birth_year, 
	a.birth_month, 
	a.birth_day, 
	a.has_finished_tutorial, 
	a.has_watch_new_incident_tutorial, 
	a.thumbnail_url,
	IFNULL(
		(
		  SELECT JSON_ARRAYAGG(
			JSON_OBJECT(
			  'inre_id', i.inre_id,
			  'media_url', i.media_url,
			  'description', i.description,
			  'event_type', i.event_type,
			  'subcategory_name', i.subcategory_name
			)
		  )
		  FROM incident_reports i 
		  WHERE i.account_id = a.account_id
		),
		'[]'
	) AS incidents
	FROM account a WHERE a.account_id = ?`
	var stc Profile
	var rawIncidents string

	err := r.db.QueryRow(query, accountID).Scan(
		&stc.AccountID,
		&stc.Nickname,
		&stc.FirstName,
		&stc.LastName,
		&stc.PhoneNumber,
		&stc.Status,
		&stc.Credibility,
		&stc.IsPrivateProfile,
		&stc.Score,
		&stc.IsPremium,
		&stc.CounterTotalIncidentsCreated,
		&stc.CounterTotalVotesMade,
		&stc.CounterTotalCommentsMade,
		&stc.CounterTotalLocations,
		&stc.CounterTotalFlags,
		&stc.CounterTotalMedals,
		&stc.BirthYear,
		&stc.BirthMonth,
		&stc.BirthDay,
		&stc.HasFinishedTutorial,
		&stc.HasWatchNewIncidentTutorial,
		&stc.ThumbnailUrl,
		&rawIncidents)

	if err != nil {
		return Profile{}, fmt.Errorf("error scanning row: %w", err)
	}

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

func (r *mysqlRepository) UpdateTotalIncidents(accountID int64) error {
	query := `UPDATE account SET counter_total_incidents_created = counter_total_incidents_created + 1 WHERE account_id = ?`
	_, err := r.db.Exec(query, accountID)
	if err != nil {
		return fmt.Errorf("error updating total incidents %w", err)
	}

	return nil
}
