package account

import "time"

type Account struct {
	AccountID int64 `db:"account_id" json:"account_id"`
}

type History struct {
	HisID     int64     `db:"his_id" json:"his_id"`
	AccountID int64     `db:"account_id" json:"account_id"`
	InclID    int64     `db:"incl_id" json:"incl_id"`
	CreatedAt time.Time `db:"created_at" json:"created_at"`
}
