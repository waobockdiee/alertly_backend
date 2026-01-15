package referrals

import (
	"database/sql"
	"fmt"
	"time"
)

// Repository define las operaciones de base de datos para el sistema de referrals
type Repository interface {
	// Influencers
	GetInfluencerByCode(code string) (*Influencer, error)
	GetInfluencerByID(id int64) (*Influencer, error)
	UpsertInfluencer(inf *Influencer) error
	GetAllActiveInfluencers() ([]Influencer, error)

	// Conversions
	CreateConversion(conv *ReferralConversion) error
	GetConversionByUserID(userID int64) (*ReferralConversion, error)
	GetConversionsByCode(code string) ([]ReferralConversion, error)

	// Premium Conversions
	CreatePremiumConversion(conv *ReferralPremiumConversion) error
	GetPremiumConversionsByCode(code string) ([]ReferralPremiumConversion, error)

	// Metrics
	GetTotalRegistrationsByCode(code string) (int, error)
	GetTotalPremiumByCode(code string) (int, error)
	GetTotalEarningsByCode(code string) (float64, error)
	GetCurrentMonthRegistrationsByCode(code string) (int, error)
	GetCurrentMonthPremiumByCode(code string) (int, error)
	GetCurrentMonthEarningsByCode(code string) (float64, error)
	GetDailyMetrics(code string, days int) ([]DailyMetric, error)
	GetInfluencerRank(code string) (int, int, error) // rank, total

	// Aggregate Metrics
	GetTotalReferrals() (int, error)
	GetTotalPremiumConversions() (int, error)
	GetTotalEarnings() (float64, error)
	GetActiveInfluencersCount() (int, error)
	GetTopPerformers(limit int) ([]TopPerformer, error)
	GetMonthlyTrend(months int) ([]MonthlyTrend, error)
	GetPlatformBreakdown() (map[string]PlatformBreakdown, error)
}

type pgRepository struct {
	db *sql.DB
}

// NewRepository crea una nueva instancia del repositorio
func NewRepository(db *sql.DB) Repository {
	return &pgRepository{db: db}
}

// ========================================
// INFLUENCERS
// ========================================

func (r *pgRepository) GetInfluencerByCode(code string) (*Influencer, error) {
	query := `
		SELECT id, web_influencer_id, referral_code, name, platform, is_active, created_at, updated_at
		FROM influencers
		WHERE referral_code = $1
	`
	var inf Influencer
	err := r.db.QueryRow(query, code).Scan(
		&inf.ID, &inf.WebInfluencerID, &inf.ReferralCode, &inf.Name,
		&inf.Platform, &inf.IsActive, &inf.CreatedAt, &inf.UpdatedAt,
	)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &inf, nil
}

func (r *pgRepository) GetInfluencerByID(id int64) (*Influencer, error) {
	query := `
		SELECT id, web_influencer_id, referral_code, name, platform, is_active, created_at, updated_at
		FROM influencers
		WHERE id = $1
	`
	var inf Influencer
	err := r.db.QueryRow(query, id).Scan(
		&inf.ID, &inf.WebInfluencerID, &inf.ReferralCode, &inf.Name,
		&inf.Platform, &inf.IsActive, &inf.CreatedAt, &inf.UpdatedAt,
	)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &inf, nil
}

func (r *pgRepository) UpsertInfluencer(inf *Influencer) error {
	query := `
		INSERT INTO influencers (web_influencer_id, referral_code, name, platform, is_active)
		VALUES ($1, $2, $3, $4, $5)
		ON CONFLICT (web_influencer_id) DO UPDATE SET
			name = EXCLUDED.name,
			platform = EXCLUDED.platform,
			is_active = EXCLUDED.is_active,
			updated_at = CURRENT_TIMESTAMP
	`
	_, err := r.db.Exec(query, inf.WebInfluencerID, inf.ReferralCode, inf.Name, inf.Platform, inf.IsActive)
	return err
}

