package getclustersbylocation

type Cluster struct {
	InclId          int64   `json:"incl_id"`
	Latitude        float32 `json:"latitude"`
	Longitude       float32 `json:"longitude"`
	InsuId          int64   `json:"insu_id"`
	CategoryCode    string  `json:"category_code"`
	SubcategoryCode string  `json:"subcategory_code"`
}

type Inputs struct {
	MinLatitude  float64 `uri:"min_latitude" binding:"required"`
	MaxLatitude  float64 `uri:"max_latitude" binding:"required"`
	MinLongitude float64 `uri:"min_longitude" binding:"required"`
	MaxLongitude float64 `uri:"max_longitude" binding:"required"`
	FromDate     string  `uri:"from_date" binding:"required,datetime=2006-01-02"`
	ToDate       string  `uri:"to_date" binding:"required,datetime=2006-01-02"`
	InsuID       int     `uri:"insu_id"`
	Categories   string  `form:"categories"`
}
