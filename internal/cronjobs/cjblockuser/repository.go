package cjblockuser

import (
	"database/sql"
	"fmt"
)

// Repository encapsula el acceso a la base de datos para el cronjob de bloqueo de usuarios.
type Repository struct {
	db *sql.DB
}

// NewRepository crea una nueva instancia de Repository.
func NewRepository(db *sql.DB) *Repository {
	return &Repository{db: db}
}

// FetchUsersToBlock obtiene los usuarios que han sido reportados más de 20 veces
// y cuyo estado actual no es 'blocked'.
func (r *Repository) FetchUsersToBlock() ([]UserToBlock, error) {
	query := `
        SELECT
            ar.account_id,
            COUNT(ar.acre_id) AS report_count
        FROM
            account_reports ar
        JOIN
            account a ON ar.account_id = a.account_id
        WHERE
            a.status != 'blocked'
        GROUP BY
            ar.account_id
        HAVING
            COUNT(ar.acre_id) > 20
    `

	rows, err := r.db.Query(query)
	if err != nil {
		return nil, fmt.Errorf("FetchUsersToBlock: %w", err)
	}
	defer rows.Close()

	var usersToBlock []UserToBlock
	for rows.Next() {
		var utb UserToBlock
		if err := rows.Scan(&utb.AccountID, &utb.ReportCount); err != nil {
			return nil, fmt.Errorf("scanning user to block: %w", err)
		}
		usersToBlock = append(usersToBlock, utb)
	}

	return usersToBlock, nil
}

// BlockUser actualiza el estado de un usuario a 'blocked'.
func (r *Repository) BlockUser(accountID int64) error {
	query := `
        UPDATE account
        SET status = 'blocked'
        WHERE account_id = $1
    `
	_, err := r.db.Exec(query, accountID)
	if err != nil {
		return fmt.Errorf("BlockUser: %w", err)
	}
	return nil
}
