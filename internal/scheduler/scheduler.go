package scheduler

import (
	"alertly/internal/cronjob"
	"alertly/internal/cronjobs/cjbadgeearn"
	"alertly/internal/cronjobs/cjblockincident"
	"alertly/internal/cronjobs/cjblockuser"
	"alertly/internal/cronjobs/cjbot_creator"
	"alertly/internal/cronjobs/cjcomments"
	"alertly/internal/cronjobs/cjinactivityreminder"
	"alertly/internal/cronjobs/cjincidentexpiration"
	"alertly/internal/cronjobs/cjincidentupdate"
	"alertly/internal/cronjobs/cjnewcluster"
	"alertly/internal/cronjobs/cjuserank"
	"alertly/internal/cronjobs/notifications"
	"alertly/internal/database"
	"log"
	"time"
)

// StartCronjobs inicia todos los cronjobs en goroutines separadas
func StartCronjobs() {
	log.Println("🕐 Starting internal cronjob scheduler...")

	// ─── EVERY 2 MINUTES ────────────────────────────────────────────────────────

	// Cronjob: new_cluster (notificaciones de nuevos incidentes cerca de places guardados)
	go func() {
		ticker := time.NewTicker(2 * time.Minute)
		defer ticker.Stop()
		log.Println("✅ Cronjob 'new_cluster' scheduled every 2 minutes")
		runNewClusterCronjob()
		for range ticker.C {
			runNewClusterCronjob()
		}
	}()

	// Cronjob: notifications (procesar notificaciones pendientes y enviar push)
	go func() {
		ticker := time.NewTicker(2 * time.Minute)
		defer ticker.Stop()
		log.Println("✅ Cronjob 'notifications' scheduled every 2 minutes")
		runNotificationsCronjob()
		for range ticker.C {
			runNotificationsCronjob()
		}
	}()

	// Cronjob: comments (notificar a usuarios sobre nuevos comentarios en sus incidentes)
	go func() {
		ticker := time.NewTicker(2 * time.Minute)
		defer ticker.Stop()
		log.Println("✅ Cronjob 'comments' scheduled every 2 minutes")
		runCommentsCronjob()
		for range ticker.C {
			runCommentsCronjob()
		}
	}()

	// Cronjob: incident_update (notificar a usuarios sobre actualizaciones en incidentes guardados)
	go func() {
		ticker := time.NewTicker(2 * time.Minute)
		defer ticker.Stop()
		log.Println("✅ Cronjob 'incident_update' scheduled every 2 minutes")
		runIncidentUpdateCronjob()
		for range ticker.C {
			runIncidentUpdateCronjob()
		}
	}()

	// ─── EVERY 1 HOUR ───────────────────────────────────────────────────────────

	// Cronjob: incident_expiration (expirar incidentes y calcular puntajes de votos)
	go func() {
		ticker := time.NewTicker(1 * time.Hour)
		defer ticker.Stop()
		log.Println("✅ Cronjob 'incident_expiration' scheduled every 1 hour")
		runIncidentExpirationCronjob()
		for range ticker.C {
			runIncidentExpirationCronjob()
		}
	}()

	// Cronjob: premium_expiration (expirar suscripciones premium vencidas)
	go func() {
		ticker := time.NewTicker(1 * time.Hour)
		defer ticker.Stop()
		log.Println("✅ Cronjob 'premium_expiration' scheduled every 1 hour")
		runPremiumExpirationCronjob()
		for range ticker.C {
			runPremiumExpirationCronjob()
		}
	}()

	// Cronjob: bot_creator_tfs (Toronto Fire Services scraper)
	go func() {
		ticker := time.NewTicker(1 * time.Hour)
		defer ticker.Stop()
		log.Println("✅ Cronjob 'bot_creator_tfs' scheduled every 1 hour")
		runBotCreatorTFSCronjob()
		for range ticker.C {
			runBotCreatorTFSCronjob()
		}
	}()

	// Cronjob: bot_creator_tps (Toronto Police Service scraper)
	go func() {
		ticker := time.NewTicker(1 * time.Hour)
		defer ticker.Stop()
		log.Println("✅ Cronjob 'bot_creator_tps' scheduled every 1 hour")
		runBotCreatorTPSCronjob()
		for range ticker.C {
			runBotCreatorTPSCronjob()
		}
	}()

	// ─── EVERY 24 HOURS ─────────────────────────────────────────────────────────

	// Cronjob: badge_earn (otorgar badges basados en actividad del usuario)
	go func() {
		ticker := time.NewTicker(24 * time.Hour)
		defer ticker.Stop()
		log.Println("✅ Cronjob 'badge_earn' scheduled every 24 hours")
		runBadgeEarnCronjob()
		for range ticker.C {
			runBadgeEarnCronjob()
		}
	}()

	// Cronjob: user_rank (actualizar rangos de usuarios según puntaje)
	go func() {
		ticker := time.NewTicker(24 * time.Hour)
		defer ticker.Stop()
		log.Println("✅ Cronjob 'user_rank' scheduled every 24 hours")
		runUserRankCronjob()
		for range ticker.C {
			runUserRankCronjob()
		}
	}()

	// Cronjob: inactivity_reminder (notificar a usuarios inactivos por 7+ días)
	go func() {
		ticker := time.NewTicker(24 * time.Hour)
		defer ticker.Stop()
		log.Println("✅ Cronjob 'inactivity_reminder' scheduled every 24 hours")
		runInactivityReminderCronjob()
		for range ticker.C {
			runInactivityReminderCronjob()
		}
	}()

	// Cronjob: block_user (bloquear usuarios con 20+ reportes)
	go func() {
		ticker := time.NewTicker(24 * time.Hour)
		defer ticker.Stop()
		log.Println("✅ Cronjob 'block_user' scheduled every 24 hours")
		runBlockUserCronjob()
		for range ticker.C {
			runBlockUserCronjob()
		}
	}()

	// Cronjob: block_incident (rechazar incidentes con 5+ flags)
	go func() {
		ticker := time.NewTicker(24 * time.Hour)
		defer ticker.Stop()
		log.Println("✅ Cronjob 'block_incident' scheduled every 24 hours")
		runBlockIncidentCronjob()
		for range ticker.C {
			runBlockIncidentCronjob()
		}
	}()

	// ⚠️ DISABLED - pending real API:
	// bot_creator_hydro (no real Toronto Hydro API available)
	// bot_creator_ttc   (not yet implemented)
	// bot_creator_weather (not yet implemented)

	log.Println("🚀 All cronjobs scheduled and running")
}

