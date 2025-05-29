package resourcegroup

import "gofr.dev/pkg/gofr"

type Service struct {
}

func New() *Service {
	return &Service{}
}

func (s *Service) GetAllResourceGroups(ctx *gofr.Context) ([]string, error) {
	// Placeholder for getting all resource groups
	return []string{"group1", "group2"}, nil
}
