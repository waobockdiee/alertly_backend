package newincident

import "time"

type Media struct {
	Uri string `json:"uri"`
}

type IncidentReport struct {
	InreId             int64   `form:"inre_id"          json:"inre_id"`
	AccountId          int64   `form:"account_id"       json:"account_id"`
	InclId             int64   `form:"incl_id"          json:"incl_id"`
	InsuId             int64   `form:"insu_id"          json:"insu_id"`
	EventType          string  `form:"event_type"       json:"event_type"`
	Description        string  `form:"description"      json:"description"`
	Address            string  `form:"address"          json:"address"`
	City               string  `form:"city"             json:"city"`
	Province           string  `form:"province"         json:"province"`
	PostalCode         string  `form:"postal_code"      json:"postal_code"`
	Latitude           float32 `form:"latitude"         json:"latitude"`
	Longitude          float32 `form:"longitude"        json:"longitude"`
	IsAnonymous        bool    `form:"is_anonymous"     json:"is_anonymous"`
	SubCategoryName    string  `form:"subcategory_name" json:"subcategory_name"`
	Media              Media   `form:"media"            json:"media"`
	MediaType          string  `form:"media_type"       json:"media_type"`
	DefaultCircleRange int     `form:"default_circle_range" json:"default_circle_range"`
	MediaUrl           string  `form:"media_url"        json:"media_url"`
	SubcategoryCode    string  `form:"subcategory_code" json:"subcategory_code"`
	CategoryCode       string  `form:"category_code"    json:"category_code"`
	Vote               *bool   `form:"vote,omitempty"   json:"vote,omitempty"`
	Credibility        float32 `form:"credibility"      json:"credibility"`
}

type Cluster struct {
	InclId          int64      `json:"incl_id"`
	AccountId       int64      `json:"account_id"`
	CreatedAt       *time.Time `json:"created_at"`
	StartTime       *time.Time `json:"start_time"`
	EndTime         *time.Time `json:"end_time"`
	InsuId          int64      `json:"insu_id"`
	MediaUrl        string     `json:"medial_url"`
	CenterLatitude  float32    `json:"center_latitude"`
	CenterLongitude float32    `json:"center_longitude"`
	IsActive        bool       `json:"is_active"`
	MediaType       string     `json:"media_type"`
	EventType       string     `json:"event_type"`
	Description     string     `json:"description"`
	Address         string     `json:"address"`
	City            string     `json:"city"`
	Province        string     `json:"province"`
	PostalCode      string     `json:"postal_code"`
	SubcategoryName string     `json:"subcategory_name"`
	SubcategoryCode string     `json:"subcategory_code"`
	CategoryCode    string     `json:"category_code"`
	Credibility     float32    `json:"credibility"`
	ScoreTrue       float32    `json:"score_true"`
	ScoreFalse      float32    `json:"score_false"`
}
