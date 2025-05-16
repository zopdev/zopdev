package models

type Metric struct {
	Point any `json:"points"`
}

func (m *Metric) GetDoubleValue() float64 {
	if m.Point == nil {
		return 0.0
	}

	v, ok := m.Point.(float64)
	if !ok {
		return 0.0
	}

	return v
}

func (m *Metric) GetIns64Value() int64 {
	if m.Point == nil {
		return 0
	}

	v, ok := m.Point.(int64)
	if !ok {
		return 0
	}

	return v
}

func (m *Metric) GetBoolValue() bool {
	if m.Point == nil {
		return false
	}

	v, ok := m.Point.(bool)
	if !ok {
		return false
	}

	return v
}

func (m *Metric) GetStringValue() string {
	if m.Point == nil {
		return ""
	}

	v, ok := m.Point.(string)
	if !ok {
		return ""
	}

	return v
}
