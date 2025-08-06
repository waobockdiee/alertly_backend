package main

import (
	"alertly/internal/cronjobs/cjdatabase"
	"alertly/internal/cronjobs/cjinactivityreminder"
	"alertly/internal/cronjobs/cjnewcluster"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/joho/godotenv"
	"github.com/robfig/cron/v3"
)

func main() {
	// Cargar variables de entorno
	var err error
	if os.Getenv("NODE_ENV") == "production" {
		err = godotenv.Load(".env.production")
	} else {
		err = godotenv.Load(".env")
	}

	if err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}

	// SET TIMEZONE
	loc, err := time.LoadLocation("America/Edmonton")
	if err != nil {
		log.Fatalf("could not load timezone: %v", err)
	}

	// Configurar la conexión a la base de datos
	dbUser := os.Getenv("DB_USER")
	dbPass := os.Getenv("DB_PASS")
	dbHost := os.Getenv("DB_HOST")
	dbPort := os.Getenv("DB_PORT")
	dbName := os.Getenv("DB_NAME")
	// dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?parseTime=true&multiStatements=true", dbUser, dbPass, dbHost, dbPort, dbName)
	dsn := fmt.Sprintf(
		"%s:%s@tcp(%s:%s)/%s"+
			"?parseTime=true"+
			"&loc=Local"+
			"&timeout=5s"+
			"&readTimeout=2s"+
			"&writeTimeout=2s",
		dbUser, dbPass, dbHost, dbPort, dbName,
	)
	cjdatabase.InitDB(dsn)
	defer cjdatabase.DB.Close()

	// Crear una instancia del scheduler de cron (opción WithSeconds para mayor precisión)
	c := cron.New(
		cron.WithSeconds(),
		cron.WithLocation(loc),
	)

	reponewcluster := cjnewcluster.NewRepository(cjdatabase.DB)
	svcnewcluster := cjnewcluster.NewService(reponewcluster)

	cjinactivityreminderRepo := cjinactivityreminder.NewRepository(cjdatabase.DB)
	cjinactivityreminderService := cjinactivityreminder.NewService(cjinactivityreminderRepo)

	//*******************************
	// EVERY 1 MINUTE
	_, err = c.AddFunc("@every 1m", func() {
		log.Println("running cjnewcluster:", time.Now())
		svcnewcluster.Run()
	})
	if err != nil {
		log.Printf("Error running cjnewcluster: %v", err)
	}

	// EVERY DAY AT 8 AM
	_, err = c.AddFunc("0 0 8 * * *", func() {
		log.Println("running cjinactivityreminder:", time.Now())
		cjinactivityreminderService.Run()
	})
	if err != nil {
		log.Printf("Error running cjinactivityreminder: %v", err)
	}

	// ******************************* GENERATE AUTO NOTIFICATION CREATION
	//*******************************
	// EVERY DAY AT 7:30 AM
	//
	_, err = c.AddFunc("0 30 7 * * *", func() {
		// _, err = c.AddFunc("@every 1m", func() {
		log.Println("running cjinactivityreminder:", time.Now())
		cjinactivityreminderService.Run()
	})
	if err != nil {
		log.Printf("Error running cjinactivityreminder: %v", err)
	}

	// ******************************* END ADDFUNC LOGIC
	// Iniciar el scheduler
	c.Start()
	log.Println("Cron scheduler started. Waiting tasks...")

	// Bloquear el proceso para mantener el cronjob activo
	select {}
}
