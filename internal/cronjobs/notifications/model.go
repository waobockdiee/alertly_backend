package notifications

import "alertly/internal/common"

type Notification struct {
	NotiID             int64           `db:"noti_id" json:"noti_id"`
	AccountID          int64           `db:"account_id" json:"account_id"`
	Title              string          `db:"title" json:"title"`
	Message            string          `db:"message" json:"message"`
	Type               string          `db:"type" json:"type"`
	IsRead             bool            `db:"is_read" json:"is_read"`
	Link               string          `db:"link" json:"link"`
	CreatedAt          common.NullTime `db:"created_at" json:"created_at"`
	UpdatedAt          common.NullTime `db:"updated_at" json:"updated_at"`
	SentAt             common.NullTime `db:"sent_at" json:"sent_at"`
	MustSendPush       bool            `db:"must_send_as_notification_push" json:"must_send_as_notification_push"`
	MustSendInApp      bool            `db:"must_send_as_notification" json:"must_send_as_notification"`
	MustBeProcessed    bool            `db:"must_be_processed" json:"must_be_processed"`
	ErrorMessage       string          `db:"error_message" json:"error_message"`
	RetryCount         int32           `db:"retry_count" json:"retry_count"`
	ReferenceID        int64           `db:"reference_id" json:"reference_id"`
	Nickname           string          `db:"nickname" json:"nickname"`
	ThumbnailURL       string          `db:"thumbnail_url" json:"thumbnail_url"`
	ClusterMediaURL    string          `db:"media_url" json:"media_url"`
	ClusterEventType   string          `db:"event_type" json:"event_type"`
	ClusterDescription string          `db:"description" json:"description"`
	ClusterAddress     string          `db:"address" json:"address"`
}

type NotificationDelivery struct {
	NodeID      int64           `db:"node_id" json:"node_id"`
	CreatedAt   common.NullTime `db:"created_at" json:"created_at"`
	IsRead      bool            `db:"is_read" json:"is_read"`
	ToAccountID int64           `db:"to_account_id" json:"to_account_id"`
	NotiID      int64           `db:"noti_id" json:"noti_id"`
	Title       string          `db:"title" json:"title"`
	Message     string          `db:"message" json:"message"`
}

type Account struct {
	AccountID int64  `db:"account_id" json:"account_id"`
	Email     string `db:"email" json:"email"`
	Nickname  string `db:"nickname" json:"nickname"`
	Thumbnail string `db:"thumbnail_url" json:"thumbnail_url"`
}
