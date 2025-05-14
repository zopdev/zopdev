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
	Points []Point `json:"points"`
}

type Value interface {
}

type Point struct {
	Value any
}

func (p *Point) GetDoubleValue() float64 {
	v, ok := p.Value.(float64)
	if !ok {
		return 0
	}

	return v
}
