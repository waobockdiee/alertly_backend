package scheduler

import (
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

	// Aquí puedes agregar más cronjobs en el futuro:
	// go runCommentsCronjob()
	// go runIncidentUpdateCronjob()
	// etc.
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
