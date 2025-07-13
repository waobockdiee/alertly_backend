package notifications

import (
	"database/sql"
	"log"
)

type Repository interface {
	Save(n Notification) (int64, error)
}

type mysqlRepository struct {
	db *sql.DB
}

func NewRepository(db *sql.DB) Repository {
	return &mysqlRepository{db: db}
}

func (r *mysqlRepository) Save(n Notification) (int64, error) {
	var notiID int64

	query := `INSERT INTO notifications(noti_id, owner_account_id, title, message, type, link, must_send_as_notification_push, must_send_as_notification, must_be_processed, error_message, reference_id)
	VALUES(?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`

	result, err := r.db.Exec(query, n.OwnerAccountID, n.Title, n.Message, n.Type, n.Link, n.MustSendAsNotificationPush, n.MustSendAsNotification, n.MustBeProcessed, n.ErrorMesssage, n.ReferenceID)

	if err != nil {
		log.Println("Error saving notification...")
		return notiID, err
	}

	notiID, err = result.LastInsertId()

	if err != nil {
		log.Println("Error getting noti_id notification...")
		return notiID, err
	}

	log.Println("Notification has been saved succesfully...")
	return notiID, nil
}
