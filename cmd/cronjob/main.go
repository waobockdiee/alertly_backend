package main

import (
	"alertly/internal/cronjobs/cjbadgeearn"
	"alertly/internal/cronjobs/cjblockincident"
	"alertly/internal/cronjobs/cjblockuser"
	"alertly/internal/cronjobs/cjcomments"
	"alertly/internal/cronjobs/cjdatabase"
	"alertly/internal/cronjobs/cjinactivityreminder"
	"alertly/internal/cronjobs/cjincidentupdate"
	"alertly/internal/cronjobs/cjincidentexpiration"
	"alertly/internal/cronjobs/cjnewcluster"
	"alertly/internal/cronjobs/cjuserank"
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

	// New Badge Earn Cronjob
	repobadgeearn := cjbadgeearn.NewRepository(cjdatabase.DB)
	svcbadgeearn := cjbadgeearn.NewService(repobadgeearn)

	// New User Rank Cronjob
	repouserank := cjuserank.NewRepository(cjdatabase.DB)
	svcuserank := cjuserank.NewService(repouserank)

	// New Comments Cronjob
	repocomments := cjcomments.NewRepository(cjdatabase.DB)
	svccomments := cjcomments.NewService(repocomments)

	// New Incident Update Cronjob
	repoincidentupdate := cjincidentupdate.NewRepository(cjdatabase.DB)
	svcincidentupdate := cjincidentupdate.NewService(repoincidentupdate)

	// Block users Cronjob
	repoblockuser := cjblockuser.NewRepository(cjdatabase.DB)
	svcblockuser := cjblockuser.NewService(repoblockuser)

	// Block incident Cronjob
	repoblockincident := cjblockincident.NewRepository(cjdatabase.DB)
	svcblockincident := cjblockincident.NewService(repoblockincident)

	// Incident Expiration Cronjob
	repoincidentexpiration := cjincidentexpiration.NewRepository(cjdatabase.DB)
	svcincidentexpiration := cjincidentexpiration.NewService(repoincidentexpiration)

	//*******************************
	// EVERY 1 MINUTE
	_, err = c.AddFunc("@every 1m", func() {
		log.Println("running cjnewcluster:", time.Now())
		svcnewcluster.Run()
	})
	if err != nil {
		log.Printf("Error running cjnewcluster: %v", err)
	}

	// EVERY 1 HOUR
	_, err = c.AddFunc("@every 1h", func() {
		log.Println("running cjblockuser:", time.Now())
		svcblockuser.Run()
	})
	if err != nil {
		log.Printf("Error running cjblockuser: %v", err)
	}

	_, err = c.AddFunc("@every 1m", func() {
		log.Println("running cjblockincident:", time.Now())
		svcblockincident.Run()
	})
	if err != nil {
		log.Printf("Error running cjblockincident: %v", err)
	}

	// EVERY DAY AT 8 AM
	_, err = c.AddFunc("0 0 8 * * *", func() {
		// _, err = c.AddFunc("@every 1m", func() {
		log.Println("running cjinactivityreminder:", time.Now())
		cjinactivityreminderService.Run()
	})
	if err != nil {
		log.Printf("Error running cjinactivityreminder: %v", err)
	}

	// EVERY 10 MINUTES (Badge Earn)
	_, err = c.AddFunc("0 */10 * * * *", func() {
		log.Println("running cjbadgeearn:", time.Now())
		svcbadgeearn.Run()
	})
	if err != nil {
		log.Printf("Error running cjbadgeearn: %v", err)
	}

	// EVERY DAY AT 8 AM (User Rank)
	_, err = c.AddFunc("0 0 8 * * *", func() {
		log.Println("running cjuserank:", time.Now())
		svcuserank.Run()
	})
	if err != nil {
		log.Printf("Error running cjuserank: %v", err)
	}

	// EVERY 2 MINUTE (Comments)
	// _, err = c.AddFunc("@every 2m", func() {
	_, err = c.AddFunc("@every 1m", func() {
		log.Println("running cjcomments:", time.Now())
		svccomments.Run()
	})
	if err != nil {
		log.Printf("Error running cjcomments: %v", err)
	}

	// EVERY 1 MINUTE (Incident Updates)
	_, err = c.AddFunc("@every 1m", func() {
		log.Println("running cjincidentupdate:", time.Now())
		svcincidentupdate.Run()
	})
	if err != nil {
		log.Printf("Error running cjincidentupdate: %v", err)
	}

	// EVERY 5 MINUTES (Incident Expiration)
	_, err = c.AddFunc("0 */5 * * * *", func() {
		log.Println("running cjincidentexpiration:", time.Now())
		svcincidentexpiration.Run()
	})
	if err != nil {
		log.Printf("Error running cjincidentexpiration: %v", err)
	}

	// ******************************* END ADDFUNC LOGIC
	// Iniciar el scheduler
	c.Start()
	log.Println("Cron scheduler started. Waiting tasks...")

	// Bloquear el proceso para mantener el cronjob activo
	select {}
}
