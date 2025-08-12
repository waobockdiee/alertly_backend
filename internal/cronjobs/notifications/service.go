package notifications

import (
	"log"
	"sync"
)

type Service interface {
	ProcessNotifications()
	processWelcomeToApp(n Notification) error
	processBadgeEarned(n Notification) error
}

type service struct {
	repo Repository
}

func NewService(repo Repository) Service {
	return &service{repo: repo}
}

func (s *service) ProcessNotifications() {
	nots, err := s.repo.GetUnprocessedNotificationsPush()
	if err != nil {
		log.Printf("Error al obtener notificaciones: %v", err)
		return
	}

	numWorkers := 5
	notificationChan := make(chan Notification, len(nots))
	var wg sync.WaitGroup

	for i := 0; i < numWorkers; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for n := range notificationChan {
				var err error
				switch n.Type {
				case "welcome_to_app":
					err = s.processWelcomeToApp(n)
				case "badge_earned":
					err = s.processBadgeEarned(n)
				default:
					log.Printf("Acción no definida para el tipo de notificación: %s", n.Type)
				}
				if err != nil {
					log.Printf(
						"Error procesando noti ID %d tipo %s: %v",
						n.NotiID, n.Type, err,
					)
				}
			}
		}()
	}

	for _, n := range nots {
		notificationChan <- n
	}
	close(notificationChan)
	wg.Wait()
}

func (s *service) processWelcomeToApp(n Notification) error {
	accounts, err := s.repo.GetProcessWelcomeToAppAccounts(n)
	if err != nil {
		return err
	}
	if len(accounts) == 0 {
		log.Printf("Sin cuentas para welcome ID %d", n.NotiID)
		return s.repo.UpdateNotificationAsProcessed(n.NotiID)
	}

	deliveries := make([]NotificationDelivery, 0, len(accounts))
	for _, acc := range accounts {
		deliveries = append(deliveries, NotificationDelivery{
			ToAccountID: acc.AccountID,
			NotiID:      n.NotiID,
			Title:       n.Title,
			Message:     n.Message,
		})
	}

	if err := s.repo.BatchSaveNewNotificationDeliveries(deliveries); err != nil {
		log.Printf("Batch insert error noti %d: %v", n.NotiID, err)
		return err
	}
	if err := s.repo.UpdateNotificationAsProcessed(n.NotiID); err != nil {
		log.Printf("Error marcando processed noti %d: %v", n.NotiID, err)
		return err
	}
	return nil
}

func (s *service) processBadgeEarned(n Notification) error {
	// Para badge_earned, creamos una notificación directa al usuario específico
	// No necesitamos buscar múltiples cuentas como en welcome_to_app

	delivery := NotificationDelivery{
		ToAccountID: n.AccountID, // Usar AccountID del modelo
		NotiID:      n.NotiID,
		Title:       n.Title,   // Usar título personalizado del cronjob
		Message:     n.Message, // Usar mensaje personalizado del cronjob
	}

	// Guardar la delivery individual
	if err := s.repo.SaveNotificationDelivery(delivery); err != nil {
		log.Printf("Error saving notification delivery for badge_earned ID %d: %v", n.NotiID, err)
		return err
	}

	// Marcar la notificación como procesada
	if err := s.repo.UpdateNotificationAsProcessed(n.NotiID); err != nil {
		log.Printf("Error marcando processed badge_earned noti %d: %v", n.NotiID, err)
		return err
	}

	log.Printf("Successfully processed badge_earned notification ID %d for account %d", n.NotiID, n.AccountID)
	return nil
}
