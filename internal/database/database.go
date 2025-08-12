// alertly/database/db.go
package database

import (
	"database/sql"
	"log"
	"time"

	_ "github.com/go-sql-driver/mysql" // Asegúrate de instalar el driver
)

var DB *sql.DB

// InitDB inicializa la conexión a la base de datos y la asigna a la variable global DB.
func InitDB(dataSourceName string) {
	var err error
	DB, err = sql.Open("mysql", dataSourceName)
	if err != nil {
		log.Fatalf("Error al conectar a la base de datos: %v", err)
	}

	// ✅ OPTIMIZACIÓN: Configuración más agresiva para alta concurrencia
	DB.SetMaxOpenConns(500)                 // Más conexiones concurrentes
	DB.SetMaxIdleConns(50)                  // Más conexiones inactivas
	DB.SetConnMaxLifetime(15 * time.Minute) // Menor tiempo de vida
	DB.SetConnMaxIdleTime(5 * time.Minute)  // Menor tiempo inactivo

	// Configuración adicional: p.ej., tamaño del pool
	if err = DB.Ping(); err != nil {
		log.Fatalf("Error al hacer ping a la base de datos: %v", err)
	}

	log.Println("✅ Database connection pool optimized for high concurrency")
}
