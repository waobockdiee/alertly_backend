package notifications

import (
	"database/sql"
	"fmt"
	"log"
)

type Repository interface {
	Save(n Notification) (int64, error)
	SaveDeviceToken(accountID int64, token string) error
	DeleteDeviceToken(accountID int64, deviceToken string) error
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

func (r *mysqlRepository) SaveDeviceToken(accountID int64, token string) error {
	query := `
        INSERT INTO device_tokens (account_id, device_token)
        VALUES (?, ?)
        ON DUPLICATE KEY UPDATE updated_at = CURRENT_TIMESTAMP;
    `
	if _, err := r.db.Exec(query, accountID, token); err != nil {
		return fmt.Errorf("SaveDeviceToken: %w", err)
	}
	return nil
}

func (r *mysqlRepository) DeleteDeviceToken(accountID int64, deviceToken string) error {
	_, err := r.db.Exec(`
	  DELETE FROM device_tokens 
	  WHERE account_id = ? AND device_token = ?`,
		accountID, deviceToken,
	)
	return err
}
