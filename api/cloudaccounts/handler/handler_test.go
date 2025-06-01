package handler

import (
	"context"
	"errors"
	"fmt"
	netHTTP "net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gorilla/mux"
	"github.com/stretchr/testify/require"

	"go.uber.org/mock/gomock"

	"gofr.dev/pkg/gofr"
	"gofr.dev/pkg/gofr/http"

	"github.com/zopdev/zopdev/api/cloudaccounts/service"
	"github.com/zopdev/zopdev/api/cloudaccounts/store"
)

var (
	errTest                          = errors.New("service error")
	errMissingIntegrationOrAccountID = errors.New("missing required fields: integration_id or account_id")
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

func TestHandler_GetStackStatus(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockService := service.NewMockCloudAccountService(ctrl)
	handler := New(mockService)

	testCases := []struct {
		name           string
		integrationID  string
		mockBehavior   func()
		expectedStatus string
		expectedError  error
	}{
		{
			name:          "success - stack exists",
			integrationID: "test-integration-123",
			mockBehavior: func() {
				mockService.EXPECT().
					GetStackStatus(gomock.Any(), "Zopdev-test-integration-123").
					Return("CREATE_COMPLETE", nil)
			},
			expectedStatus: "CREATE_COMPLETE",
			expectedError:  nil,
		},
		{
			name:          "missing integration ID",
			integrationID: "",
			mockBehavior: func() {
				// No mock needed
			},
			expectedStatus: "",
			expectedError:  http.ErrorMissingParam{Params: []string{"integrationId"}},
		},
		{
			name:          "stack does not exist",
			integrationID: "non-existent-stack",
			mockBehavior: func() {
				mockService.EXPECT().
					GetStackStatus(gomock.Any(), "Zopdev-non-existent-stack").
					Return("STACK_DOES_NOT_EXIST", nil)
			},
			expectedStatus: "STACK_DOES_NOT_EXIST",
			expectedError:  nil,
		},
		{
			name:          "service error",
			integrationID: "test-integration-123",
			mockBehavior: func() {
				mockService.EXPECT().
					GetStackStatus(gomock.Any(), "Zopdev-test-integration-123").
					Return("", errTest)
			},
			expectedStatus: "",
			expectedError:  errTest,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			tc.mockBehavior()

			req := httptest.NewRequest(netHTTP.MethodGet, fmt.Sprintf("/cloud-accounts/aws/stack-status/%s", tc.integrationID), netHTTP.NoBody)
			if tc.integrationID != "" {
				req = mux.SetURLVars(req, map[string]string{
					"integrationId": tc.integrationID,
				})
			}

			ctx := &gofr.Context{
				Context: context.Background(),
				Request: http.NewRequest(req),
			}

			resp, err := handler.GetStackStatus(ctx)

			if tc.expectedError != nil {
				require.Error(t, err)
				require.Equal(t, tc.expectedError, err)
			} else {
				require.NoError(t, err)

				statusMap, ok := resp.(map[string]string)
				require.True(t, ok, "response should be a map[string]string")
				require.Equal(t, tc.expectedStatus, statusMap["status"])
			}
		})
	}
}

func TestHandler_GetCloudAccountConnectionInfo(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockService := service.NewMockCloudAccountService(ctrl)
	handler := New(mockService)

	testCases := []struct {
		name           string
		provider       string
		mockBehavior   func()
		expectedStatus int
		expectedError  error
	}{
		{
			name:     "success",
			provider: "aws",
			mockBehavior: func() {
				mockService.EXPECT().
					GetCloudAccountConnectionInfo(gomock.Any(), "aws").
					Return(service.AWSIntegrationINFO{
						CloudformationURL: "https://test-url",
						IntegrationID:     "test-id",
					}, nil)
			},
			expectedStatus: netHTTP.StatusOK,
			expectedError:  nil,
		},
		{
			name:           "missing provider",
			provider:       "",
			mockBehavior:   func() {},
			expectedStatus: netHTTP.StatusBadRequest,
			expectedError:  http.ErrorMissingParam{Params: []string{"provider"}},
		},
		{
			name:           "invalid provider",
			provider:       "gcp",
			mockBehavior:   func() {},
			expectedStatus: netHTTP.StatusBadRequest,
			expectedError:  http.ErrorInvalidParam{Params: []string{"provider"}},
		},
		{
			name:     "service error",
			provider: "aws",
			mockBehavior: func() {
				mockService.EXPECT().
					GetCloudAccountConnectionInfo(gomock.Any(), "aws").
					Return(service.AWSIntegrationINFO{}, errTest)
			},
			expectedStatus: netHTTP.StatusInternalServerError,
			expectedError:  errTest,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			tc.mockBehavior()

			req := httptest.NewRequest(netHTTP.MethodGet, fmt.Sprintf("/cloud-accounts/%s/connection-info", tc.provider), netHTTP.NoBody)
			if tc.provider != "" {
				req = mux.SetURLVars(req, map[string]string{
					"provider": tc.provider,
				})
			}

			ctx := &gofr.Context{
				Context: context.Background(),
				Request: http.NewRequest(req),
			}

			_, err := handler.GetCloudAccountConnectionInfo(ctx)

			if tc.expectedError != nil {
				require.Error(t, err)
				require.Equal(t, tc.expectedError, err)
			} else {
				require.NoError(t, err)
			}
		})
	}
}

func TestHandler_CreateCloudAccountConnection(t *testing.T) {
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
			name: "success",
			requestBody: `{"roleArn": "arn:aws:iam::123456789012:role/test-role", "integrationId": "test-id", "accountId":
"123456789012", "cloudAccountName": "Test Account"}`,
			mockBehavior: func() {
				mockService.EXPECT().
					CreateCloudAccountConnection(gomock.Any(), gomock.Any()).
					Return(&store.CloudAccount{
						Name:     "Test Account",
						Provider: "aws",
					}, nil)
			},
			expectedStatus: netHTTP.StatusOK,
			expectedError:  nil,
		},
		{
			name:        "missing required fields",
			requestBody: `{"roleArn": "arn:aws:iam::123456789012:role/test-role"}`,
			mockBehavior: func() {
				mockService.EXPECT().
					CreateCloudAccountConnection(gomock.Any(), gomock.Any()).
					Return(nil, errMissingIntegrationOrAccountID)
			},
			expectedStatus: netHTTP.StatusBadRequest,
			expectedError:  errMissingIntegrationOrAccountID,
		},
		{
			name: "service error",
			requestBody: `{"roleArn": "arn:aws:iam::123456789012:role/test-role", "integrationId": "test-id", "accountId": 
"123456789012", "cloudAccountName": "Test Account"}`,
			mockBehavior: func() {
				mockService.EXPECT().
					CreateCloudAccountConnection(gomock.Any(), gomock.Any()).
					Return(nil, errTest)
			},
			expectedStatus: netHTTP.StatusInternalServerError,
			expectedError:  errTest,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			tc.mockBehavior()

			req := httptest.NewRequest(netHTTP.MethodPost, "/cloud-accounts/connection", strings.NewReader(tc.requestBody))
			req.Header.Set("Content-Type", "application/json")

			ctx := &gofr.Context{
				Context: context.Background(),
				Request: http.NewRequest(req),
			}

			_, err := handler.CreateCloudAccountConnection(ctx)

			if tc.expectedError != nil {
				require.Error(t, err)
				require.Equal(t, tc.expectedError, err)
			} else {
				require.NoError(t, err)
			}
		})
	}
}
