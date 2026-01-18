package getcategories

type Category struct {
	IncaId      int64  `json:"inca_id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Icon        string `json:"icon"`
	Code        string `json:"code"`
}
