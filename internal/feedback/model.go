package feedback

type Feedback struct {
	AccountID   int64  `json:"account_id" db:"account_id"`
	Subject     string `json:"subject" db:"subject"`
	Description string `json:"description" db:"description"`
}
