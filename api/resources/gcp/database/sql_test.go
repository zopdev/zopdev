package sql

import (
	"testing"

	googleapi "google.golang.org/api/googleapi"
	sqladmin "google.golang.org/api/sqladmin/v1beta4"
)

type mockCaller struct {
	showError bool
}

func (m *mockCaller) Do(opts ...googleapi.CallOption) (*sqladmin.InstancesListResponse, error) {
	if m.showError {
		return nil, &googleapi.Error{Code: 500, Message: "Internal Server Error"}
	}

	return &sqladmin.InstancesListResponse{
		Items: []*sqladmin.DatabaseInstance{
			{
				Name:            "test-instance",
				Region:          "us-central1",
				DatabaseVersion: "MYSQL_5_7",
			},
		},
	}, nil
}

func TestGetAllInstance_Success(t *testing.T) {
	getCaller := func(projectId string) Caller {
		return &mockCaller{}
	}

	projectId := "test-project"
	instances, err := GetAllInstance(getCaller, projectId)
	if err != nil {
		t.Errorf("expected no error, got %v", err)
	}

	if len(instances) == 0 {
		t.Error("expected at least one instance, got none")
	}

	for _, instance := range instances {
		if instance.InstanceName == "" {
			t.Error("expected instance name to be non-empty")
		}
		if instance.Region == "" {
			t.Error("expected region to be non-empty")
		}
		if instance.InstanceType == "" {
			t.Error("expected instance type to be non-empty")
		}
	}
}

func TestGetAllInstance_Error(t *testing.T) {
	getCaller := func(projectId string) Caller {
		return &mockCaller{showError: true}
	}

	projectId := "test-project"
	instances, err := GetAllInstance(getCaller, projectId)
	if err == nil {
		t.Error("expected error, got none")
	}

	if len(instances) != 0 {
		t.Error("expected no instances, got some")
	}
}
