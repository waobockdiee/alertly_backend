package common

import (
	"fmt"
	"strings"
	"time"
)

const customTimeLayout = "2006-01-02 15:04:05.000000"

type CustomTime struct {
	time.Time
}

// UnmarshalJSON implementa la interfaz json.Unmarshaler para CustomTime.
// Se utiliza time.ParseInLocation para interpretar la fecha en la zona local.
func (ct *CustomTime) UnmarshalJSON(b []byte) error {
	// Convertir el slice de bytes a string y eliminar las comillas
	s := strings.Trim(string(b), "\"")
	// Parsear la fecha usando el layout definido en la zona local
	t, err := time.ParseInLocation(customTimeLayout, s, time.Local)
	if err != nil {
		return fmt.Errorf("error al parsear la fecha %s: %w", s, err)
	}
	ct.Time = t
	return nil
}

// TimeAgo devuelve una cadena legible que indica el tiempo transcurrido desde t.
func TimeAgo(t time.Time) string {
	now := time.Now()
	diff := now.Sub(t)

	if diff < time.Hour {
		minutes := int(diff.Minutes())
		if minutes < 1 {
			minutes = 1 // Mínimo 1 minuto
		}
		if minutes == 1 {
			return "1 minute ago"
		}
		return fmt.Sprintf("%d minutes ago", minutes)
	} else if diff < 24*time.Hour {
		hours := int(diff.Hours())
		if hours == 1 {
			return "1 hour ago"
		}
		return fmt.Sprintf("%d hours ago", hours)
	} else if diff < 30*24*time.Hour {
		days := int(diff.Hours() / 24)
		if days == 1 {
			return "1 day ago"
		}
		return fmt.Sprintf("%d days ago", days)
	} else if diff < 365*24*time.Hour {
		months := int(diff.Hours() / (24 * 30))
		if months == 1 {
			return "1 month ago"
		}
		return fmt.Sprintf("%d months ago", months)
	} else {
		years := int(diff.Hours() / (24 * 365))
		if years == 1 {
			return "1 year ago"
		}
		return fmt.Sprintf("%d years ago", years)
	}
}
