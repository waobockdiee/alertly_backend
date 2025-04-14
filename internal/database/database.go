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

	DB.SetMaxOpenConns(200) // O ajustar según la carga esperada
	DB.SetMaxIdleConns(20)  // Más conexiones inactivas disponibles
	DB.SetConnMaxLifetime(30 * time.Minute)

	// Configuración adicional: p.ej., tamaño del pool
	if err = DB.Ping(); err != nil {
		log.Fatalf("Error al hacer ping a la base de datos: %v", err)
	}
}
