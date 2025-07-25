package reportincident

type Report struct {
	AccountID int64  `json:"account_id" db:"account_id"`
	Reason    string `json:"reason" db:"reason"`
	InclID    int64  `json:"incl_id" db:"incl_id"`
	InreID    int64  `json:"inre_id" db:"inre_id"`
}
