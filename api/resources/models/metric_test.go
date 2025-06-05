package models

import (
	"testing"
	
	"github.com/stretchr/testify/assert"
)

func Test_Metric(t *testing.T) {
	tests := []struct {
		name     string
		point    any
		expected any
	}{
		{name: "int64", point: int64(123), expected: int64(123)},
		{name: "int64", point: int32(123), expected: int64(0)},
		{name: "int64", point: nil, expected: int64(0)},
		{name: "float64", point: 123.45, expected: 123.45},
		{name: "float64", point: 123, expected: 0.0},
		{name: "float64", point: nil, expected: 0.0},
		{name: "bool", point: true, expected: true},
		{name: "bool", point: 12, expected: false},
		{name: "bool", point: nil, expected: false},
		{name: "string", point: "test", expected: "test"},
		{name: "string", point: 12, expected: ""},
		{name: "string", point: nil, expected: ""},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			m := Metric{Point: test.point}

			switch test.name {
			case "int64":
				assert.Equal(t, test.expected, m.GetIns64Value())
			case "float64":
				assert.InDelta(t, test.expected, m.GetDoubleValue(), 0.0001)
			case "bool":
				assert.Equal(t, test.expected, m.GetBoolValue())
			case "string":
				assert.Equal(t, test.expected, m.GetStringValue())
			}
		})
	}
}
