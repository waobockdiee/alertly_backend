package main

import (
	"alertly/internal/cronjobs/cjdatabase"
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
	c := cron.New(cron.WithSeconds())

	reponewcluster := cjnewcluster.NewRepository(cjdatabase.DB)
	svcnewcluster := cjnewcluster.NewService(reponewcluster)

	_, err := c.AddFunc("@every 1m", func() {
		log.Println("running cjnewcluster:", time.Now())
		svcnewcluster.Run()
	})
	if err != nil {
		log.Printf("Error running cjnewcluster: %v", err)
	}
	// Iniciar el scheduler
	c.Start()
	log.Println("Cron scheduler started. Waiting tasks...")

	// Bloquear el proceso para mantener el cronjob activo
	select {}
}
