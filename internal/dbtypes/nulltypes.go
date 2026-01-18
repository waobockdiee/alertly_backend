package dbtypes

import (
	"database/sql/driver"
	"encoding/json"
	"fmt"
	"strconv"
)

// NullString es un tipo que maneja strings NULL de PostgreSQL
// y se serializa correctamente a JSON (como string, no como objeto)
type NullString struct {
	String string
	Valid  bool // Valid es true si String no es NULL
}

// Scan implementa la interfaz sql.Scanner para NullString
func (ns *NullString) Scan(value interface{}) error {
	if value == nil {
		ns.String, ns.Valid = "", false
		return nil
	}

	switch v := value.(type) {
	case string:
		ns.String = v
		ns.Valid = true
		return nil
	case []byte:
		ns.String = string(v)
		ns.Valid = true
		return nil
	default:
		ns.String, ns.Valid = "", false
		return fmt.Errorf("cannot scan type %T into NullString", value)
	}
}

// Value implementa la interfaz driver.Valuer para NullString
func (ns NullString) Value() (driver.Value, error) {
	if !ns.Valid {
		return nil, nil
	}
	return ns.String, nil
}

// MarshalJSON serializa como string primitivo (no objeto)
func (ns NullString) MarshalJSON() ([]byte, error) {
	if !ns.Valid {
		return []byte("null"), nil
	}
	return json.Marshal(ns.String)
}

// UnmarshalJSON deserializa desde JSON
func (ns *NullString) UnmarshalJSON(b []byte) error {
	if string(b) == "null" {
		ns.Valid = false
		return nil
	}

	var s string
	if err := json.Unmarshal(b, &s); err != nil {
		return err
	}

	ns.String = s
	ns.Valid = true
	return nil
}

// NullInt64 es un tipo que maneja integers NULL de PostgreSQL
// y se serializa correctamente a JSON (como number, no como objeto)
type NullInt64 struct {
	Int64 int64
	Valid bool // Valid es true si Int64 no es NULL
}

// Scan implementa la interfaz sql.Scanner para NullInt64
func (ni *NullInt64) Scan(value interface{}) error {
	if value == nil {
		ni.Int64, ni.Valid = 0, false
		return nil
	}

	switch v := value.(type) {
	case int64:
		ni.Int64 = v
		ni.Valid = true
		return nil
	case int32:
		ni.Int64 = int64(v)
		ni.Valid = true
		return nil
	case int:
		ni.Int64 = int64(v)
		ni.Valid = true
		return nil
	case float64:
		ni.Int64 = int64(v)
		ni.Valid = true
		return nil
	case []byte:
		// PostgreSQL puede devolver bytes
		i, err := strconv.ParseInt(string(v), 10, 64)
		if err != nil {
			return err
		}
		ni.Int64 = i
		ni.Valid = true
		return nil
	case string:
		i, err := strconv.ParseInt(v, 10, 64)
		if err != nil {
			return err
		}
		ni.Int64 = i
		ni.Valid = true
		return nil
	default:
		ni.Int64, ni.Valid = 0, false
		return fmt.Errorf("cannot scan type %T into NullInt64", value)
	}
}

// Value implementa la interfaz driver.Valuer para NullInt64
func (ni NullInt64) Value() (driver.Value, error) {
	if !ni.Valid {
		return nil, nil
	}
	return ni.Int64, nil
}

// MarshalJSON serializa como number primitivo (no objeto)
func (ni NullInt64) MarshalJSON() ([]byte, error) {
	if !ni.Valid {
		return []byte("null"), nil
	}
	return json.Marshal(ni.Int64)
}

// UnmarshalJSON deserializa desde JSON
func (ni *NullInt64) UnmarshalJSON(b []byte) error {
	if string(b) == "null" {
		ni.Valid = false
		return nil
	}

	var i int64
	if err := json.Unmarshal(b, &i); err != nil {
		return err
	}

	ni.Int64 = i
	ni.Valid = true
	return nil
}

// NullFloat64 es un tipo que maneja floats NULL de PostgreSQL
// y se serializa correctamente a JSON (como number, no como objeto)
type NullFloat64 struct {
	Float64 float64
	Valid   bool // Valid es true si Float64 no es NULL
}

// Scan implementa la interfaz sql.Scanner para NullFloat64
func (nf *NullFloat64) Scan(value interface{}) error {
	if value == nil {
		nf.Float64, nf.Valid = 0, false
		return nil
	}

	switch v := value.(type) {
	case float64:
		nf.Float64 = v
		nf.Valid = true
		return nil
	case float32:
		nf.Float64 = float64(v)
		nf.Valid = true
		return nil
	case int64:
		nf.Float64 = float64(v)
		nf.Valid = true
		return nil
	case int32:
		nf.Float64 = float64(v)
		nf.Valid = true
		return nil
	case int:
		nf.Float64 = float64(v)
		nf.Valid = true
		return nil
	case []byte:
		f, err := strconv.ParseFloat(string(v), 64)
		if err != nil {
			return err
		}
		nf.Float64 = f
		nf.Valid = true
		return nil
	case string:
		f, err := strconv.ParseFloat(v, 64)
		if err != nil {
			return err
		}
		nf.Float64 = f
		nf.Valid = true
		return nil
	default:
		nf.Float64, nf.Valid = 0, false
		return fmt.Errorf("cannot scan type %T into NullFloat64", value)
	}
}

// Value implementa la interfaz driver.Valuer para NullFloat64
func (nf NullFloat64) Value() (driver.Value, error) {
	if !nf.Valid {
		return nil, nil
	}
	return nf.Float64, nil
}

// MarshalJSON serializa como number primitivo (no objeto)
func (nf NullFloat64) MarshalJSON() ([]byte, error) {
	if !nf.Valid {
		return []byte("null"), nil
	}
	return json.Marshal(nf.Float64)
}

// UnmarshalJSON deserializa desde JSON
func (nf *NullFloat64) UnmarshalJSON(b []byte) error {
	if string(b) == "null" {
		nf.Valid = false
		return nil
	}

	var f float64
	if err := json.Unmarshal(b, &f); err != nil {
		return err
	}

	nf.Float64 = f
	nf.Valid = true
	return nil
}
