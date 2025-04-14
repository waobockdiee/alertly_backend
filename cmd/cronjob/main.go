package main

import (
	"alertly/internal/cronjobs/cjdatabase"
	"alertly/internal/cronjobs/notifications"
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
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?parseTime=true", dbUser, dbPass, dbHost, dbPort, dbName)
	cjdatabase.InitDB(dsn)
	defer cjdatabase.DB.Close()

	// Crear una instancia del scheduler de cron (opción WithSeconds para mayor precisión)
	c := cron.New(cron.WithSeconds())

	// Procesamiento de notificaciones cada 1 minuto
	_, err := c.AddFunc("@every 1m", func() {
		repo := notifications.NewRepository(cjdatabase.DB)
		service := notifications.NewService(repo)
		log.Println("Ejecutando procesamiento de notificaciones:", time.Now())
		service.ProcessNotifications()
	})
	if err != nil {
		log.Fatalf("Error programando el procesamiento de notificaciones: %v", err)
	}

	// Iniciar el scheduler
	c.Start()
	log.Println("Cron scheduler iniciado. Esperando tareas...")

	// Bloquear el proceso para mantener el cronjob activo
	select {}
}
