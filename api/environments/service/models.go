package service

type DeploymentSpaceResponse struct {
	Name string                   `json:"name"`
	Next *DeploymentSpaceResponse `json:"next,omitempty"`
}
