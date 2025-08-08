package cjbadgeearn

import "time"

// Badge representa la estructura de una insignia definida en el JSON.
type Badge struct {
	Title       string `json:"title"`
	Category    string `json:"category"`
	Image       string `json:"image"`
	Description string `json:"description"`
	Number      int    `json:"number"`
	Code        string `json:"code"`
}

// UserActivity representa los contadores de actividad de un usuario de la tabla account.
type UserActivity struct {
	AccountID                 int64
	CounterTotalIncidentsCreated int
	IncidentAsUpdate          int
	Crime                     int
	TrafficAccident           int
	MedicalEmergency          int
	FireIncident              int
	Vandalism                 int
	SuspiciousActivity        int
	InfrastructureIssues      int
	ExtremeWeather            int
	CommunityEvents           int
	DangerousWildlifeSighting int
	PositiveActions           int
	LostPet                   int
}

// EarnedBadge representa una insignia que un usuario ya ha ganado de account_achievements.
type EarnedBadge struct {
	AccountID     int64
	Name          string
	Type          string
	BadgeThreshold int
	CreatedAt     time.Time
}
