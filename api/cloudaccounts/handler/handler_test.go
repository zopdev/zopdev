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

	"github.com/zopdev/zopdev/api/cloudaccounts/service"
	"github.com/zopdev/zopdev/api/cloudaccounts/store"
)

var (
	errTest = errors.New("service error")
)

func TestHandler_AddCloudAccount(t *testing.T) {
	ctrl := gomock.NewController(t)

	defer ctrl.Finish()

	mockService := service.NewMockCloudAccountService(ctrl)
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
			requestBody: `{"name":"Test Account","provider":"gcp","credentials":{"project_id":"test-project-id"}}`,
			mockBehavior: func() {
				mockService.EXPECT().
					AddCloudAccount(gomock.Any(), gomock.Any()).
					Return(&store.CloudAccount{Name: "Test Account", Provider: "gcp"}, nil)
			},
			expectedStatus: netHTTP.StatusOK,
			expectedError:  nil,
		},
		{
			name:           "missing parameters",
			requestBody:    `{}`,
			mockBehavior:   func() {},
			expectedStatus: netHTTP.StatusBadRequest,
			expectedError:  http.ErrorMissingParam{Params: []string{"name", "provider"}},
		},
		{
			name:           "invalid provider",
			requestBody:    `{"name":"Test Account","provider":"aws","credentials":{}}`,
			mockBehavior:   func() {},
			expectedStatus: netHTTP.StatusBadRequest,
			expectedError:  http.ErrorInvalidParam{Params: []string{"provider"}},
		},
		{
			name:        "service error",
			requestBody: `{"name":"Test Account","provider":"gcp","credentials":{"project_id":"test-project-id"}}`,
			mockBehavior: func() {
				mockService.EXPECT().
					AddCloudAccount(gomock.Any(), gomock.Any()).
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
			req := httptest.NewRequest(netHTTP.MethodPost, "/add/{id}", strings.NewReader(tc.requestBody))
			req.Header.Set("Content-Type", "application/json")

			ctx := &gofr.Context{Context: context.Background(), Request: http.NewRequest(req)}

			_, err := handler.AddCloudAccount(ctx)

			if tc.expectedError != nil {
				require.Error(t, err)
				require.Equal(t, tc.expectedError, err)
			} else {
				require.NoError(t, err)
			}
		})
	}
}

func TestHandler_ListCloudAccounts(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockService := service.NewMockCloudAccountService(ctrl)

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
					FetchAllCloudAccounts(gomock.Any()).
					Return([]store.CloudAccount{
						{Name: "Test Account", Provider: "gcp"},
					}, nil)
			},
			expectedStatus: netHTTP.StatusOK,
			expectedError:  nil,
		},
		{
			name: "service error",
			mockBehavior: func() {
				mockService.EXPECT().
					FetchAllCloudAccounts(gomock.Any()).
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

			_, err := handler.ListCloudAccounts(ctx)

			if tc.expectedError != nil {
				require.Error(t, err)
				require.Equal(t, tc.expectedError, err)
			} else {
				require.NoError(t, err)
			}
		})
	}
}
