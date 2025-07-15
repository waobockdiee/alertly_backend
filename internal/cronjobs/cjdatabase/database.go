// alertly/database/db.go
package cjdatabase

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
		log.Fatalf("Error connecting DB: %v", err)
	}
	DB.SetMaxOpenConns(50)                  // Ajusta este valor según la carga esperada
	DB.SetMaxIdleConns(10)                  // Número de conexiones ociosas que se mantendrán
	DB.SetConnMaxIdleTime(5 * time.Minute)  //Evita que conexiones ociosas vivan indefinidamente
	DB.SetConnMaxLifetime(30 * time.Minute) //

	// Configuración adicional: p.ej., tamaño del pool
	if err = DB.Ping(); err != nil {
		log.Fatalf("Error pinging DB: %v", err)
	}
}
