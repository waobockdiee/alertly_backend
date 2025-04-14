package common

import "log"

// InitLogger inicializa la configuración del logger.
func InitLogger() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)
}
