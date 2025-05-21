package client

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
	"gofr.dev/pkg/gofr"
	"gofr.dev/pkg/gofr/container"
)

func Test_GetCloudCredentials_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	cont, mocks := container.NewMockContainer(t, container.WithMockHTTPService("cloud-account"))
	ctx := &gofr.Context{
		Container: cont,
	}

	cloudAccID := int64(12345)
	body := []byte(`{"data": {"id": 12345, "name": "Test Cloud Account"}}`)

	resp := generateHTTPResponse(body, http.StatusOK)
	defer resp.Body.Close()

	mocks.HTTPService.EXPECT().Get(ctx, "cloud-accounts/12345/credentials", nil).
		Return(resp, nil)

	credentials, err := GetCloudCredentials(ctx, cloudAccID)
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	// Check if credentials are not nil
	if credentials == nil {
		t.Error("Expected credentials to be not nil")
	}
}

var errService = errors.New("service error")

func Test_GetCloudCredentials_Error(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	cont, mocks := container.NewMockContainer(t, container.WithMockHTTPService("cloud-account"))
	ctx := &gofr.Context{
		Container: cont,
	}

	testCases := []struct {
		name          string
		cloudAccID    int64
		expectedError error
		serviceError  error
		statusCode    int
		body          []byte
	}{
		{
			name:          "Cloud Account Id not present",
			cloudAccID:    0,
			expectedError: errFailedToGetCloudCredentials,
			statusCode:    http.StatusBadRequest,
		},
		{
			name:          "error calling the API",
			cloudAccID:    12345,
			serviceError:  errService,
			expectedError: errService,
			statusCode:    0,
		},
		{
			name:          "invalid response from API",
			cloudAccID:    12345,
			expectedError: errInvalidResponse,
			statusCode:    http.StatusOK,
			body:          []byte(`{"data" : {"id" : 123}`),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			resp := generateHTTPResponse(tc.body, tc.statusCode)
			mocks.HTTPService.EXPECT().Get(ctx, fmt.Sprintf("cloud-accounts/%d/credentials", tc.cloudAccID), nil).
				Return(resp, tc.serviceError)

			credentials, err := GetCloudCredentials(ctx, tc.cloudAccID)
			if err == nil {
				t.Errorf("Expected error %v, got nil", tc.expectedError)
			}

			if tc.expectedError == nil {
				assert.NotNil(t, credentials)
			}

			require.ErrorIs(t, err, tc.expectedError)

			resp.Body.Close()
		})
	}
}

func generateHTTPResponse(body []byte, statusCode int) *http.Response {
	return &http.Response{
		StatusCode: statusCode,
		Body:       io.NopCloser(bytes.NewReader(body)),
	}
}
