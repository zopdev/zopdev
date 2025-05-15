package models

type SQLInstance struct {
	InstanceName string `json:"instance_name"`
	ProjectID    string `json:"project_id"`
	Region       string `json:"region"`
	Zone         string `json:"zone"`
	Version      string `json:"version"`
	CreationTime string `json:"creation_time"`
}

type Metric struct {
	Point any `json:"points"`
}

func (m *Metric) GetDoubleValue() float64 {
	v, ok := m.Point.(float64)
	if !ok {
		return 0
	}

	return v
}
