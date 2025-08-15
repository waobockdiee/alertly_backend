package cjincidentupdate

import (
	"alertly/internal/common"
	"alertly/internal/cronjobs/shared"
	"fmt"
	"log"

	"github.com/sideshow/apns2/payload"
)

// Service orquesta la obtención, envío y marcado de notificaciones de updates de incidentes.
type Service struct {
	repo *Repository
}

// NewService crea una nueva instancia de Service.
func NewService(r *Repository) *Service {
	return &Service{repo: r}
}

// Run procesa las notificaciones pendientes de updates de incidentes.
func (s *Service) Run() {
	const batchSize = 100

	// 1. Obtener notificaciones pendientes
	notifs, err := s.repo.FetchPendingIncidentUpdateNotifications(batchSize)
	if err != nil {
		log.Printf("cjincidentupdate: Error fetching pending notifications: %v", err)
		return
	}
	if len(notifs) == 0 {
		return // No hay notificaciones pendientes
	}

	var allDeliveries []shared.Delivery
	var processedNotifIDs []int64

	// 2. Procesar cada notificación de update de incidente
	for _, notif := range notifs {
		// Obtener detalles del cluster
		clusterDetails, err := s.repo.GetClusterDetails(notif.ClusterID)
		if err != nil {
			log.Printf("cjincidentupdate: Error getting cluster details for cluster ID %d: %v", notif.ClusterID, err)
			processedNotifIDs = append(processedNotifIDs, notif.NotificationID) // Marcar como procesada para no reintentar
			continue
		}

		// Obtener creadores de incidentes para el clúster
		creatorIDs, err := s.repo.GetIncidentCreators(notif.ClusterID)
		if err != nil {
			log.Printf("cjincidentupdate: Error getting incident creators for cluster %d: %v", notif.ClusterID, err)
			continue
		}

		// Obtener usuarios que guardaron el clúster
		savedUserIDs, err := s.repo.GetSavedClusterUsers(notif.ClusterID)
		if err != nil {
			log.Printf("cjincidentupdate: Error getting saved cluster users for cluster %d: %v", notif.ClusterID, err)
			continue
		}

		// Combinar y eliminar duplicados, excluyendo al reportero del update
		recipientAccountIDs := make(map[int64]bool)
		for _, id := range creatorIDs {
			if id != notif.ReporterID {
				recipientAccountIDs[id] = true
			}
		}
		for _, id := range savedUserIDs {
			if id != notif.ReporterID {
				recipientAccountIDs[id] = true
			}
		}

		var uniqueRecipientIDs []int64
		for id := range recipientAccountIDs {
			uniqueRecipientIDs = append(uniqueRecipientIDs, id)
		}

		// Si no hay destinatarios, continuar con la siguiente notificación
		if len(uniqueRecipientIDs) == 0 {
			processedNotifIDs = append(processedNotifIDs, notif.NotificationID)
			continue
		}

		// Obtener tokens de dispositivo para los receptores
		recipients, err := s.repo.GetDeviceTokensForAccounts(uniqueRecipientIDs)
		if err != nil {
			log.Printf("cjincidentupdate: Error getting device tokens for accounts: %v", err)
			continue
		}

		// Enviar notificaciones push y preparar registros de entrega
		for _, recipient := range recipients {
			// Personalizar el mensaje
			title := fmt.Sprintf("Incident Update in %s", clusterDetails.SubcategoryName)
			message := fmt.Sprintf("New information has been added to a %s incident you're following in %s.", clusterDetails.SubcategoryName, clusterDetails.City)

			// Limitar la longitud del mensaje para push
			if len(message) > 200 {
				message = message[:197] + "..."
			}

			err := common.SendPush(
				common.ExpoPushMessage{Title: title, Body: message},
				recipient.DeviceToken,
				payload.NewPayload().AlertTitle(title).AlertBody(message),
			)
			if err != nil {
				log.Printf("cjincidentupdate: Error sending push to account %d (%s): %v", recipient.AccountID, recipient.DeviceToken, err)
				continue
			}

			allDeliveries = append(allDeliveries, shared.Delivery{
				NotificationID: notif.NotificationID,
				AccountID:      recipient.AccountID,
				Title:          title,
				Message:        message,
			})
		}

		processedNotifIDs = append(processedNotifIDs, notif.NotificationID)
	}

	// 3. Insertar registros de entrega
	if len(allDeliveries) > 0 {
		if err := s.repo.InsertNotificationDeliveries(allDeliveries); err != nil {
			log.Printf("cjincidentupdate: Error inserting notification deliveries: %v", err)
		}
	}

	// 4. Marcar notificaciones como procesadas
	for _, id := range processedNotifIDs {
		if err := s.repo.MarkNotificationProcessed(id); err != nil {
			log.Printf("cjincidentupdate: Error marking notification %d as processed: %v", id, err)
		}
	}

	log.Printf("cjincidentupdate: Processed %d incident update notifications.", len(notifs))
}

