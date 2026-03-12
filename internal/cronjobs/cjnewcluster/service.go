package cjnewcluster

import (
	"alertly/internal/common"
	"alertly/internal/cronjobs/shared"
	"database/sql"
	"fmt"
	"log"
	"strings"

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
	log.Printf("📬 cjnewcluster found %d pending notifications", len(notifs))
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
			continue
		}
		log.Printf("👥 cjnewcluster cluster %d has %d subscribed users", n.ClusterID, len(users))

		// 3. Send pushes and prepare deliveries
		for _, u := range users {
			title := "New Incident Near You"
			body := fmt.Sprintf("A new '%s' incident has been reported near your saved location: '%s'.", u.SubcategoryName, u.LocationTitle)

			pushData := map[string]interface{}{
				"screen": "ViewIncidentScreen",
				"inclId": fmt.Sprintf("%d", n.ClusterID),
			}

			err := common.SendPush(
				common.ExpoPushMessage{
					Title: title,
					Body:  body,
					Data:  pushData,
				},
				u.DeviceToken,
				payload.NewPayload().
					AlertTitle(title).
					AlertBody(body).
					Custom("screen", "ViewIncidentScreen").
					Custom("inclId", n.ClusterID),
			)
			if err != nil {
				if strings.Contains(err.Error(), "BadDeviceToken") {
					if delErr := s.repo.DeleteStaleToken(u.DeviceToken); delErr != nil {
						log.Printf("cjnewcluster delete stale token: %v", delErr)
					} else {
						log.Printf("🗑️ cjnewcluster deleted stale token for user %d", u.AccountID)
					}
				}
				log.Printf("cjnewcluster push to user %d error: %v", u.AccountID, err)
				continue
			}
			log.Printf("✅ Push sent to user %d (token: %s...)", u.AccountID, u.DeviceToken[:20])
			allDeliveries = append(allDeliveries, shared.Delivery{
				NotificationID: n.ID,
				AccountID:      u.AccountID,
				Title:          title,
				Message:        body,
			})
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
