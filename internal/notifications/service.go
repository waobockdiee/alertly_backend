package notifications

import (
	"alertly/internal/alerts"
	"errors"
)

type Service interface {
	handleNotification(nType string, accountID int64, referenceID int64) (alerts.Alert, error)
}

type service struct {
	repo Repository
}

func NewService(repo Repository) Service {
	return &service{repo: repo}
}

func (s *service) handleNotification(nType string, accountID int64, referenceID int64) (alerts.Alert, error) {
	var n alerts.Alert
	var err error

	n.Type = nType
	n.AccountID = accountID
	n.ReferenceID = referenceID

	switch n.Type {
	case "welcome_to_app":
		n.Title = "Welcome to Alertly! A community where your actions make a real difference."
		n.Message = ""
		n.Link = "HomeScreen"
		n.MustSendPush = false
		n.MustSendInApp = true
		n.MustBeProcessed = false
		n.ErrorMessage = ""
	case "new_comment":
		n.Title = "New Comment Received! Check it out."
		n.Message = ""
		n.Link = "ViewIncidentScreen"
		n.MustSendPush = true
		n.MustSendInApp = true
		n.MustBeProcessed = true
		n.ErrorMessage = ""
	case "new_cluster": // new incident
		n.Title = "New Incident Reported! Stay informed."
		n.Message = ""
		n.Link = "ViewIncidentScreen"
		n.MustSendPush = true
		n.MustSendInApp = true
		n.MustBeProcessed = true
		n.ErrorMessage = ""
	case "new_incident_cluster": // an update of an incident
		n.Title = "New Incident Update! Stay informed."
		n.Message = ""
		n.Link = "ViewIncidentScreen"
		n.MustSendPush = true
		n.MustSendInApp = true
		n.MustBeProcessed = true
		n.ErrorMessage = ""
	case "earn_citizen_score":
		n.Title = "Congratulations! You've Earned Citizen Points."
		n.Message = "ProfileScreen"
		n.Link = ""
		n.MustSendPush = false
		n.MustSendInApp = true
		n.MustBeProcessed = false
		n.ErrorMessage = ""
	case "membership_expiration_10_days":
		n.Title = "Reminder: Your Membership Expires in 10 Days."
		n.Message = "ProfileScreen" // deberia ir el screen de membership
		n.Link = ""
		n.MustSendPush = true
		n.MustSendInApp = true
		n.MustBeProcessed = false
		n.ErrorMessage = ""
	case "membership_expiration_1_day":
		n.Title = "Urgent: Your Membership Expires Tomorrow!"
		n.Message = "ProfileScreen" // deberia ir el screen de membership
		n.Link = ""
		n.MustSendPush = true
		n.MustSendInApp = true
		n.MustBeProcessed = false
		n.ErrorMessage = ""
	case "welcome_to_membership":
		n.Title = "Welcome to Alertly Membership! Enjoy Exclusive Benefits."
		n.Message = ""
		n.Link = "ProfileScreen" // deberia ir el screen de membership
		n.MustSendPush = false
		n.MustSendInApp = true
		n.MustBeProcessed = false
		n.ErrorMessage = ""
	case "password_reset":
		n.Title = "Password Reset Requested. Secure Your Account."
		n.Message = ""
		n.Link = "ProfileScreen"
		n.MustSendPush = true
		n.MustSendInApp = true
		n.MustBeProcessed = false
		n.ErrorMessage = ""
	case "new_friend_request": // aun no esta funcionando la logica de friends
		n.Title = "You Have a New Friend Request. Connect Now!"
		n.Message = ""
		n.Link = ""
		n.MustSendPush = true
		n.MustSendInApp = true
		n.MustBeProcessed = false
		n.ErrorMessage = ""
	case "user_mentioned":
		n.Title = "You've Been Mentioned! See What They Said."
		n.Message = ""
		n.Link = "ViewIncidentScreen"
		n.MustSendPush = true
		n.MustSendInApp = true
		n.MustBeProcessed = false
		n.ErrorMessage = ""
	case "badge_earned":
		// No sobrescribir título/mensaje - usar los personalizados del cronjob
		// n.Title ya viene del cronjob con el rango específico
		// n.Message ya viene del cronjob con el score específico
		n.Link = "ProfileScreen"
		n.MustSendPush = true
		n.MustSendInApp = true
		n.MustBeProcessed = true // ✅ Cambiar a true para que se procese
		n.ErrorMessage = ""
	case "app_update":
		n.Title = "Alertly App Update Available. Upgrade Now!"
		n.Message = ""
		n.Link = ""
		n.MustSendPush = true
		n.MustSendInApp = true
		n.MustBeProcessed = true
		n.ErrorMessage = ""
	case "promotion":
		n.Title = "Special Promotion: Don't Miss Out on Exclusive Offers."
		n.Message = ""
		n.Link = ""
		n.MustSendPush = true
		n.MustSendInApp = true
		n.MustBeProcessed = true
		n.ErrorMessage = ""
	case "system_maintenance":
		n.Title = "Scheduled Maintenance: Service Updates Coming Soon."
		n.Message = ""
		n.Link = ""
		n.MustSendPush = true
		n.MustSendInApp = true
		n.MustBeProcessed = true
		n.ErrorMessage = ""
	default:
		n.Title = "Notification from Alertly."
		n.Message = ""
		n.Link = ""
		n.MustSendPush = true
		n.MustSendInApp = true
		n.MustBeProcessed = false
		n.ErrorMessage = ""
		err = errors.New("undefined notification type: " + nType)
	}

	return n, err
}
