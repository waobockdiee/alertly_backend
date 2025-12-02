package scheduler

import (
	"alertly/internal/cronjobs/cjbot_creator"
	"alertly/internal/cronjobs/cjnewcluster"
	"alertly/internal/cronjobs/notifications"
	"alertly/internal/database"
	"log"
	"time"
)

// StartCronjobs inicia todos los cronjobs en goroutines separadas
func StartCronjobs() {
	log.Println("🕐 Starting internal cronjob scheduler...")

	// Cronjob: new_cluster (notificaciones de nuevos incidentes cerca de places guardados)
	go func() {
		log.Println("🚀 Goroutine for new_cluster cronjob started")
		ticker := time.NewTicker(2 * time.Minute) // Ejecutar cada 2 minutos
		defer ticker.Stop()

		log.Println("✅ Cronjob 'new_cluster' scheduled every 2 minutes")

		// Ejecutar inmediatamente la primera vez
		log.Println("🔥 About to run new_cluster cronjob for the first time...")
		runNewClusterCronjob()
		log.Println("🔥 First run of new_cluster cronjob completed")

		for range ticker.C {
			runNewClusterCronjob()
		}
	}()

	// Cronjob: notifications (procesar notificaciones pendientes y enviar push)
	go func() {
		log.Println("🚀 Goroutine for notifications cronjob started")
		ticker := time.NewTicker(2 * time.Minute) // Ejecutar cada 2 minutos
		defer ticker.Stop()

		log.Println("✅ Cronjob 'notifications' scheduled every 2 minutes")

		// Ejecutar inmediatamente la primera vez
		log.Println("🔥 About to run notifications cronjob for the first time...")
		runNotificationsCronjob()
		log.Println("🔥 First run of notifications cronjob completed")

		for range ticker.C {
			runNotificationsCronjob()
		}
	}()

	// Cronjob: bot_creator_tfs (Toronto Fire Services scraper)
	go func() {
		log.Println("🚀 Goroutine for bot_creator_tfs cronjob started")
		ticker := time.NewTicker(1 * time.Hour) // Ejecutar cada 1 hora
		defer ticker.Stop()

		log.Println("✅ Cronjob 'bot_creator_tfs' scheduled every 1 hour")

		// Ejecutar inmediatamente la primera vez
		log.Println("🔥 About to run bot_creator_tfs cronjob for the first time...")
		runBotCreatorTFSCronjob()
		log.Println("🔥 First run of bot_creator_tfs cronjob completed")

		for range ticker.C {
			runBotCreatorTFSCronjob()
		}
	}()

	// Cronjob: bot_creator_hydro (Toronto Hydro outages scraper)
	// ⚠️ DISABLED: No real API available - only mock data
	// TODO: Enable when real Toronto Hydro API endpoint is found
	// go func() {
	// 	log.Println("🚀 Goroutine for bot_creator_hydro cronjob started")
	// 	ticker := time.NewTicker(1 * time.Hour) // Ejecutar cada 1 hora
	// 	defer ticker.Stop()

	// 	log.Println("✅ Cronjob 'bot_creator_hydro' scheduled every 1 hour")

	// 	// Ejecutar inmediatamente la primera vez
	// 	log.Println("🔥 About to run bot_creator_hydro cronjob for the first time...")
	// 	runBotCreatorHydroCronjob()
	// 	log.Println("🔥 First run of bot_creator_hydro cronjob completed")

	// 	for range ticker.C {
	// 		runBotCreatorHydroCronjob()
	// 	}
	// }()

	log.Println("⚠️  bot_creator_hydro DISABLED (no real API available)")

	// Cronjob: bot_creator_tps (Toronto Police Service scraper)
	go func() {
		log.Println("🚀 Goroutine for bot_creator_tps cronjob started")
		ticker := time.NewTicker(1 * time.Hour) // Ejecutar cada 1 hora
		defer ticker.Stop()

		log.Println("✅ Cronjob 'bot_creator_tps' scheduled every 1 hour")

		// Ejecutar inmediatamente la primera vez
		log.Println("🔥 About to run bot_creator_tps cronjob for the first time...")
		runBotCreatorTPSCronjob()
		log.Println("🔥 First run of bot_creator_tps cronjob completed")

		for range ticker.C {
			runBotCreatorTPSCronjob()
		}
	}()

	// Aquí puedes agregar más cronjobs en el futuro:
	// go runBotCreatorTTCCronjob() (cuando esté implementado)
	// go runBotCreatorWeatherCronjob() (cuando esté implementado)
}

func runNewClusterCronjob() {
	defer func() {
		if r := recover(); r != nil {
			log.Printf("⚠️ Panic in new_cluster cronjob: %v", r)
		}
	}()

	log.Println("🔄 Running new_cluster cronjob...")
	repo := cjnewcluster.NewRepository(database.DB)
	svc := cjnewcluster.NewService(repo)
	svc.Run()
	log.Println("✅ new_cluster cronjob completed")
}

func runNotificationsCronjob() {
	defer func() {
		if r := recover(); r != nil {
			log.Printf("⚠️ Panic in notifications cronjob: %v", r)
		}
	}()

	log.Println("🔄 Running notifications cronjob...")
	repo := notifications.NewRepository(database.DB)
	svc := notifications.NewService(repo)
	svc.ProcessNotifications()
	log.Println("✅ notifications cronjob completed")
}

func runBotCreatorTFSCronjob() {
	defer func() {
		if r := recover(); r != nil {
			log.Printf("⚠️ Panic in bot_creator_tfs cronjob: %v", r)
		}
	}()

	log.Println("🔄 Running bot_creator_tfs cronjob...")
	repo := cjbot_creator.NewRepository(database.DB)
	svc := cjbot_creator.NewService(repo)
	svc.RunTFS()
	log.Println("✅ bot_creator_tfs cronjob completed")
}

func runBotCreatorHydroCronjob() {
	defer func() {
		if r := recover(); r != nil {
			log.Printf("⚠️ Panic in bot_creator_hydro cronjob: %v", r)
		}
	}()

	log.Println("🔄 Running bot_creator_hydro cronjob...")
	repo := cjbot_creator.NewRepository(database.DB)
	svc := cjbot_creator.NewService(repo)
	svc.RunHydro()
	log.Println("✅ bot_creator_hydro cronjob completed")
}

func runBotCreatorTPSCronjob() {
	defer func() {
		if r := recover(); r != nil {
			log.Printf("⚠️ Panic in bot_creator_tps cronjob: %v", r)
		}
	}()

	log.Println("🔄 Running bot_creator_tps cronjob...")
	repo := cjbot_creator.NewRepository(database.DB)
	svc := cjbot_creator.NewService(repo)
	svc.RunTPS()
	log.Println("✅ bot_creator_tps cronjob completed")
}
