package main

import (
	"alertly/internal/account"
	"alertly/internal/activate"
	"alertly/internal/achievements"
	"alertly/internal/analytics"
	"alertly/internal/auth"
	"alertly/internal/comments"
	"alertly/internal/common"

	// "alertly/internal/config" // No longer needed
	"alertly/internal/cronjob"
	"alertly/internal/database"
	"alertly/internal/editprofile"
	"alertly/internal/emails"
	"alertly/internal/feedback"
	"alertly/internal/getcategories"
	"alertly/internal/getclusterby"
	"alertly/internal/getclusterbyradius"
	"alertly/internal/getclustersbylocation"
	"alertly/internal/getincidentsasreels"
	"alertly/internal/getsubcategoriesbycategoryid"
	"alertly/internal/health"
	"alertly/internal/invitefriend"
	"alertly/internal/logging"
	"alertly/internal/media"
	"alertly/internal/middleware"
	"alertly/internal/myplaces"
	"alertly/internal/newincident"
	"alertly/internal/notifications"
	"alertly/internal/profile"
	"alertly/internal/referrals"
	"alertly/internal/reportincident"
	"alertly/internal/saveclusteraccount"
	"alertly/internal/scheduler"
	"alertly/internal/signup"
	"alertly/internal/tutorial"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	// "github.com/joho/godotenv" // No longer needed
)

