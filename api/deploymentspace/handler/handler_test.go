package handler

import (
	"context"
	"encoding/json"
	"errors"
	netHTTP "net/http"
	"strings"
	"testing"

	"github.com/gorilla/mux"

	"github.com/stretchr/testify/require"

	"go.uber.org/mock/gomock"

	"github.com/zopdev/zopdev/api/deploymentspace/service"
	"gofr.dev/pkg/gofr"
	"gofr.dev/pkg/gofr/http"
)

var (
	errTest = errors.New("service error")
	errJSON = errors.New("invalid character 'i' looking for beginning of value")
)

func TestHandler_Add(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockService := service.NewMockDeploymentSpaceService(ctrl)
	handler := New(mockService)

	deploymentSpace := service.DeploymentSpace{
		CloudAccount: service.CloudAccount{
			ID:   1,
			Name: "test-cloud-account",
		},
		Type: service.Type{
			Name: "test-type",
		},
	}

	bytes, err := json.Marshal(&deploymentSpace)
	if err != nil {
		return
	}

	testCases := []struct {
		name           string
		pathParam      string
		requestBody    string
		mockBehavior   func()
		expectedStatus int
		expectedError  error
	}{
		{
			name:        "success",
			pathParam:   "123",
			requestBody: string(bytes),
			mockBehavior: func() {
				mockService.EXPECT().
					Add(gomock.Any(), &deploymentSpace, 123).
					Return(&service.DeploymentSpace{
						Type: service.Type{Name: "test-type"},
						CloudAccount: service.CloudAccount{
							Name: "test-cloud-account",
							ID:   1,
						},
					}, nil)
			},
			expectedStatus: netHTTP.StatusOK,
			expectedError:  nil,
		},
		{
			name:           "error binding request body",
			pathParam:      "123",
			requestBody:    `invalid-json`,
			mockBehavior:   func() {},
			expectedStatus: netHTTP.StatusBadRequest,
			expectedError:  errJSON,
		},
		{
			name:        "service layer error",
			pathParam:   "123",
			requestBody: string(bytes),
			mockBehavior: func() {
				mockService.EXPECT().
					Add(gomock.Any(), &deploymentSpace, 123).
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
			req, _ := netHTTP.NewRequestWithContext(context.Background(), netHTTP.MethodPost, "/add/{id}", strings.NewReader(tc.requestBody))
			req = mux.SetURLVars(req, map[string]string{"id": tc.pathParam})
			req.Header.Set("Content-Type", "application/json")

			// Add path parameter to the request
			ctx := &gofr.Context{
				Context: context.Background(),
				Request: http.NewRequest(req),
			}

			_, err := handler.Add(ctx)

			if tc.expectedError != nil {
				require.Error(t, err)
				require.Equal(t, tc.expectedError.Error(), err.Error())
			} else {
				require.NoError(t, err)
			}
		})
	}
}
