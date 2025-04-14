package getcategories

import "database/sql"

type Category struct {
	IncaId      int64           `json:"inca_id"`
	Name        string          `json:"name"`
	Description string          `json:"description"`
	Icon        *sql.NullString `json:"icon"`
	Code        string          `json:"code"`
}
