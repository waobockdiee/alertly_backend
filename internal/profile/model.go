package profile

import "database/sql"

type ProfileCard struct {
	AccountID    int64  `json:"account_id"`
	Nickname     string `json:"nickname"`
	FirstName    string `json:"first_name"`
	LastName     string `json:"last_name"`
	ThumbnailUrl string `json:"thumbnail_url"`
	PhoneNumber  string `json:"phone_number"`
}

type Incident struct {
	InreId          int64  `json:"inre_id"`
	MedialUrl       string `json:"media_url"`
	Description     string `json:"description"`
	EventType       string `json:"event_type"`
	SubcategoryName string `json:"subcategory_name"`
}

type Range struct {
	Title           string `json:"title"`
	BackgroundColor string `json:"background_color"`
	TextColor       string `json:"text_color"`
	Description     string `json:"description"`
}

type AccountMedal struct {
	AcmeID int64 `json:"acme_id"`
}

type Profile struct {
	AccountID                    int64          `json:"account_id"`
	Nickname                     string         `json:"nickname"`
	FirstName                    string         `json:"first_name"`
	LastName                     string         `json:"last_name"`
	PhoneNumber                  sql.NullString `json:"phone_number"`
	Range                        Range          `json:"range"`
	Status                       string         `json:"status"`
	Credibility                  int            `json:"credibility"`
	IsPrivateProfile             bool           `json:"is_private_profile"`
	Score                        int            `json:"score"`
	IsPremium                    bool           `json:"is_premium"`
	CounterTotalIncidentsCreated int            `json:"counter_total_incidents_created"`
	CounterTotalVotesMade        int            `json:"counter_total_votes_made"`
	CounterTotalCommentsMade     int            `json:"counter_total_comments_made"`
	CounterTotalLocations        int            `json:"counter_total_locations"`
	CounterTotalFlags            int            `json:"counter_total_flags"`
	CounterTotalMedals           int            `json:"counter_total_medals"`
	BirthYear                    string         `json:"birth_year"`
	BirthMonth                   string         `json:"birth_month"`
	BirthDay                     string         `json:"birth_day"`
	HasFinishedTutorial          bool           `json:"has_finished_tutorial"`
	HasWatchNewIncidentTutorial  bool           `json:"has_watch_new_incident_tutorial"`
	ThumbnailUrl                 string         `json:"thumbnail_url"`
	Incidents                    []Incident     `json:"incidents"`
	Medals                       []AccountMedal `json:"medals"`
}

type FullProfile struct {
}
