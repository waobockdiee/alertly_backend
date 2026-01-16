package dbtypes

import (
	"database/sql/driver"
	"encoding/json"
	"testing"
)

func TestNullBool_Scan(t *testing.T) {
	tests := []struct {
		name     string
		input    interface{}
		wantBool bool
		wantValid bool
		wantErr  bool
	}{
		{
			name:     "nil value",
			input:    nil,
			wantBool: false,
			wantValid: false,
			wantErr:  false,
		},
		{
			name:     "bool true",
			input:    true,
			wantBool: true,
			wantValid: true,
			wantErr:  false,
		},
		{
			name:     "bool false",
			input:    false,
			wantBool: false,
			wantValid: true,
			wantErr:  false,
		},
		{
			name:     "int64 one",
			input:    int64(1),
			wantBool: true,
			wantValid: true,
			wantErr:  false,
		},
		{
			name:     "int64 zero",
			input:    int64(0),
			wantBool: false,
			wantValid: true,
			wantErr:  false,
		},
		{
			name:     "int32 one",
			input:    int32(1),
			wantBool: true,
			wantValid: true,
			wantErr:  false,
		},
		{
			name:     "int one",
			input:    1,
			wantBool: true,
			wantValid: true,
			wantErr:  false,
		},
		{
			name:     "string '1'",
			input:    "1",
			wantBool: true,
			wantValid: true,
			wantErr:  false,
		},
		{
			name:     "string '1 ' (with trailing space)",
			input:    "1 ",
			wantBool: true,
			wantValid: true,
			wantErr:  false,
		},
		{
			name:     "string '0'",
			input:    "0",
			wantBool: false,
			wantValid: true,
			wantErr:  false,
		},
		{
			name:     "string 'true'",
			input:    "true",
			wantBool: true,
			wantValid: true,
			wantErr:  false,
		},
		{
			name:     "string 'TRUE'",
			input:    "TRUE",
			wantBool: true,
			wantValid: true,
			wantErr:  false,
		},
		{
			name:     "string 't'",
			input:    "t",
			wantBool: true,
			wantValid: true,
			wantErr:  false,
		},
		{
			name:     "string 'false'",
			input:    "false",
			wantBool: false,
			wantValid: true,
			wantErr:  false,
		},
		{
			name:     "bytes []byte('1')",
			input:    []byte("1"),
			wantBool: true,
			wantValid: true,
			wantErr:  false,
		},
		{
			name:     "bytes []byte('1 ')",
			input:    []byte("1 "),
			wantBool: true,
			wantValid: true,
			wantErr:  false,
		},
		{
			name:     "bytes []byte('t')",
			input:    []byte("t"),
			wantBool: true,
			wantValid: true,
			wantErr:  false,
		},
		{
			name:     "unsupported type float64",
			input:    float64(1.0),
			wantBool: false,
			wantValid: false,
			wantErr:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var nb NullBool
			err := nb.Scan(tt.input)

			if (err != nil) != tt.wantErr {
				t.Errorf("NullBool.Scan() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if nb.Bool != tt.wantBool {
				t.Errorf("NullBool.Scan() Bool = %v, want %v", nb.Bool, tt.wantBool)
			}

			if nb.Valid != tt.wantValid {
				t.Errorf("NullBool.Scan() Valid = %v, want %v", nb.Valid, tt.wantValid)
			}
		})
	}
}

func TestNullBool_Value(t *testing.T) {
	tests := []struct {
		name     string
		nb       NullBool
		want     driver.Value
		wantErr  bool
	}{
		{
			name: "valid true",
			nb:   NullBool{Bool: true, Valid: true},
			want: true,
			wantErr: false,
		},
		{
			name: "valid false",
			nb:   NullBool{Bool: false, Valid: true},
			want: false,
			wantErr: false,
		},
		{
			name: "invalid (NULL)",
			nb:   NullBool{Bool: false, Valid: false},
			want: nil,
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.nb.Value()

			if (err != nil) != tt.wantErr {
				t.Errorf("NullBool.Value() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if got != tt.want {
				t.Errorf("NullBool.Value() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNullBool_MarshalJSON(t *testing.T) {
	tests := []struct {
		name    string
		nb      NullBool
		want    string
		wantErr bool
	}{
		{
			name:    "valid true",
			nb:      NullBool{Bool: true, Valid: true},
			want:    "true",
			wantErr: false,
		},
		{
			name:    "valid false",
			nb:      NullBool{Bool: false, Valid: true},
			want:    "false",
			wantErr: false,
		},
		{
			name:    "invalid (NULL)",
			nb:      NullBool{Bool: false, Valid: false},
			want:    "null",
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.nb.MarshalJSON()

			if (err != nil) != tt.wantErr {
				t.Errorf("NullBool.MarshalJSON() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if string(got) != tt.want {
				t.Errorf("NullBool.MarshalJSON() = %s, want %s", string(got), tt.want)
			}
		})
	}
}

func TestNullBool_UnmarshalJSON(t *testing.T) {
	tests := []struct {
		name      string
		input     string
		wantBool  bool
		wantValid bool
		wantErr   bool
	}{
		{
			name:      "true",
			input:     "true",
			wantBool:  true,
			wantValid: true,
			wantErr:   false,
		},
		{
			name:      "false",
			input:     "false",
			wantBool:  false,
			wantValid: true,
			wantErr:   false,
		},
		{
			name:      "null",
			input:     "null",
			wantBool:  false,
			wantValid: false,
			wantErr:   false,
		},
		{
			name:      "invalid json",
			input:     "notabool",
			wantBool:  false,
			wantValid: false,
			wantErr:   true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var nb NullBool
			err := json.Unmarshal([]byte(tt.input), &nb)

			if (err != nil) != tt.wantErr {
				t.Errorf("NullBool.UnmarshalJSON() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !tt.wantErr {
				if nb.Bool != tt.wantBool {
					t.Errorf("NullBool.UnmarshalJSON() Bool = %v, want %v", nb.Bool, tt.wantBool)
				}

				if nb.Valid != tt.wantValid {
					t.Errorf("NullBool.UnmarshalJSON() Valid = %v, want %v", nb.Valid, tt.wantValid)
				}
			}
		})
	}
}

func TestBoolToInt(t *testing.T) {
	tests := []struct {
		name  string
		input bool
		want  int
	}{
		{
			name:  "true converts to 1",
			input: true,
			want:  1,
		},
		{
			name:  "false converts to 0",
			input: false,
			want:  0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := BoolToInt(tt.input)
			if got != tt.want {
				t.Errorf("BoolToInt(%v) = %d, want %d", tt.input, got, tt.want)
			}
		})
	}
}
