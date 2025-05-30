package models

import "testing"

func TestSettings_Scan(t *testing.T) {
	var s Settings

	settingsJSON := []byte(`{"key1":"value1","key2":2}`)

	err := s.Scan(settingsJSON)
	if err != nil {
		t.Errorf("Scan() error = %v, wantErr nil", err)
	}

	if s["key1"] != "value1" || s["key2"] != float64(2) {
		t.Errorf("Scan() got = %v, want map[key1:value1 key2:2]", s)
	}

	// Test nil value
	var s2 Settings

	err = s2.Scan(nil)
	if err != nil {
		t.Errorf("Scan(nil) error = %v, wantErr nil", err)
	}

	// Test non-byte input
	var s3 Settings

	err = s3.Scan(123)
	if err == nil {
		t.Errorf("Scan(non-bytes) error = nil, wantErr driver.ErrSkip")
	}
}

func TestSettings_Value(t *testing.T) {
	s := Settings{"key1": "value1", "key2": 2}
	val, err := s.Value()

	if err != nil {
		t.Errorf("Value() error = %v, wantErr nil", err)
	}

	expected := []byte(`{"key1":"value1","key2":2}`)
	if string(val.([]byte)) != string(expected) {
		t.Errorf("Value() got = %s, want %s", val, expected)
	}

	// Test empty settings
	var s2 Settings

	val2, err := s2.Value()
	if err != nil {
		t.Errorf("Value() on empty settings error = %v, wantErr nil", err)
	}

	if string(val2.([]byte)) != "null" {
		t.Errorf("Value() on empty settings got = %s, want null", val2)
	}
}
