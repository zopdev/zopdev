package client

import (
	"bytes"
	"errors"
	"fmt"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
	"io"
	"net/http"
	"testing"

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

	cloudAccId := int64(12345)
	body := []byte(`{"data": {"id": 12345, "name": "Test Cloud Account"}}`)

	resp := generateHTTPResponse(body, http.StatusOK)
	defer resp.Body.Close()

	mocks.HTTPService.EXPECT().Get(ctx, "cloud-accounts/12345/credentials", nil).
		Return(resp, nil)

	credentials, err := GetCloudCredentials(ctx, cloudAccId)
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	// Check if credentials are not nil
	if credentials == nil {
		t.Error("Expected credentials to be not nil")
	}
}

var serviceError = errors.New("service error")

func Test_GetCloudCredentials_Error(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	cont, mocks := container.NewMockContainer(t, container.WithMockHTTPService("cloud-account"))
	ctx := &gofr.Context{
		Container: cont,
	}

	testCases := []struct {
		name          string
		cloudAccId    int64
		expectedError error
		serviceError  error
		statusCode    int
		body          []byte
	}{
		{
			name:          "Cloud Account Id not present",
			cloudAccId:    0,
			expectedError: errFailedToGetCloudCredentials,
			statusCode:    http.StatusBadRequest,
		},
		{
			name:          "error calling the API",
			cloudAccId:    12345,
			serviceError:  serviceError,
			expectedError: serviceError,
			statusCode:    0,
		},
		{
			name:          "invalid response from API",
			cloudAccId:    12345,
			expectedError: errInvalidResponse,
			statusCode:    http.StatusOK,
			body:          []byte(`{"data" : {"id" : 123}`),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			resp := generateHTTPResponse(tc.body, tc.statusCode)
			mocks.HTTPService.EXPECT().Get(ctx, fmt.Sprintf("cloud-accounts/%d/credentials", tc.cloudAccId), nil).
				Return(resp, tc.serviceError)

			credentials, err := GetCloudCredentials(ctx, tc.cloudAccId)
			if err == nil {
				t.Errorf("Expected error %v, got nil", tc.expectedError)
			}

			if tc.expectedError == nil {
				assert.NotNil(t, credentials)
			}

			assert.ErrorIs(t, err, tc.expectedError)

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
