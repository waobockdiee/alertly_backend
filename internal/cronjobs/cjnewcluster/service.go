package cjnewcluster

import (
	"alertly/internal/common"
	"fmt"
	"log"
)

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

	// 2. Build clusterIDs slice and fetch recipients map
	clusterIDs := make([]int64, len(notifs))
	for i, n := range notifs {
		clusterIDs[i] = n.ClusterID
	}
	recMap, err := s.repo.FetchRecipients(clusterIDs)
	if err != nil {
		log.Printf("cjnewcluster fetch recipients: %v", err)
		return
	}

	var deliveries []Delivery
	// 3. Send pushes and prepare deliveries
	for _, n := range notifs {
		for _, accID := range recMap[n.ClusterID] {
			// 3a. Get device tokens for this user
			tokens, err := s.repo.GetDeviceTokensForAccount(int64(accID))
			if err != nil {
				log.Printf("cjnewcluster tokens for account %d: %v", accID, err)
				continue
			}
			// 3b. Send an Expo push to each token
			for _, tok := range tokens {
				msg := common.ExpoPushMessage{
					To:    tok,
					Title: "Nuevo cluster creado",
					Body:  fmt.Sprintf("Se ha generado el cluster #%%d", n.ClusterID),
					Data:  map[string]interface{}{"clusterId": n.ClusterID},
				}
				if err := common.SendExpoPush(msg); err != nil {
					log.Printf("cjnewcluster expo push to %%s error: %%v", tok, err)
				}
			}
			// 3c. Queue delivery record
			deliveries = append(deliveries, Delivery{NotificationID: n.ID, AccountID: accID})
		}
	}

	// 4. Insert deliveries and mark as processed
	if err := s.repo.InsertDeliveries(deliveries); err != nil {
		log.Printf("cjnewcluster insert deliveries: %v", err)
	}
	if err := s.repo.MarkProcessed(clusterIDs); err != nil {
		log.Printf("cjnewcluster mark processed: %v", err)
	}
}
