package notifications

import (
	"alertly/internal/common"
	"alertly/internal/cronjobs/shared"
	"fmt"
	"log"
	"sync"

	"github.com/sideshow/apns2/payload"
)

type Service interface {
	ProcessNotifications()
	processWelcomeToApp(n Notification) error
	processBadgeEarned(n Notification) error
	processIncidentResult(n Notification) error
	processNewCluster(n Notification) error
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
				case "incident_result_win", "incident_result_loss":
					err = s.processIncidentResult(n)
				case "new_cluster", "new_incident_cluster":
					err = s.processNewCluster(n)
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

	// Obtener device tokens del usuario
	deviceTokens, err := shared.GetDeviceTokensForAccount(s.repo.GetDB(), n.AccountID)
	if err != nil {
		log.Printf("badge_earned: Error getting device tokens for account %d: %v", n.AccountID, err)
		// Continuar de todos modos para guardar la notificación in-app
	}

	// Enviar push notification con screen ProfileScreen
	pushData := map[string]interface{}{
		"screen": "ProfileScreen",
	}

	for _, token := range deviceTokens {
		err := common.SendPush(
			common.ExpoPushMessage{
				Title: n.Title,
				Body:  n.Message,
				Data:  pushData,
			},
			token,
			payload.NewPayload().
				AlertTitle(n.Title).
				AlertBody(n.Message).
				Custom("screen", "ProfileScreen"),
		)
		if err != nil {
			log.Printf("badge_earned: Error sending push to token %s: %v", token[:20], err)
			// Continuar con otros tokens
		} else {
			log.Printf("✅ badge_earned push sent to account %d", n.AccountID)
		}
	}

	delivery := NotificationDelivery{
		ToAccountID: n.AccountID,
		NotiID:      n.NotiID,
		Title:       n.Title,
		Message:     n.Message,
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

func (s *service) processIncidentResult(n Notification) error {
	// incident_result_win y incident_result_loss deben enviar al ViewIncidentScreen con inclId

	// Obtener device tokens del usuario
	deviceTokens, err := shared.GetDeviceTokensForAccount(s.repo.GetDB(), n.AccountID)
	if err != nil {
		log.Printf("incident_result: Error getting device tokens for account %d: %v", n.AccountID, err)
		// Continuar de todos modos para guardar la notificación in-app
	}

	// Enviar push notification con screen ViewIncidentScreen + inclId
	pushData := map[string]interface{}{
		"screen": "ViewIncidentScreen",
		"inclId": fmt.Sprintf("%d", n.ReferenceID),
	}

	for _, token := range deviceTokens {
		err := common.SendPush(
			common.ExpoPushMessage{
				Title: n.Title,
				Body:  n.Message,
				Data:  pushData,
			},
			token,
			payload.NewPayload().
				AlertTitle(n.Title).
				AlertBody(n.Message).
				Custom("screen", "ViewIncidentScreen").
				Custom("inclId", n.ReferenceID),
		)
		if err != nil {
			log.Printf("incident_result: Error sending push to token %s: %v", token[:20], err)
			// Continuar con otros tokens
		} else {
			log.Printf("✅ incident_result push sent to account %d for incident %d", n.AccountID, n.ReferenceID)
		}
	}

	delivery := NotificationDelivery{
		ToAccountID: n.AccountID,
		NotiID:      n.NotiID,
		Title:       n.Title,
		Message:     n.Message,
	}

	// Guardar la delivery individual
	if err := s.repo.SaveNotificationDelivery(delivery); err != nil {
		log.Printf("Error saving notification delivery for incident_result ID %d: %v", n.NotiID, err)
		return err
	}

	// Marcar la notificación como procesada
	if err := s.repo.UpdateNotificationAsProcessed(n.NotiID); err != nil {
		log.Printf("Error marcando processed incident_result noti %d: %v", n.NotiID, err)
		return err
	}

	log.Printf("Successfully processed incident_result notification ID %d for account %d", n.NotiID, n.AccountID)
	return nil
}

func (s *service) processNewCluster(n Notification) error {
	// new_cluster y new_incident_cluster notifican sobre nuevos incidentes/actualizaciones
	// Envía push al owner y guarda delivery para notificación in-app

	// Obtener device tokens del usuario
	deviceTokens, err := shared.GetDeviceTokensForAccount(s.repo.GetDB(), n.AccountID)
	if err != nil {
		log.Printf("new_cluster: Error getting device tokens for account %d: %v", n.AccountID, err)
	}

	// Enviar push notification con screen ViewIncidentScreen + inclId
	pushData := map[string]interface{}{
		"screen": "ViewIncidentScreen",
		"inclId": fmt.Sprintf("%d", n.ReferenceID),
	}

	for _, token := range deviceTokens {
		err := common.SendPush(
			common.ExpoPushMessage{
				Title: n.Title,
				Body:  n.Message,
				Data:  pushData,
			},
			token,
			payload.NewPayload().
				AlertTitle(n.Title).
				AlertBody(n.Message).
				Custom("screen", "ViewIncidentScreen").
				Custom("inclId", n.ReferenceID),
		)
		if err != nil {
			log.Printf("new_cluster: Error sending push to token %s: %v", token[:20], err)
		} else {
			log.Printf("✅ new_cluster push sent to account %d for cluster %d", n.AccountID, n.ReferenceID)
		}
	}

	delivery := NotificationDelivery{
		ToAccountID: n.AccountID,
		NotiID:      n.NotiID,
		Title:       n.Title,
		Message:     n.Message,
	}

	if err := s.repo.SaveNotificationDelivery(delivery); err != nil {
		log.Printf("Error saving notification delivery for new_cluster ID %d: %v", n.NotiID, err)
		return err
	}

	if err := s.repo.UpdateNotificationAsProcessed(n.NotiID); err != nil {
		log.Printf("Error marcando processed new_cluster noti %d: %v", n.NotiID, err)
		return err
	}

	log.Printf("Successfully processed new_cluster notification ID %d for account %d", n.NotiID, n.AccountID)
	return nil
}
