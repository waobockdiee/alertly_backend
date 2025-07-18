package cronjobs

import (
	"database/sql"
	"fmt"
)

type Repository interface {
	SetClusterToInactiveAndSetAccountScore() error
	GetDeviceTokensForAccount(accountID int64) ([]string, error)
}

type mysqlRepository struct {
	db *sql.DB
}

func NewRepository(db *sql.DB) Repository {
	return &mysqlRepository{db: db}
}

func (r *Repository) GetDeviceTokensForAccount(accountID int64) ([]string, error) {
	rows, err := r.db.Query(
		`SELECT device_token
           FROM device_tokens
          WHERE account_id = ?`,
		accountID,
	)
	if err != nil {
		return nil, fmt.Errorf("GetDeviceTokensForAccount: %w", err)
	}
	defer rows.Close()

	var tokens []string
	for rows.Next() {
		var t string
		if err := rows.Scan(&t); err != nil {
			return nil, fmt.Errorf("scanning device_token: %w", err)
		}
		tokens = append(tokens, t)
	}
	return tokens, nil
}

/*
actualiza is_active = false, esto quiere decir que ya pasaron 48 horas y el incidente ya no se mostrara en el mapa.
credibilidad actualizada del usuario al asignarla a la columna credibility
credibilidad_usuario = (credibilidad_usuario * 0.8) + (credibilidad_cluster * 0.2)
*/
func (r *mysqlRepository) SetClusterToInactiveAndSetAccountScore() error {
	query := `
	UPDATE incident_clusters ic
	JOIN incident_reports ir ON ir.incl_id = ic.incl_id
	JOIN account a ON a.account_id = ir.account_id
	SET 
		ic.is_active = false,
		a.credibility = (a.credibility * 0.8) + (ir.credibility * 0.2)
	WHERE 
		TIMESTAMPDIFF(HOUR, ic.created_at, NOW()) > 48 
		AND ic.is_active = true;
	`
	_, err := r.db.Exec(query)

	if err != nil {
		return err
	}

	return nil
}
