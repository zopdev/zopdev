package handler

import (
	"bytes"
	"context"
	"errors"
	"github.com/stretchr/testify/assert"
	"github.com/zopdev/zopdev/api/resources/providers/models"
	"github.com/zopdev/zopdev/api/resources/service"
	"net/http"
	"net/http/httptest"
	"testing"

	"go.uber.org/mock/gomock"
	"gofr.dev/pkg/gofr"
	gofrHttp "gofr.dev/pkg/gofr/http"
)

func TestHandler_GetCloudSQLInstances(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	req := httptest.NewRequest(http.MethodPost, "/cloud/sql", bytes.NewBuffer([]byte(`{"cloudType": 0, "creds": {"projectId": "test-project"}}`)))
	req.Header.Set("content-type", "application/json")

	mockSvc := NewMockService(ctrl)
	// to test the binding of the request body
	mockReq := service.Request{
		CloudType: 0,
		Creds:     map[string]any{"projectId": "test-project"},
	}
	ctx := &gofr.Context{
		Context: context.Background(),
		Request: gofrHttp.NewRequest(req),
	}
	mockResp := []models.SQLInstance{
		{Name: "sql-instance-1"}, {Name: "sql-instance-2"},
	}
	errMock := errors.New("error")

	// Test success case
	mockSvc.EXPECT().GetAllSQLInstances(ctx, mockReq).Return(mockResp, nil)

	h := New(mockSvc)
	resp, err := h.GetCloudSQLInstances(ctx)
	assert.NoError(t, err)
	assert.Equal(t, mockResp, resp)

	// Test error case
	mockSvc.EXPECT().GetAllSQLInstances(ctx, mockReq).Return(nil, errMock)
	resp, err = h.GetCloudSQLInstances(ctx)
	assert.Error(t, err)

}
