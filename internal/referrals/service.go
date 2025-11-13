package referrals

import (
	"fmt"
	"math"
	"time"
)

// Service define la lógica de negocio para el sistema de referrals
type Service interface {
	ValidateReferralCode(req ValidateReferralCodeRequest) (ValidateReferralCodeResponse, error)
	RegisterConversion(req RegisterConversionRequest) (RegisterConversionResponse, error)
	RegisterPremiumConversion(req RegisterPremiumConversionRequest) (RegisterPremiumConversionResponse, error)
	GetInfluencerMetrics(code string) (InfluencerMetricsResponse, error)
	GetAggregateMetrics() (AggregateMetricsResponse, error)
	SyncInfluencer(req SyncInfluencerRequest) (SyncInfluencerResponse, error)
}

type service struct {
	repo Repository
}

// NewService crea una nueva instancia del servicio
func NewService(repo Repository) Service {
	return &service{repo: repo}
}

// ========================================
// ENDPOINT 1: Validate Referral Code
// ========================================

func (s *service) ValidateReferralCode(req ValidateReferralCodeRequest) (ValidateReferralCodeResponse, error) {
	// Buscar el influencer por código
	influencer, err := s.repo.GetInfluencerByCode(req.ReferralCode)
	if err != nil {
		return ValidateReferralCodeResponse{}, fmt.Errorf("error querying influencer: %w", err)
	}

	// Si no existe o no está activo
	if influencer == nil || !influencer.IsActive {
		return ValidateReferralCodeResponse{
			Valid:   false,
			Message: "Invalid referral code",
		}, nil
	}

	// Código válido
	return ValidateReferralCodeResponse{
		Valid:            true,
		InfluencerID:     &influencer.ID,
		InfluencerName:   influencer.Name,
		PremiumTrialDays: 14, // Beneficio por usar código de referral
	}, nil
}

// ========================================
// ENDPOINT 2: Register Conversion
// ========================================

func (s *service) RegisterConversion(req RegisterConversionRequest) (RegisterConversionResponse, error) {
	// 1. Validar que el código sea válido
	influencer, err := s.repo.GetInfluencerByCode(req.ReferralCode)
	if err != nil {
		return RegisterConversionResponse{}, fmt.Errorf("error querying influencer: %w", err)
	}
	if influencer == nil || !influencer.IsActive {
		return RegisterConversionResponse{
			Success: false,
			Error:   "Invalid referral code",
		}, nil
	}

	// 2. Verificar que el usuario NO haya usado otro código antes
	existingConversion, err := s.repo.GetConversionByUserID(req.UserID)
	if err != nil {
		return RegisterConversionResponse{}, fmt.Errorf("error checking existing conversion: %w", err)
	}
	if existingConversion != nil {
		return RegisterConversionResponse{
			Success: false,
			Error:   "User already registered with a referral code",
		}, nil
	}

	// 3. Parsear timestamp
	registeredAt, err := time.Parse(time.RFC3339, req.RegisteredAt)
	if err != nil {
		return RegisterConversionResponse{
			Success: false,
			Error:   "Invalid timestamp format. Expected ISO 8601 (RFC3339)",
		}, nil
	}

	// 4. Crear la conversión
	conversion := &ReferralConversion{
		ReferralCode: req.ReferralCode,
		UserID:       req.UserID,
		RegisteredAt: registeredAt,
		Platform:     req.Platform,
		Earnings:     0.10, // $0.10 CAD por registro
	}

	err = s.repo.CreateConversion(conversion)
	if err != nil {
		return RegisterConversionResponse{}, fmt.Errorf("error creating conversion: %w", err)
	}

	return RegisterConversionResponse{
		Success:                 true,
		Message:                 "Conversion recorded",
		InfluencerEarningsAdded: 0.10,
	}, nil
}

// ========================================
// ENDPOINT 3: Register Premium Conversion
// ========================================

