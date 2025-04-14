package newincident

import "time"

type Media struct {
	Uri string `json:"uri"`
}

type IncidentReport struct {
	InreId             int64   `form:"inre_id"`
	AccountId          int64   `form:"account_id"`
	InclId             int64   `form:"incl_id"`
	InsuId             int64   `form:"insu_id"`
	EventType          string  `form:"event_type"`
	Description        string  `form:"description"`
	Address            string  `form:"address"`
	City               string  `form:"city"`
	Province           string  `form:"province"`
	PostalCode         string  `form:"postal_code"`
	Latitude           float32 `form:"latitude"`
	Longitude          float32 `form:"longitude"`
	IsAnonymous        bool    `form:"is_anonymous"`
	SubCategoryName    string  `form:"subcategory_name"`
	Media              Media   `form:"media"`
	MediaType          string  `form:"media_type"`
	DefaultCircleRange int     `form:"default_circle_range"`
	MediaUrl           string  `form:"media_url"`
	SubcategoryCode    string  `form:"subcategory_code"`
	CategoryCode       string  `form:"code"`
}

type Cluster struct {
	InclId          int64      `json:"incl_id"`
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
	CategoryCode    string     `json:"code"`
}
