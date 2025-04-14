package achievements

import "time"

type Achievement struct {
	AcacID      int64     `json:"acac_id"`
	AccountID   int64     `json:"account_id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	Created     time.Time `json:"created"`
	ShowInModal bool      `json:"show_in_modal"`
	Type        string    `json:"type"`
	TextToShow  string    `json:"text_to_show"`
	IconUrl     string    `json:"icon_url"`
}
