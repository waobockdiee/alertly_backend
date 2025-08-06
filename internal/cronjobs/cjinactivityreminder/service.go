// service.go
package cjinactivityreminder

import (
	"alertly/internal/common"
	"alertly/internal/cronjobs/shared"
	"log"

	"github.com/sideshow/apns2/payload"
)

const (
	// Tamaño del batch para fetch y envíos
	batchSize = 100
)

// Service orchestrates fetching, sending, and marking notifications
type Service struct {
	repo *Repository
}

// NewService creates a new Service instance
func NewService(r *Repository) *Service {
	return &Service{repo: r}
}

func (s *Service) Run() {
	// 1) Generar notificaciones pendientes (> X días sin actividad)
	if err := s.repo.GenerateInactivityNotifications(); err != nil {
		log.Printf("[Inactivity] generación de notificaciones falló: %v", err)
		// seguimos para procesar las que existan
	}

	// 2) Recorrer en batches todas las notificaciones pendientes
	var lastID int64 = 0
	var allSentIDs []int64
	var allDeliveries []shared.Delivery

	for {
		notis, err := s.repo.FetchPending(batchSize, lastID)
		if err != nil {
			log.Printf("[Inactivity] error fetch pending: %v", err)
			break
		}
		if len(notis) == 0 {
			break
		}

		// Para la siguiente iteración usamos keyset pagination
		lastID = notis[len(notis)-1].NotificationID

		// Enviar cada push y acumular deliveries
		title := "We miss you at Alertly"
		body := "It’s been a while! Come back and see what’s new today."
		for _, n := range notis {
			if err := common.SendPush(
				common.ExpoPushMessage{Title: title, Body: body},
				n.DeviceToken,
				payload.NewPayload().AlertTitle(title).AlertBody(body),
			); err != nil {
				log.Printf("[Inactivity] error push acct %d: %v", n.AccountID, err)
				continue
			}
			// Acumular delivery exitoso
			allSentIDs = append(allSentIDs, n.NotificationID)
			allDeliveries = append(allDeliveries, shared.Delivery{
				NotificationID: n.NotificationID,
				AccountID:      n.AccountID,
				Title:          title,
				Message:        body,
			})
		}
	}

	// 3) Insertar todos los registros en notification_deliveries en un batch
	if len(allDeliveries) > 0 {
		if err := s.repo.InsertDeliveries(allDeliveries); err != nil {
			log.Printf("[Inactivity] error insert deliveries: %v", err)
		}
	}

	// 4) Marcar todas las notificaciones enviadas como procesadas
	if len(allSentIDs) > 0 {
		if err := s.repo.MarkProcessed(allSentIDs); err != nil {
			log.Printf("[Inactivity] error mark processed: %v", err)
		}
	}

	log.Printf("[Inactivity] reminders sent: %d", len(allSentIDs))
}
