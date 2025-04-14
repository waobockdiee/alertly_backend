package getclusterby

import (
	"alertly/internal/comments"
	"alertly/internal/common"
	"database/sql"
)

type Cluster struct {
	InclId                 int64              `json:"incl_id"`
	CreatedAt              sql.NullTime       `json:"created_at"`
	StartTime              sql.NullTime       `json:"start_time"`
	EndTime                sql.NullTime       `json:"end_time"`
	InsuId                 int64              `json:"insu_id"`
	MediaUrl               string             `json:"media_url"`
	CenterLatitude         float64            `json:"center_latitude"`
	CenterLongitude        float64            `json:"center_longitude"`
	IsActive               bool               `json:"is_active"`
	MediaType              string             `json:"media_type"`
	EventType              string             `json:"event_type"`
	Description            string             `json:"description"`
	Address                string             `json:"address"`
	City                   string             `json:"city"`
	Province               string             `json:"province"`
	PostalCode             string             `json:"postal_code"`
	SubcategoryName        string             `json:"subcategory_name"`
	CategoryCode           string             `json:"category_code"`
	SubcategoryCode        string             `json:"subcategory_code"`
	IncidentCount          int                `json:"incident_count"`
	CounterTotalComments   int                `json:"counter_total_comments"`
	CounterTotalVotes      int                `json:"counter_total_votes"`
	CounterTotalViews      int                `json:"counter_total_views"`
	CounterTotalFlags      int                `json:"counter_total_flags"`
	CounterTotalVotesTrue  int                `json:"counter_total_votes_true"`
	CounterTotalVotesFalse int                `json:"counter_total_votes_false"`
	Incidents              []Incident         `json:"incidents"`
	Comments               []comments.Comment `json:"comments"`
	CredibilityPercent     float64            `json:"credibility_percent"`
	GetAccountAlreadyVoted bool               `json:"get_account_already_voted"`
	GetAccountAlreadySaved bool               `json:"get_account_already_saved"`
}

type Comment struct {
	IncoId           int64           `json:"inco_id"`
	Comment          string          `json:"comment"`
	CreatedAt        common.NullTime `json:"created_at"`
	CounterFlags     int             `json:"counter_flags"`
	CommentStatus    string          `json:"comment_status"`
	AccountId        int64           `json:"account_id"`
	Nickname         string          `json:"nickname"`
	FirstName        string          `json:"first_name"`
	LastName         string          `json:"last_name"`
	IsPrivateProfile int8            `json:"is_private_profile"`
	ThumbnailUrl     string          `json:"thumbnail_url"`
}

type Incident struct {
	InreId           int64             `json:"inre_id"`
	MediaUrl         string            `json:"media_url"`
	Latitude         float32           `json:"center_latitude"`
	Longitude        float32           `json:"center_longitude"`
	AccountId        int64             `json:"account_id"`
	Nickname         string            `json:"nickname"`
	FirstName        string            `json:"first_name"`
	LastName         string            `json:"last_name"`
	IsPrivateProfile int8              `json:"is_private_profile"`
	ThumbnailUrl     string            `json:"thumbnail_url"`
	Description      string            `json:"description"`
	EventType        string            `json:"event_type"`
	SubcategortyName string            `json:"subcategory_name"`
	IsAnonymous      string            `json:"is_anonymous"`
	TimeDiff         string            `json:"time_diff"`
	CreatedAt        common.CustomTime `json:"created_at"`
}
