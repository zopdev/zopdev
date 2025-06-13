package service

import (
	"reflect"

	"go.uber.org/mock/gomock"
	"gofr.dev/pkg/gofr"
)

// GetStackStatus mocks base method.
func (m *MockCloudAccountService) GetStackStatus(ctx *gofr.Context, integrationID string) (string, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetStackStatus", ctx, integrationID)
	ret0, _ := ret[0].(string)
	ret1, _ := ret[1].(error)

	return ret0, ret1
}

// GetStackStatus indicates an expected call of GetStackStatus.
func (mr *MockCloudAccountServiceMockRecorder) GetStackStatus(ctx, integrationID interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetStackStatus",
		reflect.TypeOf((*MockCloudAccountService)(nil).GetStackStatus), ctx, integrationID)
}
