package handler

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
	"gofr.dev/pkg/gofr"
	gofrHttp "gofr.dev/pkg/gofr/http"

	"github.com/zopdev/zopdev/api/resources/providers/models"
	"github.com/zopdev/zopdev/api/resources/service"
)

var errMock = errors.New("mock error")

func TestHandler_GetResources(t *testing.T) {
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
		typeQuery    string
		expectedErr  error
		expectedResp any
		mockCall     func()
	}{
		{
			name:         "valid request",
			typeQuery:    "cloudAccId=1&type=sql",
			expectedResp: mockResp,
			mockCall: func() {
				mockSvc.EXPECT().GetResources(ctx, int64(1), []string{"sql"}).
					Return(mockResp, nil)
			},
		},
		{
			name:         "multiple resource types",
			typeQuery:    "cloudAccId=1&type=sql&type=redis",
			expectedResp: mockResp,
			mockCall: func() {
				mockSvc.EXPECT().GetResources(ctx, int64(1), []string{"sql", "redis"}).
					Return(mockResp, nil)
			},
		},
		{
			name:        "error in service",
			typeQuery:   "cloudAccId=1&type=sql",
			expectedErr: errMock,
			mockCall: func() {
				mockSvc.EXPECT().GetResources(ctx, int64(1), []string{"sql"}).
					Return(nil, errMock)
			},
		},
		{
			name:        "invalid id",
			typeQuery:   "cloudAccId=a&type=sql",
			expectedErr: gofrHttp.ErrorInvalidParam{Params: []string{"id"}},
			mockCall:    func() {},
		},
		{
			name:        "missing id",
			typeQuery:   "type=sql",
			expectedErr: gofrHttp.ErrorMissingParam{Params: []string{"id"}},
			mockCall:    func() {},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			tc.mockCall()

			req := httptest.NewRequest(http.MethodGet,
				fmt.Sprintf(`/cloud-account/resources?%s`, tc.typeQuery), http.NoBody)
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
	resDetails := service.ResourceDetails{CloudAccID: 123, Name: "sql-instance-1", Type: "sql", State: service.START}
	h := New(mockSvc)

	testCases := []struct {
		name     string
		reqBody  string
		expErr   error
		mockCall func()
	}{
		{
			name:    "Success",
			reqBody: `{"cloudAccID": 123, "name": "sql-instance-1", "type": "sql", "state": "START"}`,
			mockCall: func() {
				mockSvc.EXPECT().ChangeState(ctx, resDetails).Return(nil)
			},
		},
		{
			name:    "Error from service",
			reqBody: `{"cloudAccID": 123, "name": "sql-instance-1", "type": "sql", "state" : "START"}`,
			expErr:  errMock,
			mockCall: func() {
				mockSvc.EXPECT().ChangeState(ctx, resDetails).Return(errMock)
			},
		},
		{
			name:     "Invalid request body",
			reqBody:  `"cloudAccID": 123, "name": "sql-instance-1"}`,
			expErr:   gofrHttp.ErrorInvalidParam{Params: []string{"request body"}},
			mockCall: func() {},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			tc.mockCall()

			req := httptest.NewRequest(http.MethodPost, "/cloud-account/1/resources/start",
				bytes.NewBufferString(tc.reqBody))
			req.Header.Set("content-type", "application/json")

			ctx.Request = gofrHttp.NewRequest(req)

			resp, err := h.ChangeState(ctx)

			if tc.expErr != nil {
				assert.Equal(t, tc.expErr, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, resDetails, resp)
			}
		})
	}
}
