package account

import "time"

type MyInfo struct {
	AccountID int64  `db:"account_id" json:"account_id"`
	Email     string `db:"email" json:"email"`
	IsPremium bool   `db:"is_premium" json:"is_premium"`
}
type Account struct {
	AccountID int64 `db:"account_id" json:"account_id"`
}

type History struct {
	HisID       int64     `db:"his_id" json:"his_id"`
	AccountID   int64     `db:"account_id" json:"account_id"`
	InclID      int64     `db:"incl_id" json:"incl_id"`
	CreatedAt   time.Time `db:"created_at" json:"created_at"`
	Address     string    `db:"address" json:"address"`
	Description string    `db:"description" json:"description"`
}

type Counter struct {
	Counter int `json:"counter"`
}
