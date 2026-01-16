package dbtypes

import (
	"database/sql/driver"
	"encoding/json"
	"fmt"
	"strings"
)

// NullBool es un tipo que maneja booleanos en PostgreSQL que pueden ser:
// - BOOLEAN (true/false)
// - SMALLINT (0/1)
// - CHAR(1) ('0'/'1', puede tener espacios: '1 ')
// - VARCHAR con valores como "true", "false", "t", "f"
type NullBool struct {
	Bool  bool
	Valid bool // Valid es true si Bool no es NULL
}

// Scan implementa la interfaz sql.Scanner para NullBool
func (nb *NullBool) Scan(value interface{}) error {
	if value == nil {
		nb.Bool, nb.Valid = false, false
		return nil
	}

	switch v := value.(type) {
	case bool:
		nb.Bool = v
		nb.Valid = true
		return nil
	case int64:
		nb.Bool = v == 1
		nb.Valid = true
		return nil
	case int32:
		nb.Bool = v == 1
		nb.Valid = true
		return nil
	case int:
		nb.Bool = v == 1
		nb.Valid = true
		return nil
	case []byte:
		// Manejar bytes que vienen de PostgreSQL CHAR/VARCHAR
		s := strings.TrimSpace(string(v))
		nb.Bool = s == "1" || s == "t" || s == "true" || s == "TRUE"
		nb.Valid = true
		return nil
	case string:
		// Manejar strings directamente
		s := strings.TrimSpace(v)
		nb.Bool = s == "1" || s == "t" || s == "true" || s == "TRUE"
		nb.Valid = true
		return nil
	default:
		nb.Bool, nb.Valid = false, false
		return fmt.Errorf("cannot scan type %T into NullBool", value)
	}
}

// Value implementa la interfaz driver.Valuer para NullBool
// Esto permite que NullBool se use en queries INSERT/UPDATE
func (nb NullBool) Value() (driver.Value, error) {
	if !nb.Valid {
		return nil, nil
	}
	return nb.Bool, nil
}

// UnmarshalJSON implementa la deserialización JSON para NullBool
func (nb *NullBool) UnmarshalJSON(b []byte) error {
	if string(b) == "null" {
		nb.Valid = false
		return nil
	}

	var boolVal bool
	if err := json.Unmarshal(b, &boolVal); err != nil {
		return err
	}

	nb.Bool = boolVal
	nb.Valid = true
	return nil
}

// MarshalJSON implementa la serialización JSON para NullBool
func (nb NullBool) MarshalJSON() ([]byte, error) {
	if !nb.Valid {
		return []byte("null"), nil
	}
	return json.Marshal(nb.Bool)
}

// BoolToInt convierte un bool a int para insertar en columnas SMALLINT de PostgreSQL
// Uso: db.Exec("INSERT ... VALUES ($1)", dbtypes.BoolToInt(myBool))
func BoolToInt(b bool) int {
	if b {
		return 1
	}
	return 0
}
