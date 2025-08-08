package cjuserank

import (
	"database/sql"
	"fmt"
	"time"
)

// Repository encapsula el acceso a la base de datos para el cronjob de rangos de usuario.
type Repository struct {
	db *sql.DB
}

// NewRepository crea una nueva instancia de Repository.
func NewRepository(db *sql.DB) *Repository {
	return &Repository{db: db}
}

// FetchUsersScore obtiene el score de todos los usuarios activos que reciben notificaciones.
func (r *Repository) FetchUsersScore() ([]UserScore, error) {
	query := `
        SELECT
            account_id,
            score
        FROM
            account
        WHERE
            status = 'active' AND receive_notifications = 1
    `

	rows, err := r.db.Query(query)
	if err != nil {
		return nil, fmt.Errorf("FetchUsersScore: %w", err)
	}
	defer rows.Close()

	var usersScore []UserScore
	for rows.Next() {
		var us UserScore
		if err := rows.Scan(&us.AccountID, &us.Score); err != nil {
			return nil, fmt.Errorf("scanning user score: %w", err)
		}
		usersScore = append(usersScore, us)
	}

	return usersScore, nil
}

// GetEarnedRanksForUser obtiene todos los rangos que un usuario ya ha ganado.
func (r *Repository) GetEarnedRanksForUser(accountID int64) ([]EarnedRank, error) {
	query := `
        SELECT
            account_id, name, type, badge_threshold, created
        FROM
            account_achievements
        WHERE
            account_id = ? AND type = 'user_rank'
    `

	rows, err := r.db.Query(query, accountID)
	if err != nil {
		return nil, fmt.Errorf("GetEarnedRanksForUser: %w", err)
	}
	defer rows.Close()

	var earnedRanks []EarnedRank
	for rows.Next() {
		var er EarnedRank
		if err := rows.Scan(&er.AccountID, &er.Name, &er.Type, &er.BadgeThreshold, &er.CreatedAt); err != nil {
			return nil, fmt.Errorf("scanning earned rank: %w", err)
		}
		earnedRanks = append(earnedRanks, er)
	}

	return earnedRanks, nil
}

// InsertEarnedRank inserta un nuevo registro en account_achievements para un rango.
func (r *Repository) InsertEarnedRank(accountID int64, rank RankItem) error {
	query := `
        INSERT INTO account_achievements
            (account_id, name, description, type, icon_url, badge_threshold, created)
        VALUES
            (?, ?, ?, ?, ?, ?, ?)
    `
	_, err := r.db.Exec(
		query,
		accountID,
		rank.Title,
		rank.Description,
		"user_rank", // Usamos un tipo fijo para los rangos de usuario
		rank.Image, // Asumiendo que 'Image' del JSON es 'icon_url'
		rank.ScoreMin, // Usamos ScoreMin como el badge_threshold
		time.Now(),
	)
	if err != nil {
		return fmt.Errorf("InsertEarnedRank: %w", err)
	}
	return nil
}

// InsertNotification inserta una nueva notificación de tipo 'badge_earned' (o similar).
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
		"badge_earned", // O un nuevo tipo como "user_rank_achieved"
		time.Now(),
		1, // must_be_processed = 1 para que el cronjob de notificaciones la procese
	)
	if err != nil {
		return fmt.Errorf("InsertNotification: %w", err)
	}
	return nil
}
