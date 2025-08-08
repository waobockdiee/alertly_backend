package cjuserank

import "time"

// RankItem representa la estructura de un rango de usuario.
type RankItem struct {
	Title          string `json:"title"`
	Code           string `json:"code"`
	BackgroundColor string `json:"background_color"`
	TextColor      string `json:"text_color"`
	Description    string `json:"description"`
	Image          string `json:"image"` // Usaremos la ruta de la imagen
	ScoreMin       int    `json:"score_min"`
	ScoreMax       int    `json:"score_max"`
}

// UserScore representa el score de un usuario de la tabla account.
type UserScore struct {
	AccountID int64
	Score     int
}

// EarnedRank representa un rango que un usuario ya ha ganado de account_achievements.
type EarnedRank struct {
	AccountID     int64
	Name          string
	Type          string
	BadgeThreshold int // Reutilizamos este campo para almacenar el score mínimo del rango
	CreatedAt     time.Time
}