func (s *service) RegisterPremiumConversion(req RegisterPremiumConversionRequest) (RegisterPremiumConversionResponse, error) {
	// 1. Buscar el referral_code asociado a este user_id
	conversion, err := s.repo.GetConversionByUserID(req.UserID)
	if err != nil {
		return RegisterPremiumConversionResponse{}, fmt.Errorf("error querying conversion: %w", err)
	}
	if conversion == nil {
		return RegisterPremiumConversionResponse{
			Success: false,
			Error:   "User not found or no referral code associated",
		}, nil
	}

	referralCode := conversion.ReferralCode
	conversionID := conversion.ID

	// Si se proporcionó un referral_code en el request, verificar que coincida
	if req.ReferralCode != "" && req.ReferralCode != referralCode {
		return RegisterPremiumConversionResponse{
			Success: false,
			Error:   "Provided referral code does not match user's conversion",
		}, nil
	}

	// 2. Parsear timestamp
	convertedAt, err := time.Parse(time.RFC3339, req.ConvertedAt)
	if err != nil {
		return RegisterPremiumConversionResponse{
			Success: false,
			Error:   "Invalid timestamp format. Expected ISO 8601 (RFC3339)",
		}, nil
	}

	// 3. Calcular comisión (15% del monto)
	commissionPercentage := 15.00
	commission := math.Round(req.Amount*(commissionPercentage/100)*100) / 100 // Redondear a 2 decimales

	// 4. Crear la conversión premium
	premiumConversion := &ReferralPremiumConversion{
		ReferralCode:         referralCode,
		UserID:               req.UserID,
		ConversionID:         &conversionID,
		SubscriptionType:     req.SubscriptionType,
		Amount:               req.Amount,
		Commission:           commission,
		CommissionPercentage: commissionPercentage,
		ConvertedAt:          convertedAt,
	}

	err = s.repo.CreatePremiumConversion(premiumConversion)
	if err != nil {
		return RegisterPremiumConversionResponse{}, fmt.Errorf("error creating premium conversion: %w", err)
	}

	return RegisterPremiumConversionResponse{
		Success:              true,
		Message:              "Premium conversion recorded",
		InfluencerCommission: commission,
		CommissionPercentage: commissionPercentage,
	}, nil
}

// ========================================
// ENDPOINT 4: Get Influencer Metrics
// ========================================

func (s *service) GetInfluencerMetrics(code string) (InfluencerMetricsResponse, error) {
	// 1. Validar que el código exista
	influencer, err := s.repo.GetInfluencerByCode(code)
	if err != nil {
		return InfluencerMetricsResponse{}, fmt.Errorf("error querying influencer: %w", err)
	}
	if influencer == nil {
		return InfluencerMetricsResponse{}, fmt.Errorf("referral code not found")
	}

	// 2. Obtener totales históricos
	totalRegistrations, err := s.repo.GetTotalRegistrationsByCode(code)
	if err != nil {
		return InfluencerMetricsResponse{}, fmt.Errorf("error getting total registrations: %w", err)
	}

	totalPremium, err := s.repo.GetTotalPremiumByCode(code)
	if err != nil {
		return InfluencerMetricsResponse{}, fmt.Errorf("error getting total premium: %w", err)
	}

	totalEarnings, err := s.repo.GetTotalEarningsByCode(code)
	if err != nil {
		return InfluencerMetricsResponse{}, fmt.Errorf("error getting total earnings: %w", err)
	}

	// 3. Obtener métricas del mes actual
	currentMonthRegistrations, err := s.repo.GetCurrentMonthRegistrationsByCode(code)
	if err != nil {
		return InfluencerMetricsResponse{}, fmt.Errorf("error getting current month registrations: %w", err)
	}

	currentMonthPremium, err := s.repo.GetCurrentMonthPremiumByCode(code)
	if err != nil {
		return InfluencerMetricsResponse{}, fmt.Errorf("error getting current month premium: %w", err)
	}

	currentMonthEarnings, err := s.repo.GetCurrentMonthEarningsByCode(code)
	if err != nil {
		return InfluencerMetricsResponse{}, fmt.Errorf("error getting current month earnings: %w", err)
	}

	// 4. Proyectar earnings del mes
	now := time.Now()
	daysElapsed := now.Day()
	daysInMonth := time.Date(now.Year(), now.Month()+1, 0, 0, 0, 0, 0, time.UTC).Day()

	var projectedEarnings float64
	if daysElapsed > 0 {
		projectedEarnings = (currentMonthEarnings / float64(daysElapsed)) * float64(daysInMonth)
		projectedEarnings = math.Round(projectedEarnings*100) / 100
	}

	// 5. Obtener métricas diarias (últimos 30 días)
	dailyMetrics, err := s.repo.GetDailyMetrics(code, 30)
	if err != nil {
		return InfluencerMetricsResponse{}, fmt.Errorf("error getting daily metrics: %w", err)
	}

	// 6. Calcular ranking
	rank, totalInfluencers, err := s.repo.GetInfluencerRank(code)
	if err != nil {
		return InfluencerMetricsResponse{}, fmt.Errorf("error getting rank: %w", err)
	}

	// 7. Construir respuesta
	return InfluencerMetricsResponse{
		ReferralCode:              code,
		InfluencerID:              influencer.ID,
		TotalRegistrations:        totalRegistrations,
		TotalPremiumConversions:   totalPremium,
		TotalEarnings:             math.Round(totalEarnings*100) / 100,
		CurrentMonthRegistrations: currentMonthRegistrations,
		CurrentMonthPremium:       currentMonthPremium,
		CurrentMonthEarnings:      math.Round(currentMonthEarnings*100) / 100,
		ProjectedMonthEarnings:    projectedEarnings,
		DailyMetrics:              dailyMetrics,
		Rank:                      rank,
		TotalInfluencers:          totalInfluencers,
	}, nil
}

