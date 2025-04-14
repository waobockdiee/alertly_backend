package getsubcategoriesbycategoryid

import "database/sql"

type Subcategory struct {
	InsuId             int64           `json:"insu_id"`
	IncaId             int64           `json:"inca_id"`
	Name               string          `json:"name"`
	Description        string          `json:"description"`
	Icon               *sql.NullString `json:"icon"`
	Code               string          `json:"code"`
	MinCircleRange     *sql.NullInt64  `json:"min_circle_range"`
	MaxCircleRange     *sql.NullInt64  `json:"max_circle_range"`
	DefaultCircleRange *sql.NullInt64  `json:"default_circle_range"`
	CategoryCode       string          `json:"category_code"`
	SubcategoryCode    string          `json:"subcategory_code"`
}
