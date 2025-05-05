package handler

import (
	"context"
	"errors"
	netHTTP "net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
	"gofr.dev/pkg/gofr"
	"gofr.dev/pkg/gofr/http"

	"github.com/zopdev/zopdev/api/applications/service"
	"github.com/zopdev/zopdev/api/applications/store"
)

var (
	errTest = errors.New("service error")
)

func TestHandler_AddApplication(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockService := service.NewMockApplicationService(ctrl)
	handler := New(mockService)

	testCases := []struct {
		name           string
		requestBody    string
		mockBehavior   func()
		expectedStatus int
		expectedError  error
	}{
		{
			name:        "success",
			requestBody: `{"name":"Test Application"}`,
			mockBehavior: func() {
				mockService.EXPECT().
					AddApplication(gomock.Any(), gomock.Any()).
					Return(&store.Application{Name: "Test Application"}, nil)
			},
			expectedStatus: netHTTP.StatusOK,
			expectedError:  nil,
		},
		{
			name:           "missing name",
			requestBody:    `{}`,
			mockBehavior:   func() {},
			expectedStatus: netHTTP.StatusBadRequest,
			expectedError:  http.ErrorInvalidParam{Params: []string{"name"}},
		},
		{
			name:        "service error",
			requestBody: `{"name":"Test Application"}`,
			mockBehavior: func() {
				mockService.EXPECT().
					AddApplication(gomock.Any(), gomock.Any()).
					Return(nil, errTest)
			},
			expectedStatus: netHTTP.StatusInternalServerError,
			expectedError:  errTest,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			tc.mockBehavior()

			// Prepare HTTP request
			req := httptest.NewRequest(netHTTP.MethodPost, "/add", strings.NewReader(tc.requestBody))
			req.Header.Set("Content-Type", "application/json")

			ctx := &gofr.Context{Context: context.Background(), Request: http.NewRequest(req)}

			_, err := handler.AddApplication(ctx)

			if tc.expectedError != nil {
				require.Error(t, err)
				require.Equal(t, tc.expectedError, err)
			} else {
				require.NoError(t, err)
			}
		})
	}
}

func TestHandler_ListApplications(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockService := service.NewMockApplicationService(ctrl)
	handler := New(mockService)

	testCases := []struct {
		name           string
		mockBehavior   func()
		expectedStatus int
		expectedError  error
	}{
		{
			name: "success",
			mockBehavior: func() {
				mockService.EXPECT().
					FetchAllApplications(gomock.Any()).
					Return([]store.Application{
						{Name: "Test Application"},
					}, nil)
			},
			expectedStatus: netHTTP.StatusOK,
			expectedError:  nil,
		},
		{
			name: "service error",
			mockBehavior: func() {
				mockService.EXPECT().
					FetchAllApplications(gomock.Any()).
					Return(nil, errTest)
			},
			expectedStatus: netHTTP.StatusInternalServerError,
			expectedError:  errTest,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			tc.mockBehavior()

			// Prepare HTTP request
			req := httptest.NewRequest(netHTTP.MethodGet, "/list", netHTTP.NoBody)

			ctx := &gofr.Context{Context: context.Background(), Request: http.NewRequest(req)}

			_, err := handler.ListApplications(ctx)

			if tc.expectedError != nil {
				require.Error(t, err)
				require.Equal(t, tc.expectedError, err)
			} else {
				require.NoError(t, err)
			}
		})
	}
}
