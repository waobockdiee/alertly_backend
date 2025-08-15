package main

import (
	"alertly/internal/account"
	"alertly/internal/activate"
	"alertly/internal/analytics"
	"alertly/internal/auth"
	"alertly/internal/comments"
	"alertly/internal/common"
	"alertly/internal/config"
	"alertly/internal/cronjob"
	"alertly/internal/database"
	"alertly/internal/editprofile"
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
	"alertly/internal/middleware"
	"alertly/internal/myplaces"
	"alertly/internal/newincident"
	"alertly/internal/notifications"
	"alertly/internal/profile"
	"alertly/internal/reportincident"
	"alertly/internal/saveclusteraccount"
	"alertly/internal/signup"
	"alertly/internal/tutorial"
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	log.Println("🚀 Starting Alertly Backend...")

	// ✅ PRODUCCIÓN: Configuración segura desde variables de ambiente
	var cfg *config.ProductionConfig

	if os.Getenv("GIN_MODE") == "release" {
		// ✅ Modo producción: Solo variables de ambiente
		log.Println("📦 Loading production configuration from environment variables...")
		cfg = config.LoadProductionConfig()
	} else {
		// ✅ Modo desarrollo: Mantener compatibilidad con .env
		log.Println("🔧 Development mode: Loading from .env file...")
		var err error
		if os.Getenv("NODE_ENV") == "production" {
			err = godotenv.Load(".env.production")
		} else {
			err = godotenv.Load(".env")
		}

		if err != nil {
			log.Printf("⚠️ Warning: .env file not found, using environment variables: %v", err)
		}

		// Crear configuración desde variables de ambiente (compatible con .env)
		cfg = &config.ProductionConfig{
			DBUser: getEnvOrDefault("DB_USER", ""),
			DBPass: getEnvOrDefault("DB_PASS", ""),
			DBHost: getEnvOrDefault("DB_HOST", "localhost"),
			DBPort: getEnvOrDefault("DB_PORT", "3306"),
			DBName: getEnvOrDefault("DB_NAME", ""),
			Port:   getEnvOrDefault("PORT", "8080"),
		}
	}

	// ✅ Configurar logging para producción
	logging.SetupProductionLogging()

	// ✅ Configurar base de datos con la nueva configuración
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?parseTime=true",
		cfg.DBUser, cfg.DBPass, cfg.DBHost, cfg.DBPort, cfg.DBName)

	log.Printf("🗄️ Connecting to database at %s:%s...", cfg.DBHost, cfg.DBPort)
	database.InitDB(dsn)
	log.Println("✅ Database connected successfully")
	defer database.DB.Close()

	// ✅ OPTIMIZACIÓN: Iniciar cache cleanup
	common.StartCacheCleanup()
	log.Println("✅ Cache cleanup started")

	router := gin.Default()

	// ✅ PRODUCCIÓN: Configurar middlewares de seguridad
	if os.Getenv("GIN_MODE") == "release" {
		// ✅ Logging optimizado para producción
		router.Use(logging.ProductionLogger())

		// ✅ Security headers para producción
		router.Use(middleware.SecurityHeadersMiddleware())

		// ✅ Rate limit headers informativos
		router.Use(middleware.RateLimitHeadersMiddleware())
	}

	// ✅ CORS optimizado con SafeCORSMiddleware (O(1) lookup, misma funcionalidad)
	router.Use(middleware.SafeCORSMiddleware())

	wd, err := os.Getwd()
	if err != nil {
		log.Fatalf("Failed to get working directory: %v", err)
	}

	uploadsPath := filepath.Join(wd, "uploads")
	log.Println("Serving uploads from:", uploadsPath)
	// ✅ CORRECCIÓN: Usar ruta absoluta para servir archivos estáticos
	router.Static("/uploads", uploadsPath)

	// ✅ HEALTH CHECKS: Endpoints de monitoreo (sin rate limiting)
	router.GET("/health", health.HealthHandler(database.DB))
	router.GET("/health/ready", health.ReadinessHandler(database.DB))
	router.GET("/health/live", health.LivenessHandler())
	log.Println("✅ Health check endpoints configured")

	// ✅ OPTIMIZACIÓN: Rate limiting para endpoints públicos
	router.Use(middleware.RateLimitMiddleware())

	router.POST("/account/signup", signup.RegisterUserHandler)
	router.POST("/account/signin", auth.SignIn)
	router.POST("/account/activate", activate.ActivateAccount)

	api := router.Group("/api")
	api.Use(middleware.TokenAuthMiddleware())
	// ✅ OPTIMIZACIÓN: Rate limiting más estricto para endpoints autenticados
	api.Use(middleware.RateLimitMiddlewareStrict())

	api.GET("/account/validate", auth.ValidateSession)
	router.GET("/category/get_all", getcategories.GetCategories)
	router.GET("/category/getsubcategoriesbycategoryid/:id", getsubcategoriesbycategoryid.GetSubcategoriesByCategoryId)

	// ✅ PÚBLICO: Endpoint para landing pages web (con rate limiting estricto)
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
	api.POST("/account/myplaces/add", myplaces.Add)
	api.GET("/account/myplaces/get", myplaces.GetByAccountId)
	api.GET("/account/myplaces/get_by_id/:afl_id", myplaces.GetById)
	api.POST("/account/myplaces/update", myplaces.Update)
	api.POST("account/set_has_finished_tutorial", account.SetHasFinishedTutorial)
	api.POST("/account/myplaces/full_update", myplaces.FullUpdate)
	api.GET("/account/myplaces/delete/:afl_id", myplaces.Delete)
	api.GET("/account/profile/get_by_id/:account_id", profile.GetById)
	api.GET("/account/cluster/toggle_save/:incl_id", saveclusteraccount.ToggleSaveClusterAccount)
	api.POST("/cluster/send_comment", middleware.ProfanityFilterMiddleware(), comments.SaveClusterComment)
	api.GET("/saved/get_my_list", saveclusteraccount.GetMyList)
	api.GET("/saved/delete/:acs_id", saveclusteraccount.DeleteFollowIncident)
	api.POST("/account/report/:account_id", profile.ReportAccount)
	api.GET("/account/get_my_info", account.GetMyInfo)
	api.POST("/account/update_premium_status", account.UpdatePremiumStatus)
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

	// Cronjob endpoints (for manual execution and monitoring)
	api.POST("/cronjob/premium/expire", cronjob.RunPremiumExpirationCheck)
	api.GET("/cronjob/premium/stats", cronjob.GetPremiumStats)
	api.POST("/cronjob/premium/warnings", cronjob.SendExpirationWarnings)
	api.GET("/analytics/location", analyticsHandler.GetLocationAnalytics)

	// comunitacions with apple APN (to send push notifications)
	// ✅ MOVIDO: Endpoints de device tokens sin rate limiting estricto
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

	// ✅ PRODUCCIÓN: Iniciar servidor con configuración segura
	port := ":" + cfg.Port
	log.Printf("🚀 Alertly Backend starting on port %s", cfg.Port)
	log.Printf("🌍 Environment: %s", os.Getenv("GIN_MODE"))
	log.Printf("🔗 Health check: http://localhost%s/health", port)

	if err := router.Run(port); err != nil {
		log.Fatalf("❌ Failed to start server: %v", err)
	}
}

// ✅ Helper function para compatibilidad con desarrollo
func getEnvOrDefault(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
