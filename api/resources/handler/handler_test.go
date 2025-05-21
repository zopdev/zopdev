package handler

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
	"gofr.dev/pkg/gofr"
	gofrHttp "gofr.dev/pkg/gofr/http"

	"github.com/zopdev/zopdev/api/resources/providers/models"
)

var errMock = errors.New("mock error")

func TestHandler_GetCloudSQLInstances(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockSvc := NewMockService(ctrl)
	ctx := &gofr.Context{
		Context: context.Background(),
	}
	mockResp := []models.Instance{
		{Name: "sql-instance-1"}, {Name: "sql-instance-2"},
	}
	h := New(mockSvc)

	testCases := []struct {
		name         string
		id           string
		typeQuery    string
		expectedErr  error
		expectedResp any
		mockCall     func()
	}{
		{
			name:         "valid request",
			id:           "1",
			typeQuery:    "type=sql",
			expectedResp: mockResp,
			mockCall: func() {
				mockSvc.EXPECT().GetResources(ctx, int64(1), []string{"sql"}).
					Return(mockResp, nil)
			},
		},
		{
			name:         "multiple resource types",
			id:           "1",
			typeQuery:    "type=sql&type=redis",
			expectedResp: mockResp,
			mockCall: func() {
				mockSvc.EXPECT().GetResources(ctx, int64(1), []string{"sql", "redis"}).
					Return(mockResp, nil)
			},
		},
		{
			name:        "error in service",
			id:          "1",
			typeQuery:   "type=sql",
			expectedErr: errMock,
			mockCall: func() {
				mockSvc.EXPECT().GetResources(ctx, int64(1), []string{"sql"}).
					Return(nil, errMock)
			},
		},
		{
			name:        "invalid id",
			id:          "invalid",
			typeQuery:   "type=sql",
			expectedErr: gofrHttp.ErrorInvalidParam{Params: []string{"id"}},
			mockCall:    func() {},
		},
		{
			name:        "missing id",
			id:          "",
			typeQuery:   "type=sql",
			expectedErr: gofrHttp.ErrorMissingParam{Params: []string{"id"}},
			mockCall:    func() {},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			tc.mockCall()

			req := httptest.NewRequest(http.MethodGet,
				fmt.Sprintf(`/cloud-account/1/resources?%s`, tc.typeQuery), http.NoBody)
			req = mux.SetURLVars(req, map[string]string{"id": tc.id})
			req.Header.Set("content-type", "application/json")

			ctx.Request = gofrHttp.NewRequest(req)

			resp, err := h.GetResources(ctx)

			assert.Equal(t, tc.expectedErr, err)
			assert.Equal(t, tc.expectedResp, resp)
		})
	}
}