// ─── RUNNER FUNCTIONS ────────────────────────────────────────────────────────

func runNewClusterCronjob() {
	defer func() {
		if r := recover(); r != nil {
			log.Printf("⚠️ Panic in new_cluster cronjob: %v", r)
		}
	}()
	repo := cjnewcluster.NewRepository(database.DB)
	svc := cjnewcluster.NewService(repo)
	svc.Run()
}

func runNotificationsCronjob() {
	defer func() {
		if r := recover(); r != nil {
			log.Printf("⚠️ Panic in notifications cronjob: %v", r)
		}
	}()
	repo := notifications.NewRepository(database.DB)
	svc := notifications.NewService(repo)
	svc.ProcessNotifications()
}

func runCommentsCronjob() {
	defer func() {
		if r := recover(); r != nil {
			log.Printf("⚠️ Panic in comments cronjob: %v", r)
		}
	}()
	repo := cjcomments.NewRepository(database.DB)
	svc := cjcomments.NewService(repo)
	svc.Run()
}

func runIncidentUpdateCronjob() {
	defer func() {
		if r := recover(); r != nil {
			log.Printf("⚠️ Panic in incident_update cronjob: %v", r)
		}
	}()
	repo := cjincidentupdate.NewRepository(database.DB)
	svc := cjincidentupdate.NewService(repo)
	svc.Run()
}

func runIncidentExpirationCronjob() {
	defer func() {
		if r := recover(); r != nil {
			log.Printf("⚠️ Panic in incident_expiration cronjob: %v", r)
		}
	}()
	repo := cjincidentexpiration.NewRepository(database.DB)
	svc := cjincidentexpiration.NewService(repo)
	svc.Run()
}

func runPremiumExpirationCronjob() {
	defer func() {
		if r := recover(); r != nil {
			log.Printf("⚠️ Panic in premium_expiration cronjob: %v", r)
		}
	}()
	svc := cronjob.NewPremiumExpirationService(database.DB)
	if err := svc.CheckAndExpirePremiumAccounts(); err != nil {
		log.Printf("❌ Error in premium_expiration cronjob: %v", err)
	}
}

func runBotCreatorTFSCronjob() {
	defer func() {
		if r := recover(); r != nil {
			log.Printf("⚠️ Panic in bot_creator_tfs cronjob: %v", r)
		}
	}()
	repo := cjbot_creator.NewRepository(database.DB)
	svc := cjbot_creator.NewService(repo)
	svc.RunTFS()
}

func runBotCreatorTPSCronjob() {
	defer func() {
		if r := recover(); r != nil {
			log.Printf("⚠️ Panic in bot_creator_tps cronjob: %v", r)
		}
	}()
	repo := cjbot_creator.NewRepository(database.DB)
	svc := cjbot_creator.NewService(repo)
	svc.RunTPS()
}

func runBadgeEarnCronjob() {
	defer func() {
		if r := recover(); r != nil {
			log.Printf("⚠️ Panic in badge_earn cronjob: %v", r)
		}
	}()
	repo := cjbadgeearn.NewRepository(database.DB)
	svc := cjbadgeearn.NewService(repo)
	svc.Run()
}

func runUserRankCronjob() {
	defer func() {
		if r := recover(); r != nil {
			log.Printf("⚠️ Panic in user_rank cronjob: %v", r)
		}
	}()
	repo := cjuserank.NewRepository(database.DB)
	svc := cjuserank.NewService(repo)
	svc.Run()
}

func runInactivityReminderCronjob() {
	defer func() {
		if r := recover(); r != nil {
			log.Printf("⚠️ Panic in inactivity_reminder cronjob: %v", r)
		}
	}()
	repo := cjinactivityreminder.NewRepository(database.DB)
	svc := cjinactivityreminder.NewService(repo)
	svc.Run()
}

func runBlockUserCronjob() {
	defer func() {
		if r := recover(); r != nil {
			log.Printf("⚠️ Panic in block_user cronjob: %v", r)
		}
	}()
	repo := cjblockuser.NewRepository(database.DB)
	svc := cjblockuser.NewService(repo)
	svc.Run()
}

func runBlockIncidentCronjob() {
	defer func() {
		if r := recover(); r != nil {
			log.Printf("⚠️ Panic in block_incident cronjob: %v", r)
		}
	}()
	repo := cjblockincident.NewRepository(database.DB)
	svc := cjblockincident.NewService(repo)
	svc.Run()
}
