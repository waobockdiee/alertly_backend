package cjblockincident

import (
	"database/sql"
	"fmt"
)

// Repository encapsula el acceso a la base de datos para el cronjob de bloqueo de incidentes.
type Repository struct {
	db *sql.DB
}

// NewRepository crea una nueva instancia de Repository.
func NewRepository(db *sql.DB) *Repository {
	return &Repository{db: db}
}

// FetchIncidentsToReject obtiene los incidentes que han sido reportados más de 5 veces
// y cuyo estado actual no es 'rejected'.
func (r *Repository) FetchIncidentsToReject() ([]IncidentToReject, error) {
	query := `
        SELECT
            inf.inre_id,
            COUNT(inf.infl_id) AS flag_count
        FROM
            incident_flags inf
        JOIN
            incident_reports ir ON inf.inre_id = ir.inre_id
        WHERE
            (ir.status IS NULL OR ir.status != 'rejected')
        GROUP BY
            inf.inre_id
        HAVING
            COUNT(inf.infl_id) > 5
    `

	rows, err := r.db.Query(query)
	if err != nil {
		return nil, fmt.Errorf("FetchIncidentsToReject: %w", err)
	}
	defer rows.Close()

	var incidentsToReject []IncidentToReject
	for rows.Next() {
		var itr IncidentToReject
		if err := rows.Scan(&itr.IncidentID, &itr.FlagCount); err != nil {
			return nil, fmt.Errorf("scanning incident to reject: %w", err)
		}
		incidentsToReject = append(incidentsToReject, itr)
	}

	return incidentsToReject, nil
}

// RejectIncident actualiza el estado de un incidente a 'rejected'.
func (r *Repository) RejectIncident(incidentID int64) error {
	query := `
        UPDATE incident_reports
        SET status = 'rejected'
        WHERE inre_id = ?
    `
	_, err := r.db.Exec(query, incidentID)
	if err != nil {
		return fmt.Errorf("RejectIncident: %w", err)
	}
	return nil
}
