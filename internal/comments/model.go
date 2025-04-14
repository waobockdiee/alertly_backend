package comments

import "time"

type InComment struct {
	IncoID        int64     `json:"inco_id"`
	AccountID     int64     `json:"account_id"`
	InclID        int64     `json:"incl_id"`
	CreatedAt     time.Time `json:"created_at"`
	Comment       string    `json:"comment"`
	CounterFlags  int       `json:"counter_flags"`
	CommentStatus bool      `json:"comment_status"`
}

type Comment struct {
	IncoID        int64     `json:"inco_id"`
	AccountID     int64     `json:"account_id"`
	InclID        int64     `json:"incl_id"`
	CreatedAt     time.Time `json:"created_at"`
	Comment       string    `json:"comment"`
	CounterFlags  int       `json:"counter_flags"`
	CommentStatus bool      `json:"comment_status"`
	Nickname      string    `json:"nickname"`
	ThumbnailUrl  string    `json:"thumbnail_url"`
}