func (r *pgRepository) GetAllActiveInfluencers() ([]Influencer, error) {
	query := `
		SELECT id, web_influencer_id, referral_code, name, platform, is_active, created_at, updated_at
		FROM influencers
		WHERE is_active = '1'
		ORDER BY created_at DESC
	`
	rows, err := r.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var influencers []Influencer
	for rows.Next() {
		var inf Influencer
		err := rows.Scan(
			&inf.ID, &inf.WebInfluencerID, &inf.ReferralCode, &inf.Name,
			&inf.Platform, &inf.IsActive, &inf.CreatedAt, &inf.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		influencers = append(influencers, inf)
	}
	return influencers, nil
}

// ========================================
// CONVERSIONS
// ========================================

func (r *pgRepository) CreateConversion(conv *ReferralConversion) error {
	query := `
		INSERT INTO referral_conversions (referral_code, user_id, registered_at, platform, earnings)
		VALUES ($1, $2, $3, $4, $5)
	`
	_, err := r.db.Exec(query, conv.ReferralCode, conv.UserID, conv.RegisteredAt, conv.Platform, conv.Earnings)
	return err
}

func (r *pgRepository) GetConversionByUserID(userID int64) (*ReferralConversion, error) {
	query := `
		SELECT id, referral_code, user_id, registered_at, platform, earnings, created_at
		FROM referral_conversions
		WHERE user_id = $1
	`
	var conv ReferralConversion
	err := r.db.QueryRow(query, userID).Scan(
		&conv.ID, &conv.ReferralCode, &conv.UserID, &conv.RegisteredAt,
		&conv.Platform, &conv.Earnings, &conv.CreatedAt,
	)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &conv, nil
}

func (r *pgRepository) GetConversionsByCode(code string) ([]ReferralConversion, error) {
	query := `
		SELECT id, referral_code, user_id, registered_at, platform, earnings, created_at
		FROM referral_conversions
		WHERE referral_code = $1
		ORDER BY registered_at DESC
	`
	rows, err := r.db.Query(query, code)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var conversions []ReferralConversion
	for rows.Next() {
		var conv ReferralConversion
		err := rows.Scan(
			&conv.ID, &conv.ReferralCode, &conv.UserID, &conv.RegisteredAt,
			&conv.Platform, &conv.Earnings, &conv.CreatedAt,
		)
		if err != nil {
			return nil, err
		}
		conversions = append(conversions, conv)
	}
	return conversions, nil
}

// ========================================
// PREMIUM CONVERSIONS
// ========================================

func (r *pgRepository) CreatePremiumConversion(conv *ReferralPremiumConversion) error {
	query := `
		INSERT INTO referral_premium_conversions
		(referral_code, user_id, conversion_id, subscription_type, amount, commission, commission_percentage, converted_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
	`
	_, err := r.db.Exec(query, conv.ReferralCode, conv.UserID, conv.ConversionID, conv.SubscriptionType,
		conv.Amount, conv.Commission, conv.CommissionPercentage, conv.ConvertedAt)
	return err
}

func (r *pgRepository) GetPremiumConversionsByCode(code string) ([]ReferralPremiumConversion, error) {
	query := `
		SELECT id, referral_code, user_id, conversion_id, subscription_type, amount,
		       commission, commission_percentage, converted_at, created_at
		FROM referral_premium_conversions
		WHERE referral_code = $1
		ORDER BY converted_at DESC
	`
	rows, err := r.db.Query(query, code)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var conversions []ReferralPremiumConversion
	for rows.Next() {
		var conv ReferralPremiumConversion
		err := rows.Scan(
			&conv.ID, &conv.ReferralCode, &conv.UserID, &conv.ConversionID, &conv.SubscriptionType,
			&conv.Amount, &conv.Commission, &conv.CommissionPercentage, &conv.ConvertedAt, &conv.CreatedAt,
		)
		if err != nil {
			return nil, err
		}
		conversions = append(conversions, conv)
	}
	return conversions, nil
}

// ========================================
// METRICS - Individual Influencer
// ========================================

func (r *pgRepository) GetTotalRegistrationsByCode(code string) (int, error) {
	query := `SELECT COUNT(*) FROM referral_conversions WHERE referral_code = $1`
	var count int
	err := r.db.QueryRow(query, code).Scan(&count)
	return count, err
}

func (r *pgRepository) GetTotalPremiumByCode(code string) (int, error) {
	query := `SELECT COUNT(*) FROM referral_premium_conversions WHERE referral_code = $1`
	var count int
	err := r.db.QueryRow(query, code).Scan(&count)
	return count, err
}

func (r *pgRepository) GetTotalEarningsByCode(code string) (float64, error) {
	query := `
		SELECT COALESCE(SUM(rc.earnings), 0) + COALESCE(SUM(rpc.commission), 0) as total
		FROM referral_conversions rc
		LEFT JOIN referral_premium_conversions rpc ON rc.referral_code = rpc.referral_code
		WHERE rc.referral_code = $1
	`
	var total float64
	err := r.db.QueryRow(query, code).Scan(&total)
	return total, err
}

func (r *pgRepository) GetCurrentMonthRegistrationsByCode(code string) (int, error) {
	query := `
		SELECT COUNT(*) FROM referral_conversions
		WHERE referral_code = $1 AND registered_at >= DATE_TRUNC('month', NOW())
	`
	var count int
	err := r.db.QueryRow(query, code).Scan(&count)
	return count, err
}

func (r *pgRepository) GetCurrentMonthPremiumByCode(code string) (int, error) {
	query := `
		SELECT COUNT(*) FROM referral_premium_conversions
		WHERE referral_code = $1 AND converted_at >= DATE_TRUNC('month', NOW())
	`
	var count int
	err := r.db.QueryRow(query, code).Scan(&count)
	return count, err
}

func (r *pgRepository) GetCurrentMonthEarningsByCode(code string) (float64, error) {
	startOfMonth := time.Now().AddDate(0, 0, -time.Now().Day()+1).Format("2006-01-02 00:00:00")
	query := `
		SELECT
			COALESCE(SUM(rc.earnings), 0) + COALESCE(SUM(rpc.commission), 0) as total
		FROM referral_conversions rc
		LEFT JOIN referral_premium_conversions rpc ON rc.user_id = rpc.user_id
		WHERE rc.referral_code = $1
			AND (rc.registered_at >= $2 OR rpc.converted_at >= $2)
	`
	var total float64
	err := r.db.QueryRow(query, code, startOfMonth, startOfMonth).Scan(&total)
	return total, err
}

func (r *pgRepository) GetDailyMetrics(code string, days int) ([]DailyMetric, error) {
	query := `
		SELECT
			DATE(registered_at) as date,
			COUNT(*) as registrations,
			0 as premium_conversions,
			SUM(earnings) as earnings
		FROM referral_conversions
		WHERE referral_code = $1 AND registered_at >= NOW() - INTERVAL '1 day' * $2
		GROUP BY DATE(registered_at)
		ORDER BY date DESC
	`
	rows, err := r.db.Query(query, code, days)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	metricsMap := make(map[string]*DailyMetric)
	for rows.Next() {
		var date string
		var registrations int
		var premiumConversions int
		var earnings float64
		err := rows.Scan(&date, &registrations, &premiumConversions, &earnings)
		if err != nil {
			return nil, err
		}
		metricsMap[date] = &DailyMetric{
			Date:               date,
			Registrations:      registrations,
			PremiumConversions: premiumConversions,
			Earnings:           earnings,
		}
	}

	// Ahora obtener premium conversions
	query2 := `
		SELECT
			DATE(converted_at) as date,
			COUNT(*) as premium_conversions,
			SUM(commission) as commission
		FROM referral_premium_conversions
		WHERE referral_code = $1 AND converted_at >= NOW() - INTERVAL '1 day' * $2
		GROUP BY DATE(converted_at)
	`
	rows2, err := r.db.Query(query2, code, days)
	if err != nil {
		return nil, err
	}
	defer rows2.Close()

	for rows2.Next() {
		var date string
		var premiumConversions int
		var commission float64
		err := rows2.Scan(&date, &premiumConversions, &commission)
		if err != nil {
			return nil, err
		}
		if metric, exists := metricsMap[date]; exists {
			metric.PremiumConversions = premiumConversions
			metric.Earnings += commission
		} else {
			metricsMap[date] = &DailyMetric{
				Date:               date,
				Registrations:      0,
				PremiumConversions: premiumConversions,
				Earnings:           commission,
			}
		}
	}

	// Convertir map a slice
	var metrics []DailyMetric
	for _, metric := range metricsMap {
		metrics = append(metrics, *metric)
	}

	return metrics, nil
}

func (r *pgRepository) GetInfluencerRank(code string) (int, int, error) {
	// Primero, obtener todos los earnings ordenados
	query := `
		SELECT
			i.referral_code,
			COALESCE(SUM(rc.earnings), 0) + COALESCE(SUM(rpc.commission), 0) as total_earnings
		FROM influencers i
		LEFT JOIN referral_conversions rc ON i.referral_code = rc.referral_code
		LEFT JOIN referral_premium_conversions rpc ON i.referral_code = rpc.referral_code
		WHERE i.is_active = '1'
		GROUP BY i.referral_code
		ORDER BY total_earnings DESC
	`
	rows, err := r.db.Query(query)
	if err != nil {
		return 0, 0, err
	}
	defer rows.Close()

	rank := 0
	totalInfluencers := 0
	currentRank := 1

	for rows.Next() {
		var refCode string
		var earnings float64
		err := rows.Scan(&refCode, &earnings)
		if err != nil {
			return 0, 0, err
		}
		totalInfluencers++
		if refCode == code {
			rank = currentRank
		}
		currentRank++
	}

	return rank, totalInfluencers, nil
}

// ========================================
// METRICS - Aggregate
// ========================================

func (r *pgRepository) GetTotalReferrals() (int, error) {
	query := `SELECT COUNT(*) FROM referral_conversions`
	var count int
	err := r.db.QueryRow(query).Scan(&count)
	return count, err
}

func (r *pgRepository) GetTotalPremiumConversions() (int, error) {
	query := `SELECT COUNT(*) FROM referral_premium_conversions`
	var count int
	err := r.db.QueryRow(query).Scan(&count)
	return count, err
}

func (r *pgRepository) GetTotalEarnings() (float64, error) {
	query := `
		SELECT
			COALESCE(SUM(rc.earnings), 0) + COALESCE(SUM(rpc.commission), 0) as total
		FROM referral_conversions rc
		LEFT JOIN referral_premium_conversions rpc ON rc.user_id = rpc.user_id
	`
	var total float64
	err := r.db.QueryRow(query).Scan(&total)
	return total, err
}

func (r *pgRepository) GetActiveInfluencersCount() (int, error) {
	query := `SELECT COUNT(*) FROM influencers WHERE is_active = '1'`
	var count int
	err := r.db.QueryRow(query).Scan(&count)
	return count, err
}

func (r *pgRepository) GetTopPerformers(limit int) ([]TopPerformer, error) {
	query := `
		SELECT
			i.id as influencer_id,
			i.referral_code,
			i.name,
			i.platform,
			COUNT(DISTINCT rc.id) as total_registrations,
			COUNT(DISTINCT rpc.id) as total_premium_conversions,
			COALESCE(SUM(rc.earnings), 0) + COALESCE(SUM(rpc.commission), 0) as total_earnings
		FROM influencers i
		LEFT JOIN referral_conversions rc ON i.referral_code = rc.referral_code
		LEFT JOIN referral_premium_conversions rpc ON i.referral_code = rpc.referral_code
		WHERE i.is_active = '1'
		GROUP BY i.id, i.referral_code, i.name, i.platform
		ORDER BY total_earnings DESC
		LIMIT $1
	`
	rows, err := r.db.Query(query, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var performers []TopPerformer
	for rows.Next() {
		var p TopPerformer
		err := rows.Scan(
			&p.InfluencerID, &p.ReferralCode, &p.Name, &p.Platform,
			&p.TotalRegistrations, &p.TotalPremiumConversions, &p.TotalEarnings,
		)
		if err != nil {
			return nil, err
		}
		performers = append(performers, p)
	}
	return performers, nil
}

func (r *pgRepository) GetMonthlyTrend(months int) ([]MonthlyTrend, error) {
	query := `
		SELECT
			TO_CHAR(registered_at, 'YYYY-MM') as month,
			COUNT(*) as registrations,
			0 as premium,
			SUM(earnings) as earnings
		FROM referral_conversions
		WHERE registered_at >= NOW() - INTERVAL '1 month' * $1
		GROUP BY TO_CHAR(registered_at, 'YYYY-MM')
		ORDER BY month DESC
	`
	rows, err := r.db.Query(query, months)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	trendsMap := make(map[string]*MonthlyTrend)
	for rows.Next() {
		var month string
		var registrations, premium int
		var earnings float64
		err := rows.Scan(&month, &registrations, &premium, &earnings)
		if err != nil {
			return nil, err
		}
		trendsMap[month] = &MonthlyTrend{
			Month:         month,
			Registrations: registrations,
			Premium:       premium,
			Earnings:      earnings,
		}
	}

	// Obtener premium conversions por mes
	query2 := `
		SELECT
			TO_CHAR(converted_at, 'YYYY-MM') as month,
			COUNT(*) as premium,
			SUM(commission) as commission
		FROM referral_premium_conversions
		WHERE converted_at >= NOW() - INTERVAL '1 month' * $1
		GROUP BY TO_CHAR(converted_at, 'YYYY-MM')
	`
	rows2, err := r.db.Query(query2, months)
	if err != nil {
		return nil, err
	}
	defer rows2.Close()

	for rows2.Next() {
		var month string
		var premium int
		var commission float64
		err := rows2.Scan(&month, &premium, &commission)
		if err != nil {
			return nil, err
		}
		if trend, exists := trendsMap[month]; exists {
			trend.Premium = premium
			trend.Earnings += commission
		} else {
			trendsMap[month] = &MonthlyTrend{
				Month:         month,
				Registrations: 0,
				Premium:       premium,
				Earnings:      commission,
			}
		}
	}

	// Convertir map a slice
	var trends []MonthlyTrend
	for _, trend := range trendsMap {
		trends = append(trends, *trend)
	}

	return trends, nil
}

func (r *pgRepository) GetPlatformBreakdown() (map[string]PlatformBreakdown, error) {
	platforms := []string{"Instagram", "TikTok", "Reddit", "Other"}
	breakdown := make(map[string]PlatformBreakdown)

	for _, platform := range platforms {
		query := `
			SELECT
				COUNT(DISTINCT i.id) as influencer_count,
				COUNT(DISTINCT rc.id) as total_registrations,
				COUNT(DISTINCT rpc.id) as total_premium,
				COALESCE(SUM(rc.earnings), 0) + COALESCE(SUM(rpc.commission), 0) as total_earnings
			FROM influencers i
			LEFT JOIN referral_conversions rc ON i.referral_code = rc.referral_code
			LEFT JOIN referral_premium_conversions rpc ON i.referral_code = rpc.referral_code
			WHERE i.platform = $1 AND i.is_active = '1'
		`
		var pb PlatformBreakdown
		err := r.db.QueryRow(query, platform).Scan(
			&pb.Influencers, &pb.Registrations, &pb.Premium, &pb.Earnings,
		)
		if err != nil {
			return nil, fmt.Errorf("error querying platform %s: %w", platform, err)
		}
		breakdown[platform] = pb
	}

	return breakdown, nil
}
