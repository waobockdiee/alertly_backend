package saveclusteraccount

import "time"

type Save struct {
	AcsID     int64     `json:"acs_id"`
	AccountID int64     `json:"account_id"`
	InclID    int64     `json:"incl_id"`
	Created   time.Time `json:"created"`
}

type MyList struct {
	AcsID     int64  `json:"acs_id"`
	AccountID int64  `json:"account_id"`
	InclID    int64  `json:"incl_id"`
	MediaUrl  string `json:"media_url"`
}
