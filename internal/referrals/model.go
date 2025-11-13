package referrals

import "time"

// ========================================
// MODELS - Database Entities
// ========================================

// Influencer representa un influencer/marketer
type Influencer struct {
	ID               int64     `db:"id" json:"id"`
	WebInfluencerID  int       `db:"web_influencer_id" json:"web_influencer_id"`
	ReferralCode     string    `db:"referral_code" json:"referral_code"`
	Name             string    `db:"name" json:"name"`
	Platform         string    `db:"platform" json:"platform"` // Instagram, TikTok, Reddit, Other
	IsActive         bool      `db:"is_active" json:"is_active"`
	CreatedAt        time.Time `db:"created_at" json:"created_at"`
	UpdatedAt        time.Time `db:"updated_at" json:"updated_at"`
}

// ReferralConversion representa un registro con código de referral
type ReferralConversion struct {
	ID           int64     `db:"id" json:"id"`
	ReferralCode string    `db:"referral_code" json:"referral_code"`
	UserID       int64     `db:"user_id" json:"user_id"`
	RegisteredAt time.Time `db:"registered_at" json:"registered_at"`
	Platform     string    `db:"platform" json:"platform"` // iOS, Android
	Earnings     float64   `db:"earnings" json:"earnings"`
	CreatedAt    time.Time `db:"created_at" json:"created_at"`
}

// ReferralPremiumConversion representa una suscripción premium de usuario referido
type ReferralPremiumConversion struct {
	ID                   int64     `db:"id" json:"id"`
	ReferralCode         string    `db:"referral_code" json:"referral_code"`
	UserID               int64     `db:"user_id" json:"user_id"`
	ConversionID         *int64    `db:"conversion_id" json:"conversion_id"`
	SubscriptionType     string    `db:"subscription_type" json:"subscription_type"` // monthly, yearly
	Amount               float64   `db:"amount" json:"amount"`
	Commission           float64   `db:"commission" json:"commission"`
	CommissionPercentage float64   `db:"commission_percentage" json:"commission_percentage"`
	ConvertedAt          time.Time `db:"converted_at" json:"converted_at"`
	CreatedAt            time.Time `db:"created_at" json:"created_at"`
}

// ReferralMetricsCache representa el cache de métricas
type ReferralMetricsCache struct {
	ReferralCode            string    `db:"referral_code" json:"referral_code"`
	TotalRegistrations      int       `db:"total_registrations" json:"total_registrations"`
	TotalPremiumConversions int       `db:"total_premium_conversions" json:"total_premium_conversions"`
	TotalEarnings           float64   `db:"total_earnings" json:"total_earnings"`
	LastUpdated             time.Time `db:"last_updated" json:"last_updated"`
}

// ========================================
// REQUEST MODELS - API Inputs
// ========================================

// ValidateReferralCodeRequest es el request para validar un código
type ValidateReferralCodeRequest struct {
	ReferralCode string `json:"referral_code" binding:"required"`
}

// RegisterConversionRequest es el request para registrar una conversión de registro
type RegisterConversionRequest struct {
	ReferralCode string `json:"referral_code" binding:"required"`
	UserID       int64  `json:"user_id" binding:"required"`
	RegisteredAt string `json:"registered_at" binding:"required"` // ISO 8601
	Platform     string `json:"platform" binding:"required,oneof=iOS Android"`
}

// RegisterPremiumConversionRequest es el request para registrar una conversión premium
type RegisterPremiumConversionRequest struct {
	UserID           int64   `json:"user_id" binding:"required"`
	ReferralCode     string  `json:"referral_code"`                               // Opcional si se puede buscar por user_id
	SubscriptionType string  `json:"subscription_type" binding:"required,oneof=monthly yearly"`
	Amount           float64 `json:"amount" binding:"required,gt=0"`
	ConvertedAt      string  `json:"converted_at" binding:"required"` // ISO 8601
}

// SyncInfluencerRequest es el request para sincronizar un influencer desde backend web
type SyncInfluencerRequest struct {
	WebInfluencerID int    `json:"web_influencer_id" binding:"required"`
	ReferralCode    string `json:"referral_code" binding:"required"`
	Name            string `json:"name" binding:"required"`
	Platform        string `json:"platform" binding:"required,oneof=Instagram TikTok Reddit Other"`
	IsActive        bool   `json:"is_active"`
}

