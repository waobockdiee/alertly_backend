package config

// Config es la estructura que contiene la configuración de la app.
type Config struct {
	Port string
}

// LoadConfig carga y retorna la configuración de la aplicación.
func LoadConfig() Config {
	// Aquí podrías leer variables de entorno o archivos.
	return Config{
		Port: "8080",
	}
}
