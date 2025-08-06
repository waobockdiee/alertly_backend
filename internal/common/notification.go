package common

import (
	"alertly/internal/alerts"
	"database/sql"
	"fmt"
)

type DBExecutor interface {
	Exec(query string, args ...interface{}) (sql.Result, error)
}

// func SaveNotification(tx *sql.Tx, nType string, accountID int64, referenceID int64) error {
// 	return nil
// }

/*
Comentado porque aun no esta definida la logica final y esta me esta dando un error obvio por cambiar cosas en la tabla notifications
*/
func SaveNotification(dbExec DBExecutor, nType string, accountID int64, referenceID int64) error {
	query := `INSERT INTO notifications(noti_id, owner_account_id, title, message, type, link, must_send_as_notification_push, must_send_as_notification, must_be_processed, error_message, reference_id)
	VALUES(?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`
	n := HandleNotification(nType, accountID, referenceID)
	_, err := dbExec.Exec(query, n.AcnoID, n.AccountID, n.Title, n.Message, n.Type, n.Link, n.MustSendPush, n.MustSendInApp, n.MustBeProcessed, n.ErrorMessage, n.ReferenceID)
	if err != nil {
		return fmt.Errorf("failed to save notification: %w", err)
	}
	return nil
}

func HandleNotification(nType string, accountID int64, referenceID int64) alerts.Alert {
	var n alerts.Alert

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
		return n
	case "new_comment":
		n.Title = "New Comment Received! Check it out."
		n.Message = ""
		n.Link = "ViewIncidentScreen"
		n.MustSendPush = true
		n.MustSendInApp = true
		n.MustBeProcessed = true
		n.ErrorMessage = ""
		return n
	case "new_cluster": // new incident
		n.Title = "New Incident Reported! Stay informed."
		n.Message = ""
		n.Link = "ViewIncidentScreen"
		n.MustSendPush = true
		n.MustSendInApp = true
		n.MustBeProcessed = true
		n.ErrorMessage = ""
		return n
	case "new_incident_cluster": // an update of an incident
		n.Title = "New Incident Update! Stay informed."
		n.Message = ""
		n.Link = "ViewIncidentScreen"
		n.MustSendPush = true
		n.MustSendInApp = true
		n.MustBeProcessed = true
		n.ErrorMessage = ""
		return n
	case "earn_citizen_score":
		n.Title = "Congratulations! You've Earned Citizen Points."
		n.Message = "ProfileScreen"
		n.Link = ""
		n.MustSendPush = false
		n.MustSendInApp = true
		n.MustBeProcessed = false
		n.ErrorMessage = ""
		return n
	case "membership_expiration_10_days":
		n.Title = "Reminder: Your Membership Expires in 10 Days."
		n.Message = "ProfileScreen" // deberia ir el screen de membership
		n.Link = ""
		n.MustSendPush = true
		n.MustSendInApp = true
		n.MustBeProcessed = false
		n.ErrorMessage = ""
		return n
	case "membership_expiration_1_day":
		n.Title = "Urgent: Your Membership Expires Tomorrow!"
		n.Message = "ProfileScreen" // deberia ir el screen de membership
		n.Link = ""
		n.MustSendPush = true
		n.MustSendInApp = true
		n.MustBeProcessed = false
		n.ErrorMessage = ""
		return n
	case "welcome_to_membership":
		n.Title = "Welcome to Alertly Membership! Enjoy Exclusive Benefits."
		n.Message = ""
		n.Link = "ProfileScreen" // deberia ir el screen de membership
		n.MustSendPush = false
		n.MustSendInApp = true
		n.MustBeProcessed = false
		n.ErrorMessage = ""
		return n
	case "password_reset":
		n.Title = "Password Reset Requested. Secure Your Account."
		n.Message = ""
		n.Link = "ProfileScreen"
		n.MustSendPush = true
		n.MustSendInApp = true
		n.MustBeProcessed = false
		n.ErrorMessage = ""
		return n
	case "new_friend_request": // aun no esta funcionando la logica de friends
		n.Title = "You Have a New Friend Request. Connect Now!"
		n.Message = ""
		n.Link = ""
		n.MustSendPush = true
		n.MustSendInApp = true
		n.MustBeProcessed = false
		n.ErrorMessage = ""
		return n
	case "user_mentioned":
		n.Title = "You've Been Mentioned! See What They Said."
		n.Message = ""
		n.Link = "ViewIncidentScreen"
		n.MustSendPush = true
		n.MustSendInApp = true
		n.MustBeProcessed = false
		n.ErrorMessage = ""
		return n
	case "badge_earned":
		n.Title = "Achievement Unlocked! New Badge Earned."
		n.Message = ""
		n.Link = ""
		n.MustSendPush = true
		n.MustSendInApp = true
		n.MustBeProcessed = false
		n.ErrorMessage = ""
		return n
	case "app_update":
		n.Title = "Alertly App Update Available. Upgrade Now!"
		n.Message = ""
		n.Link = ""
		n.MustSendPush = true
		n.MustSendInApp = true
		n.MustBeProcessed = true
		n.ErrorMessage = ""
		return n
	case "promotion":
		n.Title = "Special Promotion: Don't Miss Out on Exclusive Offers."
		n.Message = ""
		n.Link = ""
		n.MustSendPush = true
		n.MustSendInApp = true
		n.MustBeProcessed = true
		n.ErrorMessage = ""
		return n
	case "system_maintenance":
		n.Title = "Scheduled Maintenance: Service Updates Coming Soon."
		n.Message = ""
		n.Link = ""
		n.MustSendPush = true
		n.MustSendInApp = true
		n.MustBeProcessed = true
		n.ErrorMessage = ""
		return n
	case "inactivity_reminder":
		n.Title = "We Miss You at Alertly"
		n.Message = "It's been a while since you last logged in. We have exciting updates and features waiting for you!"
		n.Link = ""
		n.MustSendPush = true
		n.MustSendInApp = true
		n.MustBeProcessed = true
		n.ErrorMessage = ""
		return n
	default:
		n.Title = "Notification from Alertly."
		n.Message = ""
		n.Link = ""
		n.MustSendPush = true
		n.MustSendInApp = true
		n.MustBeProcessed = false
		n.ErrorMessage = ""
		return n
	}

}