// ========================================
// RESPONSE MODELS - API Outputs
// ========================================

// ValidateReferralCodeResponse es la respuesta de validación
type ValidateReferralCodeResponse struct {
	Valid            bool   `json:"valid"`
	InfluencerID     *int64 `json:"influencer_id,omitempty"`
	InfluencerName   string `json:"influencer_name,omitempty"`
	PremiumTrialDays int    `json:"premium_trial_days,omitempty"`
	Message          string `json:"message,omitempty"`
}

// RegisterConversionResponse es la respuesta de registro de conversión
type RegisterConversionResponse struct {
	Success                 bool    `json:"success"`
	Message                 string  `json:"message,omitempty"`
	InfluencerEarningsAdded float64 `json:"influencer_earnings_added,omitempty"`
	Error                   string  `json:"error,omitempty"`
}

// RegisterPremiumConversionResponse es la respuesta de conversión premium
type RegisterPremiumConversionResponse struct {
	Success              bool    `json:"success"`
	Message              string  `json:"message,omitempty"`
	InfluencerCommission float64 `json:"influencer_commission,omitempty"`
	CommissionPercentage float64 `json:"commission_percentage,omitempty"`
	Error                string  `json:"error,omitempty"`
}

// SyncInfluencerResponse es la respuesta de sincronización
type SyncInfluencerResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message,omitempty"`
	Error   string `json:"error,omitempty"`
}

// DailyMetric representa métricas de un día específico
type DailyMetric struct {
	Date                string  `json:"date"`
	Registrations       int     `json:"registrations"`
	PremiumConversions  int     `json:"premium_conversions"`
	Earnings            float64 `json:"earnings"`
}

// InfluencerMetricsResponse es la respuesta de métricas individuales
type InfluencerMetricsResponse struct {
	ReferralCode              string        `json:"referral_code"`
	InfluencerID              int64         `json:"influencer_id"`
	TotalRegistrations        int           `json:"total_registrations"`
	TotalPremiumConversions   int           `json:"total_premium_conversions"`
	TotalEarnings             float64       `json:"total_earnings"`
	CurrentMonthRegistrations int           `json:"current_month_registrations"`
	CurrentMonthPremium       int           `json:"current_month_premium"`
	CurrentMonthEarnings      float64       `json:"current_month_earnings"`
	ProjectedMonthEarnings    float64       `json:"projected_month_earnings"`
	DailyMetrics              []DailyMetric `json:"daily_metrics"`
	Rank                      int           `json:"rank"`
	TotalInfluencers          int           `json:"total_influencers"`
}

// TopPerformer representa un influencer top
type TopPerformer struct {
	InfluencerID            int64   `json:"influencer_id"`
	ReferralCode            string  `json:"referral_code"`
	Name                    string  `json:"name"`
	TotalRegistrations      int     `json:"total_registrations"`
	TotalPremiumConversions int     `json:"total_premium_conversions"`
	TotalEarnings           float64 `json:"total_earnings"`
	Platform                string  `json:"platform"`
}

// MonthlyTrend representa métricas mensuales
type MonthlyTrend struct {
	Month          string  `json:"month"`
	Registrations  int     `json:"registrations"`
	Premium        int     `json:"premium"`
	Earnings       float64 `json:"earnings"`
}

// PlatformBreakdown representa métricas por plataforma
type PlatformBreakdown struct {
	Influencers   int     `json:"influencers"`
	Registrations int     `json:"registrations"`
	Premium       int     `json:"premium"`
	Earnings      float64 `json:"earnings"`
}

// AggregateMetricsResponse es la respuesta de métricas agregadas
type AggregateMetricsResponse struct {
	TotalReferrals          int                          `json:"total_referrals"`
	TotalPremiumConversions int                          `json:"total_premium_conversions"`
	TotalEarningsPaid       float64                      `json:"total_earnings_paid"`
	ActiveInfluencers       int                          `json:"active_influencers"`
	ConversionRate          float64                      `json:"conversion_rate"`
	TopPerformers           []TopPerformer               `json:"top_performers"`
	MonthlyTrend            []MonthlyTrend               `json:"monthly_trend"`
	PlatformBreakdown       map[string]PlatformBreakdown `json:"platform_breakdown"`
}