// ========================================
// ENDPOINT 5: Get Aggregate Metrics
// ========================================

func (s *service) GetAggregateMetrics() (AggregateMetricsResponse, error) {
	// 1. Totales generales
	totalReferrals, err := s.repo.GetTotalReferrals()
	if err != nil {
		return AggregateMetricsResponse{}, fmt.Errorf("error getting total referrals: %w", err)
	}

	totalPremium, err := s.repo.GetTotalPremiumConversions()
	if err != nil {
		return AggregateMetricsResponse{}, fmt.Errorf("error getting total premium: %w", err)
	}

	totalEarnings, err := s.repo.GetTotalEarnings()
	if err != nil {
		return AggregateMetricsResponse{}, fmt.Errorf("error getting total earnings: %w", err)
	}

	// 2. Contar influencers activos
	activeInfluencers, err := s.repo.GetActiveInfluencersCount()
	if err != nil {
		return AggregateMetricsResponse{}, fmt.Errorf("error getting active influencers: %w", err)
	}

	// 3. Calcular tasa de conversión
	var conversionRate float64
	if totalReferrals > 0 {
		conversionRate = (float64(totalPremium) / float64(totalReferrals)) * 100
		conversionRate = math.Round(conversionRate*100) / 100
	}

	// 4. Top performers (Top 10)
	topPerformers, err := s.repo.GetTopPerformers(10)
	if err != nil {
		return AggregateMetricsResponse{}, fmt.Errorf("error getting top performers: %w", err)
	}

	// 5. Monthly trend (últimos 12 meses)
	monthlyTrend, err := s.repo.GetMonthlyTrend(12)
	if err != nil {
		return AggregateMetricsResponse{}, fmt.Errorf("error getting monthly trend: %w", err)
	}

	// 6. Platform breakdown
	platformBreakdown, err := s.repo.GetPlatformBreakdown()
	if err != nil {
		return AggregateMetricsResponse{}, fmt.Errorf("error getting platform breakdown: %w", err)
	}

	// 7. Construir respuesta
	return AggregateMetricsResponse{
		TotalReferrals:          totalReferrals,
		TotalPremiumConversions: totalPremium,
		TotalEarningsPaid:       math.Round(totalEarnings*100) / 100,
		ActiveInfluencers:       activeInfluencers,
		ConversionRate:          conversionRate,
		TopPerformers:           topPerformers,
		MonthlyTrend:            monthlyTrend,
		PlatformBreakdown:       platformBreakdown,
	}, nil
}

// ========================================
// ENDPOINT 6: Sync Influencer (BONUS)
// ========================================

func (s *service) SyncInfluencer(req SyncInfluencerRequest) (SyncInfluencerResponse, error) {
	// Crear o actualizar influencer
	influencer := &Influencer{
		WebInfluencerID: req.WebInfluencerID,
		ReferralCode:    req.ReferralCode,
		Name:            req.Name,
		Platform:        req.Platform,
		IsActive:        req.IsActive,
	}

	err := s.repo.UpsertInfluencer(influencer)
	if err != nil {
		return SyncInfluencerResponse{
			Success: false,
			Error:   fmt.Sprintf("Error syncing influencer: %v", err),
		}, nil
	}

	return SyncInfluencerResponse{
		Success: true,
		Message: fmt.Sprintf("Influencer %s synced successfully", req.ReferralCode),
	}, nil
}
