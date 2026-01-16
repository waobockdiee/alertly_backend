package tutorial

import (
	"database/sql"

	"alertly/internal/dbtypes"
)

type Repository interface {
	MarkTutorialAsFinished(accountID int64) error
}

type pgRepository struct {
	db *sql.DB
}

func NewRepository(db *sql.DB) Repository {
	return &pgRepository{db: db}
}

func (r *pgRepository) MarkTutorialAsFinished(accountID int64) error {
	// Usar dbtypes.BoolToInt para insertar en columnas SMALLINT
	query := `UPDATE account SET has_finished_tutorial = $1 WHERE account_id = $2`
	_, err := r.db.Exec(query, dbtypes.BoolToInt(true), accountID)
	return err
}
