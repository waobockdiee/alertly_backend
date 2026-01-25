package getclustersbylocation

type Cluster struct {
	InclId          int64        `json:"incl_id"`
	Latitude        float64      `json:"latitude"`
	Longitude       float64      `json:"longitude"`
	InsuId          int64        `json:"insu_id"`
	CategoryCode    string       `json:"category_code"`
	SubcategoryCode string       `json:"subcategory_code"`
	Subcategory     *Subcategory `json:"subcategory"` // ✅ Datos completos de subcategoría
}

// Subcategory contiene la información completa de la subcategoría (JOIN)
type Subcategory struct {
	InsuId             int64  `json:"insu_id"`
	IncaId             int64  `json:"inca_id"`
	Name               string `json:"name"`
	Description        string `json:"description"`
	Icon               string `json:"icon"`     // Usar COALESCE en query
	IconURI            string `json:"icon_uri"` // Usar COALESCE en query
	Code               string `json:"code"`
	MinCircleRange     int64  `json:"min_circle_range"`     // Usar COALESCE en query
	MaxCircleRange     int64  `json:"max_circle_range"`     // Usar COALESCE en query
	DefaultCircleRange int64  `json:"default_circle_range"` // Usar COALESCE en query
	CategoryCode       string `json:"category_code"`
	SubcategoryCode    string `json:"subcategory_code"`
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
