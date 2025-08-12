package common

import (
	"database/sql"
	"encoding/json"
	"time"
)

type NullTime struct {
	sql.NullTime
}

func (nt *NullTime) UnmarshalJSON(b []byte) error {
	// Si el valor es null
	if string(b) == "null" {
		nt.Valid = false
		return nil
	}

	// Remover comillas
	var s string
	if err := json.Unmarshal(b, &s); err != nil {
		return err
	}

	// Intentar con fracciones de segundo
	layout := "2006-01-02 15:04:05.000000"
	t, err := time.Parse(layout, s)
	if err != nil {
		// Si falla, probar sin fracciones
		layout = "2006-01-02 15:04:05"
		t, err = time.Parse(layout, s)
		if err != nil {
			return err
		}
	}
	nt.Time = t
	nt.Valid = true
	return nil
}

// MarshalJSON implementa la serialización JSON para NullTime
func (nt *NullTime) MarshalJSON() ([]byte, error) {
	if !nt.Valid {
		return []byte("null"), nil
	}
	return json.Marshal(nt.Time.Format("2006-01-02T15:04:05Z07:00"))
}
