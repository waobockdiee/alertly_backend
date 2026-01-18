package getsubcategoriesbycategoryid

type Subcategory struct {
	InsuId             int64  `json:"insu_id"`
	IncaId             int64  `json:"inca_id"`
	Name               string `json:"name"`
	Description        string `json:"description"`
	Icon               string `json:"icon"`
	Code               string `json:"code"`
	MinCircleRange     int64  `json:"min_circle_range"`
	MaxCircleRange     int64  `json:"max_circle_range"`
	DefaultCircleRange int64  `json:"default_circle_range"`
	CategoryCode       string `json:"category_code"`
	SubcategoryCode    string `json:"subcategory_code"`
}
