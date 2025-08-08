package cjnewcluster

import (
	"alertly/internal/common"
	"alertly/internal/cronjobs/shared"
	"database/sql"
	"fmt"
	"log"

	"github.com/sideshow/apns2/payload"
)

// Notification represents a pending notification to be processed
type Notification struct {
	ID        int64
	ClusterID int64
	CreatedAt sql.NullTime
}

// Service orchestrates fetching, sending, and marking notifications
type Service struct {
	repo      *Repository
	batchSize int64
}

// NewService creates a new Service instance
func NewService(r *Repository) *Service {
	return &Service{repo: r, batchSize: 100}
}

// Run processes pending notifications every cron tick
func (s *Service) Run() {
	// 1. Fetch pending notifications
	notifs, err := s.repo.FetchPending(s.batchSize)
	if err != nil {
		log.Printf("cjnewcluster fetch pending: %v", err)
		return
	}
	if len(notifs) == 0 {
		return
	}

	var allDeliveries []shared.Delivery
	var processedNotifIDs []int64

	// 2. Process each notification
	for _, n := range notifs {
		users, err := s.repo.FindSubscribedUsersForCluster(n.ClusterID)
		if err != nil {
			log.Printf("cjnewcluster find users for cluster %d: %v", n.ClusterID, err)
			continue // Skip to next notification on error
		}

		// 3. Send pushes and prepare deliveries
		for _, u := range users {
			title := "New Incident Near You"
			body := fmt.Sprintf("A new '%s' incident has been reported near your saved location: '%s'.", u.SubcategoryName, u.LocationTitle)

			err := common.SendPush(
				common.ExpoPushMessage{Title: title, Body: body},
				u.DeviceToken,
				payload.NewPayload().AlertTitle(title).AlertBody(body),
			)
			if err != nil {
				log.Printf("cjnewcluster expo push to %s error: %v", u.DeviceToken, err)
				continue
			}
			// Queue delivery record
			allDeliveries = append(allDeliveries, shared.Delivery{NotificationID: n.ID, AccountID: u.AccountID})
		}
		processedNotifIDs = append(processedNotifIDs, n.ID)
	}

	// 4. Insert deliveries and mark as processed
	if len(allDeliveries) > 0 {
		if err := s.repo.InsertDeliveries(allDeliveries); err != nil {
			log.Printf("cjnewcluster insert deliveries: %v", err)
		}
	}

	if len(processedNotifIDs) > 0 {
		if err := s.repo.MarkProcessed(processedNotifIDs); err != nil {
			log.Printf("cjnewcluster mark processed: %v", err)
		}
	}

	log.Printf("cjnewcluster processed %d notifications and sent %d pushes", len(processedNotifIDs), len(allDeliveries))
}
