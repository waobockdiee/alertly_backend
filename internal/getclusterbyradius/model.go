package getclusterbyradius

type Cluster struct {
	InclId          int64   `json:"incl_id"`
	Latitude        float32 `json:"latitude"`
	Longitude       float32 `json:"longitude"`
	InsuId          int64   `json:"insu_id"`
	CategoryCode    string  `json:"category_code"`
	SubcategoryCode string  `json:"subcategory_code"`
}

type Inputs struct {
	Latitude   float64 `uri:"latitude" binding:"required"`
	Longitude  float64 `uri:"longitude" binding:"required"`
	Radius     float64 `uri:"radius" binding:"required"` // en metros
	FromDate   string  `uri:"from_date" binding:"required,datetime=2006-01-02"`
	ToDate     string  `uri:"to_date" binding:"required,datetime=2006-01-02"`
	InsuID     int     `uri:"insu_id"`
	Categories string  `form:"categories"`
}
