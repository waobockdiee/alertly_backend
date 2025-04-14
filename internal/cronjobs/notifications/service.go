package notifications

import (
	"log"
	"sync"
)

type Service interface {
	ProcessNotifications()
	processWelcomeToApp(n Notification) error
}

type service struct {
	repo Repository
}

func NewService(repo Repository) Service {
	return &service{repo: repo}
}

func (s *service) ProcessNotifications() {
	notifications, err := s.repo.GetUnprocessedNotificationsPush()
	if err != nil {
		log.Printf("Error al obtener notificaciones: %v", err)
		return
	}

	// Definir número de workers para procesar las notificaciones concurrentemente
	numWorkers := 5
	notificationChan := make(chan Notification, len(notifications))
	var wg sync.WaitGroup

	// Iniciar el pool de workers
	for i := 0; i < numWorkers; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for n := range notificationChan {
				var err error
				switch n.Type {
				case "welcome_to_app":
					err = s.processWelcomeToApp(n)
				// Aquí se pueden agregar otros casos según se requiera.
				default:
					log.Printf("No hay acción definida para el tipo de notificación: %s", n.Type)
				}

				if err != nil {
					log.Printf("Error al procesar la notificación (ID: %d) de tipo %s: %v", n.NotiID, n.Type, err)
				}
			}
		}()
	}

	// Enviar las notificaciones al canal para ser procesadas
	for _, n := range notifications {
		notificationChan <- n
	}
	close(notificationChan)
	wg.Wait()
}

// processWelcomeToApp procesa la notificación de bienvenida y realiza un batch insert de deliveries.
func (s *service) processWelcomeToApp(n Notification) error {
	accounts, err := s.repo.GetProcessWelcomeToAppAccounts(n)
	if err != nil {
		return err
	}

	if len(accounts) == 0 {
		log.Printf("No se encontraron cuentas para la notificación de bienvenida (ID: %d)", n.NotiID)
		// Se marca la notificación como procesada aun si no hay cuentas,
		// para evitar re-procesarla repetidamente.
		return s.repo.UpdateNotificationAsProcessed(n.NotiID)
	}

	// Construir un slice de notification deliveries para inserción en batch
	deliveries := make([]NotificationDelivery, 0, len(accounts))
	for _, account := range accounts {
		nd := NotificationDelivery{
			ToAccountID: account.AccountID,
			NotiID:      n.NotiID,
			Title:       n.Title,   // Puedes ajustar el contenido del title según convenga.
			Message:     n.Message, // O definir un mensaje personalizado para bienvenida.
		}
		deliveries = append(deliveries, nd)
	}

	// Realizar el batch insert de las notification deliveries
	err = s.repo.BatchSaveNewNotificationDeliveries(deliveries)
	if err != nil {
		log.Printf("Error en batch insert de notification deliveries para la notificación ID %d: %v", n.NotiID, err)
		return err
	}

	// Actualizar la notificación como procesada
	err = s.repo.UpdateNotificationAsProcessed(n.NotiID)
	if err != nil {
		log.Printf("Error al actualizar la notificación como procesada (ID: %d): %v", n.NotiID, err)
		return err
	}

	return nil
}
