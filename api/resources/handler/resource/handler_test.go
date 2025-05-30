package resource

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"github.com/gorilla/mux"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
	"gofr.dev/pkg/gofr"
	gofrHttp "gofr.dev/pkg/gofr/http"

	"github.com/zopdev/zopdev/api/resources/models"
	"github.com/zopdev/zopdev/api/resources/service/resource"
)

var errMock = errors.New("mock error")

func TestHandler_GetResources(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockSvc := NewMockService(ctrl)
	ctx := &gofr.Context{
		Context: context.Background(),
	}
	mockResp := []models.Resource{
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
			typeQuery:    "type=sql",
			id:           "1",
			expectedResp: mockResp,
			mockCall: func() {
				mockSvc.EXPECT().GetAll(ctx, int64(1), []string{"sql"}).
					Return(mockResp, nil)
			},
		},
		{
			name:         "multiple resource types",
			typeQuery:    "type=sql&type=redis",
			id:           "1",
			expectedResp: mockResp,
			mockCall: func() {
				mockSvc.EXPECT().GetAll(ctx, int64(1), []string{"sql", "redis"}).
					Return(mockResp, nil)
			},
		},
		{
			name:        "error in service",
			typeQuery:   "type=sql",
			id:          "1",
			expectedErr: errMock,
			mockCall: func() {
				mockSvc.EXPECT().GetAll(ctx, int64(1), []string{"sql"}).
					Return(nil, errMock)
			},
		},
		{
			name:        "invalid id",
			typeQuery:   "type=sql",
			id:          "a",
			expectedErr: gofrHttp.ErrorInvalidParam{Params: []string{"id"}},
			mockCall:    func() {},
		},
		{
			name:        "missing id",
			typeQuery:   "type=sql",
			id:          "",
			expectedErr: gofrHttp.ErrorMissingParam{Params: []string{"id"}},
			mockCall:    func() {},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			tc.mockCall()

			req := httptest.NewRequest(http.MethodGet,
				fmt.Sprintf(`/cloud-account/{id}/resources?%s`, tc.typeQuery), http.NoBody)
			req = mux.SetURLVars(req, map[string]string{"id": tc.id})
			req.Header.Set("content-type", "application/json")

			ctx.Request = gofrHttp.NewRequest(req)

			resp, err := h.GetResources(ctx)

			assert.Equal(t, tc.expectedErr, err)
			assert.Equal(t, tc.expectedResp, resp)
		})
	}
}

func TestHandler_ChangeState(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockSvc := NewMockService(ctrl)
	ctx := &gofr.Context{
		Context: context.Background(),
	}
	resDetails := resource.ResourceDetails{CloudAccID: 123, Name: "sql-instance-1", Type: "sql", State: resource.START}
	h := New(mockSvc)

	testCases := []struct {
		name      string
		pathParam string
		reqBody   string
		expErr    error
		mockCall  func()
	}{
		{
			name:      "Success",
			pathParam: "123",
			reqBody:   `{"name": "sql-instance-1", "type": "sql", "state": "START"}`,
			mockCall: func() {
				mockSvc.EXPECT().ChangeState(ctx, resDetails).Return(nil)
			},
		},
		{
			name:      "Error from service",
			pathParam: "123",
			reqBody:   `{"name": "sql-instance-1", "type": "sql", "state" : "START"}`,
			expErr:    errMock,
			mockCall: func() {
				mockSvc.EXPECT().ChangeState(ctx, resDetails).Return(errMock)
			},
		},
		{
			name:     "Invalid request body",
			reqBody:  `""name": "sql-instance-1"}`,
			expErr:   gofrHttp.ErrorInvalidParam{Params: []string{"request body"}},
			mockCall: func() {},
		},
		{
			name:      "Invalid id",
			pathParam: "a",
			reqBody:   `{"name": "sql-instance-1", "type": "sql", "state": "START"}`,
			expErr:    gofrHttp.ErrorInvalidParam{Params: []string{"id"}},
			mockCall:  func() {},
		},
		{
			name:      "Missing id",
			pathParam: "",
			reqBody:   `{"name": "sql-instance-1", "type": "sql", "state": "START"}`,
			expErr:    gofrHttp.ErrorMissingParam{Params: []string{"id"}},
			mockCall:  func() {},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			tc.mockCall()

			req := httptest.NewRequest(http.MethodPost, "/cloud-account/1/resources/start",
				bytes.NewBufferString(tc.reqBody))
			req = mux.SetURLVars(req, map[string]string{"id": tc.pathParam})
			req.Header.Set("content-type", "application/json")

			ctx.Request = gofrHttp.NewRequest(req)

			resp, err := h.ChangeState(ctx)

			if tc.expErr != nil {
				assert.Equal(t, tc.expErr, err)
			} else {
				require.NoError(t, err)
				assert.Equal(t, resDetails, resp)
			}
		})
	}
}

func TestHandler_SyncResources(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockSvc := NewMockService(ctrl)
	ctx := &gofr.Context{
		Context: context.Background(),
	}
	mockResp := []models.Resource{
		{Name: "sql-instance-1"}, {Name: "sql-instance-2"},
	}
	h := New(mockSvc)

	testCases := []struct {
		name        string
		pathParam   string
		expectedErr error
		expectedRes any
		mockCall    func()
	}{
		{
			name:        "Success",
			pathParam:   "123",
			expectedRes: mockResp,
			mockCall: func() {
				mockSvc.EXPECT().SyncResources(ctx, int64(123)).Return(mockResp, nil)
			},
		},
		{
			name:        "Service error",
			pathParam:   "123",
			expectedErr: errMock,
			mockCall: func() {
				mockSvc.EXPECT().SyncResources(ctx, int64(123)).Return(nil, errMock)
			},
		},
		{
			name:        "Invalid id",
			pathParam:   "a",
			expectedErr: gofrHttp.ErrorInvalidParam{Params: []string{"id"}},
			mockCall:    func() {},
		},
		{
			name:        "Missing id",
			pathParam:   "",
			expectedErr: gofrHttp.ErrorMissingParam{Params: []string{"id"}},
			mockCall:    func() {},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			tc.mockCall()

			req := httptest.NewRequest(http.MethodPost, "/cloud-account/{id}/resources/sync", http.NoBody)
			req = mux.SetURLVars(req, map[string]string{"id": tc.pathParam})
			req.Header.Set("content-type", "application/json")
			ctx.Request = gofrHttp.NewRequest(req)

			resp, err := h.SyncResources(ctx)

			assert.Equal(t, tc.expectedErr, err)
			assert.Equal(t, tc.expectedRes, resp)
		})
	}
}
