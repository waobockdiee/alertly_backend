package reportincident

import (
	"database/sql"
	"log"
)

type Repository interface {
	ReportIncident(report Report) error
}

type pgRepository struct {
	db *sql.DB
}

func NewRepository(db *sql.DB) Repository {
	return &pgRepository{db: db}
}

func (r *pgRepository) ReportIncident(report Report) error {

	query := `INSERT INTO incident_flags (account_id, inre_id, reason, created_at) VALUES ($1, $2, $3, NOW())`
	_, err := r.db.Exec(query, report.AccountID, report.InreID, report.Reason)

	if err != nil {
		log.Printf("Error inserting incident flag: %v", err)
		return err
	}

	query = `UPDATE incident_clusters SET counter_total_flags = counter_total_flags + 1 WHERE incl_id = $1`
	_, err = r.db.Exec(query, report.InclID)

	if err != nil {
		log.Printf("Error inserting incident flag: %v", err)
		return err
	}

	query = `UPDATE incident_reports SET counter_total_flags = counter_total_flags + 1 WHERE inre_id = $1`
	_, err = r.db.Exec(query, report.InreID)

	return err
}
