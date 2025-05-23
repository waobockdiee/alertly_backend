package myplaces

import (
	"database/sql"
	"fmt"
)

type Repository interface {
	Get(accountId int) ([]MyPlaces, error)
	Add(myPlace MyPlaces) (int64, error)
	Update(myPlace MyPlaces) error
	GetByAccountId(accountId int64) ([]MyPlaces, error)
	GetById(accountId, aflId int64) (MyPlaces, error)
	FullUpdate(myPlace MyPlaces) error
	Delete(accountID, aflID int64) error
}

type mysqlRepository struct {
	db *sql.DB
}

func NewRepository(db *sql.DB) Repository {
	return &mysqlRepository{db: db}
}

func (r *mysqlRepository) Get(accountId int) ([]MyPlaces, error) {

	query := `
	SELECT afl_id, account_id, title, latitude, longitude, city, province, postal_code, status, radius FROM account_favorite_locations WHERE account_id = ? ORDER BY afl_id DESC;
	`
	var myPlaces []MyPlaces

	rows, err := r.db.Query(query, accountId)
	if err != nil {
		return myPlaces, err
	}
	defer rows.Close()

	for rows.Next() {
		var place MyPlaces
		if err := rows.Scan(&place.AflId, &place.AccountId, &place.Title, &place.Latitude, &place.Longitude, &place.City, &place.Province, &place.PostalCode, &place.Status, &place.Radius); err != nil {
			return myPlaces, err
		}
		myPlaces = append(myPlaces, place)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	fmt.Printf("%v", myPlaces)

	return myPlaces, nil
}

func (r *mysqlRepository) Add(myPlace MyPlaces) (int64, error) {
	query := "INSERT INTO account_favorite_locations(account_id, title, latitude, longitude, city, province, postal_code, crime, traffic_accident, medical_emergency, fire_incident, vandalism, suspicious_activity, infrastructure_issues, extreme_weather, community_events, dangerous_wildlife_sighting, positive_actions, lost_pet, radius) VALUES(?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)"

	result, err := r.db.Exec(query,
		myPlace.AccountId,
		myPlace.Title,
		myPlace.Latitude,
		myPlace.Longitude,
		myPlace.City,
		myPlace.Province,
		myPlace.PostalCode,
		myPlace.Crime,
		myPlace.TrafficAccident,
		myPlace.MedicalEmergency,
		myPlace.FireIncident,
		myPlace.Vandalism,
		myPlace.SuspiciousActivity,
		myPlace.InfrastructureIssues,
		myPlace.ExtremeWeather,
		myPlace.CommunityEvents,
		myPlace.DangerousWildlifeSighting,
		myPlace.PositiveActions,
		myPlace.LostPet,
		myPlace.Radius,
	)

	fmt.Println("DEBUG", err)

	if err != nil {
		return 0, err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}

	fmt.Printf("Saved succesfully with ID: %d\n", id)
	return id, nil
}

func (r *mysqlRepository) Update(myPlace MyPlaces) error {
	query := `UPDATE account_favorite_locations SET status = ? WHERE afl_id = ?`
	_, err := r.db.Exec(query, myPlace.Status, myPlace.AflId)
	if err != nil {
		return fmt.Errorf("error updating cluster %w", err)
	}
	return nil
}

func (r *mysqlRepository) FullUpdate(myPlace MyPlaces) error {
	query := `UPDATE account_favorite_locations 
	SET title =?,
	crime = ?,
	traffic_accident = ?,
	medical_emergency = ?,
	fire_incident = ?,
	vandalism = ?,
	suspicious_activity = ?,
	infrastructure_issues = ?,
	extreme_weather = ?,
	community_events = ?,
	dangerous_wildlife_sighting = ?,
	positive_actions = ?,
	lost_pet = ?
	WHERE afl_id = ? AND account_id = ?`
	_, err := r.db.Exec(query,
		myPlace.Title,
		myPlace.Crime,
		myPlace.TrafficAccident,
		myPlace.MedicalEmergency,
		myPlace.FireIncident,
		myPlace.Vandalism,
		myPlace.SuspiciousActivity,
		myPlace.InfrastructureIssues,
		myPlace.ExtremeWeather,
		myPlace.CommunityEvents,
		myPlace.DangerousWildlifeSighting,
		myPlace.PositiveActions,
		myPlace.LostPet,
		myPlace.AflId,
		myPlace.AccountId,
	)
	if err != nil {
		return fmt.Errorf("error updating cluster  %w", err)
	}
	return nil
}

func (r *mysqlRepository) GetById(accountId, aflId int64) (MyPlaces, error) {
	query := `SELECT afl_id, account_id, title, latitude, longitude, city, province, postal_code, status, crime, traffic_accident, medical_emergency, fire_incident, vandalism, suspicious_activity, infrastructure_issues, extreme_weather, community_events, dangerous_wildlife_sighting, positive_actions, lost_pet FROM account_favorite_locations WHERE account_id = ? AND afl_id = ?`

	var c MyPlaces
	err := r.db.QueryRow(query, accountId, aflId).Scan(&c.AflId,
		&c.AccountId,
		&c.Title,
		&c.Latitude,
		&c.Longitude,
		&c.City,
		&c.Province,
		&c.PostalCode,
		&c.Status,
		&c.Crime,
		&c.TrafficAccident,
		&c.MedicalEmergency,
		&c.FireIncident,
		&c.Vandalism,
		&c.SuspiciousActivity,
		&c.InfrastructureIssues,
		&c.ExtremeWeather,
		&c.CommunityEvents,
		&c.DangerousWildlifeSighting,
		&c.PositiveActions,
		&c.LostPet)

	if err != nil {
		return MyPlaces{}, fmt.Errorf("error scanning row: %w", err)
	}

	return c, nil
}

func (r *mysqlRepository) GetByAccountId(accountId int64) ([]MyPlaces, error) {
	query := `SELECT afl_id, account_id, title, status, city, latitude, longitude, radius FROM account_favorite_locations WHERE account_id = ? ORDER BY afl_id DESC`
	rows, err := r.db.Query(query, accountId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var places []MyPlaces
	for rows.Next() {
		var c MyPlaces
		if err := rows.Scan(
			&c.AflId,
			&c.AccountId,
			&c.Title,
			&c.Status,
			&c.City,
			&c.Latitude,
			&c.Longitude,
			&c.Radius,
		); err != nil {
			return nil, err
		}
		places = append(places, c)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}
	return places, nil
}

func (r *mysqlRepository) Delete(accountID, aflID int64) error {
	query := `DELETE FROM account_favorite_locations WHERE account_id = ? AND afl_id = ?`
	_, err := r.db.Exec(query,
		accountID,
		aflID,
	)
	if err != nil {
		return fmt.Errorf("error deleting place %w", err)
	}
	return nil
}
