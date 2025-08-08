package cjcomments

import (
	"alertly/internal/common"
	"alertly/internal/cronjobs/shared"
	"fmt"
	"log"

	"github.com/sideshow/apns2/payload"
)

// Service orquesta la lógica de notificaciones de comentarios.
type Service struct {
	repo      *Repository
	batchSize int64
}

// NewService crea una nueva instancia de Service.
func NewService(r *Repository) *Service {
	return &Service{repo: r, batchSize: 100}
}

// Run procesa las notificaciones de comentarios pendientes.
func (s *Service) Run() {
	log.Println("cjcomments: Running comment notifications cronjob...")

	// 1. Obtener notificaciones de comentarios pendientes
	notifs, err := s.repo.FetchPendingCommentNotifications(s.batchSize)
	if err != nil {
		log.Printf("cjcomments: Error fetching pending comment notifications: %v", err)
		return
	}

	if len(notifs) == 0 {
		log.Println("cjcomments: No pending comment notifications found.")
		return
	}

	var allDeliveries []shared.Delivery
	var processedNotifIDs []int64

	// 2. Procesar cada notificación de comentario
	for _, notif := range notifs {
		// Obtener detalles del comentario
		commentDetails, err := s.repo.GetCommentDetails(notif.CommentID)
		if err != nil {
			log.Printf("cjcomments: Error getting comment details for comment ID %d: %v", notif.CommentID, err)
			processedNotifIDs = append(processedNotifIDs, notif.NotificationID) // Marcar como procesada para no reintentar
			continue
		}

		// Obtener creadores de incidentes para el clúster
		creatorIDs, err := s.repo.GetIncidentCreators(commentDetails.ClusterID)
		if err != nil {
			log.Printf("cjcomments: Error getting incident creators for cluster %d: %v", commentDetails.ClusterID, err)
			continue
		}

		// Obtener usuarios que guardaron el clúster
		savedUserIDs, err := s.repo.GetSavedClusterUsers(commentDetails.ClusterID)
		if err != nil {
			log.Printf("cjcomments: Error getting saved cluster users for cluster %d: %v", commentDetails.ClusterID, err)
			continue
		}

		// Combinar y eliminar duplicados, excluyendo al comentarista
		recipientAccountIDs := make(map[int64]bool)
		for _, id := range creatorIDs {
			if id != commentDetails.CommenterID {
				recipientAccountIDs[id] = true
			}
		}
		for _, id := range savedUserIDs {
			if id != commentDetails.CommenterID {
				recipientAccountIDs[id] = true
			}
		}

		var uniqueRecipientIDs []int64
		for id := range recipientAccountIDs {
			uniqueRecipientIDs = append(uniqueRecipientIDs, id)
		}

		// Obtener tokens de dispositivo para los receptores
		recipients, err := s.repo.GetDeviceTokensForAccounts(uniqueRecipientIDs)
		if err != nil {
			log.Printf("cjcomments: Error getting device tokens for recipients: %v", err)
			continue
		}

		// Enviar notificaciones push y preparar registros de entrega
		for _, recipient := range recipients {
			// Personalizar el mensaje
			title := fmt.Sprintf("Nuevo comentario en %s", commentDetails.SubcategoryName)
			message := fmt.Sprintf("%s ha comentado en un incidente que sigues.", commentDetails.CommentText)
			if len(message) > 200 { // Limitar la longitud del mensaje para push
				message = message[:197] + "..."
			}

			err := common.SendPush(
				common.ExpoPushMessage{Title: title, Body: message},
				recipient.DeviceToken,
				payload.NewPayload().AlertTitle(title).AlertBody(message),
			)
			if err != nil {
				log.Printf("cjcomments: Error sending push to account %d (%s): %v", recipient.AccountID, recipient.DeviceToken, err)
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
			log.Printf("cjcomments: Error inserting notification deliveries: %v", err)
		}
	}

	// 4. Marcar notificaciones como procesadas
	for _, id := range processedNotifIDs {
		if err := s.repo.MarkNotificationProcessed(id); err != nil {
			log.Printf("cjcomments: Error marking notification %d as processed: %v", id, err)
		}
	}

	log.Printf("cjcomments: Processed %d comment notifications.", len(notifs))
}
