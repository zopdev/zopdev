package resource

import (
	"errors"

	"gofr.dev/pkg/gofr"

	"github.com/zopdev/zopdev/api/resources/models"
)

var errMock = errors.New("mock error")

type mockSQLClient struct {
	isError   bool
	instances []models.Resource
}

func (m *mockSQLClient) GetAllInstances(_ *gofr.Context, _ string) ([]models.Resource, error) {
	if m.isError {
		return nil, errMock
	}

	return m.instances, nil
}

func (m *mockSQLClient) StartInstance(_ *gofr.Context, _, _ string) error {
	if m.isError {
		return errMock
	}

	return nil
}

func (m *mockSQLClient) StopInstance(_ *gofr.Context, _, _ string) error {
	if m.isError {
		return errMock
	}

	return nil
}
