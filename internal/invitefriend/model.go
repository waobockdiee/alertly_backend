package invitefriend

type Invitation struct {
	AccountID          int64 `json:"account_id" db:"account_id"`
	CitizenScoreEarned uint8 `json:"citizen_score_earned" db:"citizen_score_earned"`
}