func main() {
	log.Println("Starting Alertly Backend...")

	// --- AWS Lambda Refactor ---
	// Directly read configuration from environment variables.
	// This simplifies logic and removes dependency on .env files and complex config structs for the Lambda environment.
	dbUser := os.Getenv("DB_USER")
	dbPass := os.Getenv("DB_PASS")
	dbHost := os.Getenv("DB_HOST")
	dbPort := os.Getenv("DB_PORT")
	dbName := os.Getenv("DB_NAME")
	port := os.Getenv("PORT")

	if port == "" {
		port = "8080" // Default port if not set
	}

	// Inicializar el cliente de AWS SES
	emails.InitSES()

	// Configurar logging para producción
	logging.SetupProductionLogging()

	// Construir el DSN (Data Source Name) de forma robusta
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?parseTime=true",
		dbUser, dbPass, dbHost, dbPort, dbName)

	log.Printf("Connecting to database at %s:%s...", dbHost, dbPort)
	database.InitDB(dsn)
	log.Println("Database connected successfully")
	defer database.DB.Close()

	// Initialize database schema if needed
	if err := database.CheckAndInitDatabase(database.DB); err != nil {
		log.Printf("⚠️ Warning: Could not initialize database schema: %v", err)
		// Continue anyway, the app might still work with existing tables
	}

	// OPTIMIZACIÓN: Iniciar cache cleanup
	common.StartCacheCleanup()
	log.Println("Cache cleanup started")

	// PUSH NOTIFICATIONS: Iniciar cronjobs internos
	scheduler.StartCronjobs()
	log.Println("Cronjob scheduler started")

	router := gin.Default()

	// PRODUCCIÓN: Configurar middlewares de seguridad
	// The GIN_MODE check is kept to ensure these are only applied in a production-like environment.
	if os.Getenv("GIN_MODE") == "release" {
		// Logging optimizado para producción
		router.Use(logging.ProductionLogger())

		// Security headers para producción
		router.Use(middleware.SecurityHeadersMiddleware())

		// Rate limit headers informativos
		router.Use(middleware.RateLimitHeadersMiddleware())
	}

	// CORS optimizado con SafeCORSMiddleware (O(1) lookup, misma funcionalidad)
	router.Use(middleware.SafeCORSMiddleware())

	// IMÁGENES AHORA SE SIRVEN DESDE S3
	// Ya no necesitamos servir archivos estáticos locales
	log.Println("Images served from S3: alertly-images-production")

	// HEALTH CHECKS: Endpoints de monitoreo (sin rate limiting)
	router.GET("/health", health.HealthHandler(database.DB))
	router.GET("/health/ready", health.ReadinessHandler(database.DB))
	router.GET("/health/live", health.LivenessHandler())
	log.Println("Health check endpoints configured")

	// OPTIMIZACIÓN: Rate limiting para endpoints públicos
	router.Use(middleware.RateLimitMiddleware())

	router.POST("/account/signup", signup.RegisterUserHandler)
	router.POST("/account/signin", auth.SignIn)
	router.POST("/account/activate", activate.ActivateAccount)

	api := router.Group("/api")
	api.Use(middleware.TokenAuthMiddleware())
	// OPTIMIZACIÓN: Rate limiting más estricto para endpoints autenticados
	api.Use(middleware.RateLimitMiddlewareStrict())

	api.GET("/account/validate", auth.ValidateSession)
	router.GET("/category/get_all", getcategories.GetCategories)
	router.GET("/category/getsubcategoriesbycategoryid/:id", getsubcategoriesbycategoryid.GetSubcategoriesByCategoryId)

	// PÚBLICO: Endpoint para landing pages web (con rate limiting estricto)
	publicRoutes := router.Group("/public")
	publicRoutes.Use(middleware.RateLimitMiddlewarePublic()) // Rate limiting más estricto
	publicRoutes.GET("/cluster/getbyid/:incl_id", getclusterby.ViewPublic)

	api.POST("/incident/create", newincident.Create)
	api.GET("/cluster/getbyid/:incl_id", getclusterby.View)
	router.GET("/cluster/getbylocation/:min_latitude/:max_latitude/:min_longitude/:max_longitude/:from_date/:to_date/:insu_id", getclustersbylocation.Get)
	router.GET("/cluster/getbyradius/:latitude/:longitude/:radius/:from_date/:to_date/:insu_id", getclusterbyradius.GetByRadius)
	api.GET("/cluster/getasreel/:min_latitude/:max_latitude/:min_longitude/:max_longitude", getincidentsasreels.GetReel)

	router.GET("/account/myplaces/get/:account_id", myplaces.Get)
	api.GET("/account/get_my_profile", editprofile.GetMyProfile)
	api.POST("/account/edit_fullname", editprofile.UpdateFullName)
	api.POST("/account/edit/nickname", editprofile.UpdateNickname)
	api.POST("/account/edit/birthdate", editprofile.UpdateBirthDate)
	api.POST("account/edit/desactivate_account", editprofile.DesactivateAccount)
	api.POST("account/edit/receive_notifications", editprofile.UpdateReceiveNotifications)
	api.POST("/account/edit/email", editprofile.UpdateEmail)
	api.POST("/account/edit/password", editprofile.UpdatePassword)
	api.POST("account/edit/picture", editprofile.UpdateThumbnail)
	api.GET("/account/get_history", account.GetHistory)
	api.GET("/account/get_viewed_incident_ids", account.GetViewedIncidentIds)
	api.GET("account/get_counter_histories", account.GetCounterHistories)
	api.POST("/account/clear_history", account.ClearHistory)
	api.POST("/account/delete_account", account.DeleteAccount)
	api.POST("/account/check_password", auth.CheckPasswordMatch)

	// Premium-protected: Multiple alert locations feature (Saved Places)
	premiumMW := middleware.PremiumMiddleware(database.DB)
	api.POST("/account/myplaces/add", premiumMW, myplaces.Add)
	api.GET("/account/myplaces/get", premiumMW, myplaces.GetByAccountId)
	api.GET("/account/myplaces/get_by_id/:afl_id", premiumMW, myplaces.GetById)
	api.POST("/account/myplaces/update", premiumMW, myplaces.Update)
	api.POST("account/set_has_finished_tutorial", account.SetHasFinishedTutorial)
	api.POST("/account/myplaces/full_update", premiumMW, myplaces.FullUpdate)
	api.GET("/account/myplaces/delete/:afl_id", premiumMW, myplaces.Delete)
	api.GET("/account/profile/get_by_id/:account_id", profile.GetById)
	api.GET("/account/cluster/toggle_save/:incl_id", saveclusteraccount.ToggleSaveClusterAccount)
	api.POST("/cluster/send_comment", middleware.ProfanityFilterMiddleware(), comments.SaveClusterComment)
	api.GET("/saved/get_my_list", saveclusteraccount.GetMyList)
	api.GET("/saved/delete/:acs_id", saveclusteraccount.DeleteFollowIncident)
	api.POST("/account/report/:account_id", profile.ReportAccount)
	api.GET("/account/get_my_info", account.GetMyInfo)
	api.POST("/purchase/apple/validate", account.ValidateAppleReceipt)
	api.POST("/account/update_premium_status", account.UpdatePremiumStatusHandler)
	api.POST("/send_feedback", feedback.SendFeedback)
	api.POST("/send_invitation", invitefriend.Save)
	api.POST("report_incident", reportincident.ReportIncident)

	// Tutorial
	api.POST("/tutorial/complete", tutorial.CompleteHandler)

	// Analytics endpoints
	analyticsService := analytics.NewBasicAnalytics(database.DB)
	analyticsHandler := analytics.NewHandler(analyticsService)
	api.GET("/analytics/summary", analyticsHandler.GetAnalyticsSummary)
	api.GET("/analytics/predictions", analyticsHandler.GetSimplePredictions)
	api.GET("/analytics/test", analyticsHandler.TestAnalytics)

	// Achievements endpoints (badges, ranks, citizen score)
	achievementsRepo := achievements.NewRepository(database.DB)
	achievementsService := achievements.NewService(achievementsRepo)
	achievementsHandler := achievements.NewHandler(achievementsService)
	api.GET("/achievements/pending", achievementsHandler.GetPending)
	api.PUT("/achievements/:id/mark-shown", achievementsHandler.MarkAsShown)

	// ==================================================
	// REFERRAL SYSTEM ENDPOINTS
	// ==================================================
	// Inicializar el sistema de referrals
	referralsRepo := referrals.NewRepository(database.DB)
	referralsService := referrals.NewService(referralsRepo)
	referralsHandler := referrals.NewHandler(referralsService)

	// Grupo de endpoints de referrals v1
	referralV1 := router.Group("/api/v1")
	{
		// ENDPOINT 1: Validar código de referral (PÚBLICO - sin autenticación)
		// Usado por la app móvil durante el signup
		referralV1.POST("/referral/validate", referralsHandler.ValidateReferralCode)

		// Endpoints protegidos con API Key (backend web los consume)
		referralV1Protected := referralV1.Group("")
		referralV1Protected.Use(middleware.ReferralAPIKeyMiddleware())
		{
			// ENDPOINT 2: Registrar conversión de registro
			referralV1Protected.POST("/referral/conversion", referralsHandler.RegisterConversion)

			// ENDPOINT 3: Registrar conversión premium
			referralV1Protected.POST("/referral/premium-conversion", referralsHandler.RegisterPremiumConversion)

			// ENDPOINT 4: Obtener métricas de influencer individual
			referralV1Protected.GET("/referrals/metrics", referralsHandler.GetInfluencerMetrics)

			// ENDPOINT 5: Obtener métricas agregadas
			referralV1Protected.GET("/referrals/aggregate", referralsHandler.GetAggregateMetrics)

			// ENDPOINT 6: Sincronizar influencer desde backend web
			referralV1Protected.POST("/referral/sync-influencer", referralsHandler.SyncInfluencer)
		}
	}
	log.Println("✅ Referral system endpoints registered")

	// Cronjob endpoints (for manual execution and monitoring)
	api.POST("/cronjob/premium/expire", cronjob.RunPremiumExpirationCheck)
	api.GET("/cronjob/premium/stats", cronjob.GetPremiumStats)
	api.POST("/cronjob/premium/warnings", cronjob.SendExpirationWarnings)
	api.GET("/analytics/location", analyticsHandler.GetLocationAnalytics)

	// comunitacions with apple APN (to send push notifications)
	// MOVIDO: Endpoints de device tokens sin rate limiting estricto
	router.POST("/api/device_tokens", middleware.TokenAuthMiddleware(), notifications.SaveDeviceToken)
	router.DELETE("/api/device_tokens", middleware.TokenAuthMiddleware(), notifications.DeleteDeviceToken)

	// Notification endpoints
	api.GET("/notifications", notifications.GetNotifications)
	api.GET("/notifications/unread_count", notifications.GetUnreadCount)
	api.POST("/notifications/mark_as_read", notifications.MarkAsRead)
	api.POST("/notifications/mark_all_as_read", notifications.MarkAllAsRead)
	api.DELETE("/notifications", notifications.DeleteNotification)

	// TESTING
	api.POST("/test_push", notifications.TestPushHandler)

	// Premium middleware test endpoint - Use this to verify premium validation works
	api.GET("/test_premium", middleware.PremiumMiddleware(database.DB), func(c *gin.Context) {
		accountID, _ := c.Get("AccountId")
		c.JSON(http.StatusOK, gin.H{
			"message":    "✅ Premium validation successful",
			"account_id": accountID,
			"is_premium": true,
		})
	})

	// MEDIA PROCESSING: Endpoints para procesamiento de imágenes con pixelado
	api.POST("/media/reprocess/:incl_id", media.ReprocessImageHandler)
	api.POST("/media/test_pixelation", media.TestPixelationHandler)

	// PRODUCCIÓN: Iniciar servidor con configuración segura
	serverPort := ":" + port
	log.Printf("Alertly Backend starting on port %s", port)
	log.Printf("Environment: %s", os.Getenv("GIN_MODE"))
	log.Printf("Health check: http://localhost%s/health", serverPort)

	log.Println("Starting Gin router...")
	if err := router.Run(serverPort); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}

// getEnvOrDefault is no longer needed as we read variables directly.
