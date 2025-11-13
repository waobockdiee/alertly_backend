package referrals

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// Handler maneja las peticiones HTTP para el sistema de referrals
type Handler struct {
	service Service
}

// NewHandler crea una nueva instancia del handler
func NewHandler(service Service) *Handler {
	return &Handler{service: service}
}

// ========================================
// ENDPOINT 1: POST /api/v1/referral/validate
// PÚBLICO - Sin autenticación
// ========================================

// ValidateReferralCode valida un código de referral durante el signup
// @Summary Validate referral code
// @Description Valida si un código de referral es válido y retorna información del influencer
// @Tags Referrals
// @Accept json
// @Produce json
// @Param body body ValidateReferralCodeRequest true "Referral code to validate"
// @Success 200 {object} ValidateReferralCodeResponse "Código válido"
// @Success 404 {object} ValidateReferralCodeResponse "Código inválido"
// @Failure 400 {object} map[string]string "Bad request"
// @Router /api/v1/referral/validate [post]
func (h *Handler) ValidateReferralCode(c *gin.Context) {
	var req ValidateReferralCodeRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request format", "details": err.Error()})
		return
	}

	response, err := h.service.ValidateReferralCode(req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error", "details": err.Error()})
		return
	}

	// Si el código no es válido, retornar 404
	if !response.Valid {
		c.JSON(http.StatusNotFound, response)
		return
	}

	// Código válido
	c.JSON(http.StatusOK, response)
}

// ========================================
// ENDPOINT 2: POST /api/v1/referral/conversion
// PROTEGIDO - Requiere API Key
// ========================================

// RegisterConversion registra una conversión de registro con código de referral
// @Summary Register a signup conversion
// @Description Registra cuando un usuario se registra usando un código de referral
// @Tags Referrals
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param body body RegisterConversionRequest true "Conversion data"
// @Success 201 {object} RegisterConversionResponse "Conversión registrada"
// @Failure 400 {object} RegisterConversionResponse "Bad request o código inválido"
// @Failure 401 {object} map[string]string "Unauthorized"
// @Router /api/v1/referral/conversion [post]
func (h *Handler) RegisterConversion(c *gin.Context) {
	var req RegisterConversionRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request format", "details": err.Error()})
		return
	}

	response, err := h.service.RegisterConversion(req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error", "details": err.Error()})
		return
	}

	// Si hubo un error de validación, retornar 400
	if !response.Success {
		c.JSON(http.StatusBadRequest, response)
		return
	}

	// Conversión registrada exitosamente
	c.JSON(http.StatusCreated, response)
}

// ========================================
// ENDPOINT 3: POST /api/v1/referral/premium-conversion
// PROTEGIDO - Requiere API Key
// ========================================

// RegisterPremiumConversion registra una conversión premium de usuario referido
// @Summary Register a premium conversion
// @Description Registra cuando un usuario referido se suscribe a premium
// @Tags Referrals
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param body body RegisterPremiumConversionRequest true "Premium conversion data"
// @Success 201 {object} RegisterPremiumConversionResponse "Conversión premium registrada"
// @Failure 400 {object} RegisterPremiumConversionResponse "Bad request o usuario sin código"
// @Failure 401 {object} map[string]string "Unauthorized"
// @Router /api/v1/referral/premium-conversion [post]
func (h *Handler) RegisterPremiumConversion(c *gin.Context) {
	var req RegisterPremiumConversionRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request format", "details": err.Error()})
		return
	}

	response, err := h.service.RegisterPremiumConversion(req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error", "details": err.Error()})
		return
	}

	// Si hubo un error de validación, retornar 400
	if !response.Success {
		c.JSON(http.StatusBadRequest, response)
		return
	}

	// Conversión premium registrada exitosamente
	c.JSON(http.StatusCreated, response)
}

// ========================================
// ENDPOINT 4: GET /api/v1/referrals/metrics?code={CODE}
// PROTEGIDO - Requiere API Key
// ========================================

// GetInfluencerMetrics obtiene métricas de un influencer específico
// @Summary Get influencer metrics
// @Description Obtiene métricas detalladas de un influencer por código de referral
// @Tags Referrals
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param code query string true "Referral code"
// @Success 200 {object} InfluencerMetricsResponse "Métricas del influencer"
// @Failure 400 {object} map[string]string "Missing code parameter"
// @Failure 404 {object} map[string]string "Referral code not found"
// @Failure 401 {object} map[string]string "Unauthorized"
// @Router /api/v1/referrals/metrics [get]
func (h *Handler) GetInfluencerMetrics(c *gin.Context) {
	code := c.Query("code")
	if code == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Missing required query parameter: code"})
		return
	}

	response, err := h.service.GetInfluencerMetrics(code)
	if err != nil {
		// Si el error es "not found", retornar 404
		if err.Error() == "referral code not found" {
			c.JSON(http.StatusNotFound, gin.H{"error": "Referral code not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error", "details": err.Error()})
		return
	}

	c.JSON(http.StatusOK, response)
}

// ========================================
// ENDPOINT 5: GET /api/v1/referrals/aggregate
// PROTEGIDO - Requiere API Key
// ========================================

// GetAggregateMetrics obtiene métricas agregadas de todos los influencers
// @Summary Get aggregate metrics
// @Description Obtiene métricas agregadas del sistema completo de referrals
// @Tags Referrals
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Success 200 {object} AggregateMetricsResponse "Métricas agregadas"
// @Failure 401 {object} map[string]string "Unauthorized"
// @Router /api/v1/referrals/aggregate [get]
func (h *Handler) GetAggregateMetrics(c *gin.Context) {
	response, err := h.service.GetAggregateMetrics()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error", "details": err.Error()})
		return
	}

	c.JSON(http.StatusOK, response)
}

// ========================================
// ENDPOINT 6: POST /api/v1/referral/sync-influencer
// PROTEGIDO - Requiere API Key
// ========================================

// SyncInfluencer sincroniza un influencer desde el backend web
// @Summary Sync influencer from web backend
// @Description Crea o actualiza un influencer en la base de datos local
// @Tags Referrals
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param body body SyncInfluencerRequest true "Influencer data"
// @Success 200 {object} SyncInfluencerResponse "Influencer sincronizado"
// @Failure 400 {object} map[string]string "Bad request"
// @Failure 401 {object} map[string]string "Unauthorized"
// @Router /api/v1/referral/sync-influencer [post]
func (h *Handler) SyncInfluencer(c *gin.Context) {
	var req SyncInfluencerRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request format", "details": err.Error()})
		return
	}

	response, err := h.service.SyncInfluencer(req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error", "details": err.Error()})
		return
	}

	// Si hubo un error en el servicio, retornar 400
	if !response.Success {
		c.JSON(http.StatusBadRequest, response)
		return
	}

	c.JSON(http.StatusOK, response)
}
