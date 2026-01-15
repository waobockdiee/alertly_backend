// alertly/database/db.go
package database

import (
	"database/sql"
	"log"
	"os"
	"time"

	_ "github.com/lib/pq" // PostgreSQL driver
)

var DB *sql.DB

// InitDB inicializa la conexión a la base de datos y la asigna a la variable global DB.
// Soporta DATABASE_URL (Railway standard) o el DSN tradicional.
func InitDB(dataSourceName string) {
	var err error

	// Railway provee DATABASE_URL automáticamente
	dbURL := os.Getenv("DATABASE_URL")
	if dbURL != "" {
		dataSourceName = dbURL
		log.Println("Using DATABASE_URL from environment")
	}

	DB, err = sql.Open("postgres", dataSourceName)
	if err != nil {
		log.Fatalf("Error al conectar a la base de datos: %v", err)
	}

	// ✅ OPTIMIZACIÓN: Configuración para alta concurrencia
	DB.SetMaxOpenConns(100)                 // PostgreSQL maneja mejor con menos conexiones
	DB.SetMaxIdleConns(25)                  // Conexiones inactivas
	DB.SetConnMaxLifetime(15 * time.Minute) // Tiempo de vida
	DB.SetConnMaxIdleTime(5 * time.Minute)  // Tiempo inactivo

	if err = DB.Ping(); err != nil {
		log.Fatalf("Error al hacer ping a la base de datos: %v", err)
	}

	log.Println("✅ PostgreSQL connection pool initialized")
}
