package common

import (
	"database/sql"
	"encoding/json"
	"testing"
	"time"
)

func TestNullTime_MarshalJSON(t *testing.T) {
	tests := []struct {
		name     string
		nullTime NullTime
		expected string
	}{
		{
			name:     "Valid time",
			nullTime: NullTime{sql.NullTime{Time: time.Date(2023, 12, 25, 15, 30, 0, 0, time.UTC), Valid: true}},
			expected: `"2023-12-25T15:30:00Z"`,
		},
		{
			name:     "Null time",
			nullTime: NullTime{sql.NullTime{Time: time.Time{}, Valid: false}},
			expected: `null`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := json.Marshal(tt.nullTime)
			if err != nil {
				t.Errorf("MarshalJSON() error = %v", err)
				return
			}
			if string(result) != tt.expected {
				t.Errorf("MarshalJSON() = %v, want %v", string(result), tt.expected)
			}
		})
	}
}

func TestNullTime_UnmarshalJSON(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected NullTime
		wantErr  bool
	}{
		{
			name:     "Valid time string",
			input:    `"2023-12-25T15:30:00Z"`,
			expected: NullTime{sql.NullTime{Time: time.Date(2023, 12, 25, 15, 30, 0, 0, time.UTC), Valid: true}},
			wantErr:  false,
		},
		{
			name:     "Null string",
			input:    `null`,
			expected: NullTime{sql.NullTime{Time: time.Time{}, Valid: false}},
			wantErr:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var result NullTime
			err := json.Unmarshal([]byte(tt.input), &result)
			if (err != nil) != tt.wantErr {
				t.Errorf("UnmarshalJSON() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if result.Valid != tt.expected.Valid {
				t.Errorf("UnmarshalJSON() Valid = %v, want %v", result.Valid, tt.expected.Valid)
			}
			if result.Valid && !result.Time.Equal(tt.expected.Time) {
				t.Errorf("UnmarshalJSON() Time = %v, want %v", result.Time, tt.expected.Time)
			}
		})
	}
}



