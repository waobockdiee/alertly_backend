package main

import (
	"alertly/internal/cronjobs/cjcluster"
	"alertly/internal/cronjobs/cjdatabase"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/joho/godotenv"
	"github.com/robfig/cron/v3"
)

func main() {
	// Cargar variables de entorno
	if err := godotenv.Load(); err != nil {
		log.Println("No se encontró .env, se usarán las variables de entorno del sistema")
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

	//SetClusterToInactiveAndSetAccountScore
	// actualiza el cluster cuando ya han pasado mas de 48horas de creado.
	// actualiza la credibilidad del account
	_, err := c.AddFunc("@every 1m", func() {
		repo := cjcluster.NewRepository(cjdatabase.DB)
		service := cjcluster.NewService(repo)
		log.Println("running SetClusterToInactiveAndSetAccountScore:", time.Now())
		service.SetClusterToInactiveAndSetAccountScore()
	})
	if err != nil {
		log.Fatalf("Error SetClusterToInactiveAndSetAccountScore: %v", err)
	}
	// Iniciar el scheduler
	c.Start()
	log.Println("Cron scheduler started. Waiting tasks...")

	// Bloquear el proceso para mantener el cronjob activo
	select {}
}
