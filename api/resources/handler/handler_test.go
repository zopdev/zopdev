package handler

import (
	"bytes"
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
	"gofr.dev/pkg/gofr"
	gofrHttp "gofr.dev/pkg/gofr/http"

	"github.com/zopdev/zopdev/api/resources/providers/models"
	"github.com/zopdev/zopdev/api/resources/service"
)

var errMock = errors.New("mock error")

func TestHandler_GetCloudSQLInstances(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	req := httptest.NewRequest(http.MethodPost, "/cloud/sql",
		bytes.NewBufferString(`{"cloudType": "GCP", "creds": {"projectId": "test-project"}}`))
	req.Header.Set("content-type", "application/json")

	mockSvc := NewMockService(ctrl)
	// to test the binding of the request body
	mockReq := service.Request{
		CloudType: "GCP",
		Creds:     map[string]any{"projectId": "test-project"},
	}
	ctx := &gofr.Context{
		Context: context.Background(),
		Request: gofrHttp.NewRequest(req),
	}
	mockResp := []models.Instance{
		{Name: "sql-instance-1"}, {Name: "sql-instance-2"},
	}

	// Test success case
	mockSvc.EXPECT().GetAllSQLInstances(ctx, mockReq).Return(mockResp, nil)

	h := New(mockSvc)
	resp, err := h.GetCloudSQLInstances(ctx)
	require.NoError(t, err)
	assert.Equal(t, mockResp, resp)

	// Test error case
	mockSvc.EXPECT().GetAllSQLInstances(ctx, mockReq).Return(nil, errMock)
	_, err = h.GetCloudSQLInstances(ctx)
	require.Error(t, err)
}
