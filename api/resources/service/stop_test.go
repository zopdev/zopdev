package service

import (
	"bytes"
	"io"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
	gofrHttp "gofr.dev/pkg/gofr/http"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/option"
)

func TestService_Stop(t *testing.T) {
	ctx, ctrl, mock, mGCP := InitializeTests(t)
	defer ctrl.Finish()

	s := New(mGCP)
	mockCreds := &google.Credentials{ProjectID: "test-project"}
	mockStopper := &mockSQLClient{}

	testCases := []struct {
		name      string
		input     ResourceDetails
		expErr    error
		mockCalls func()
	}{
		{
			name:  "Successfully stop SQL instance",
			input: ResourceDetails{CloudAccID: 123, Name: "test-instance", Type: SQL},
			mockCalls: func() {
				resp := &http.Response{StatusCode: http.StatusOK, Body: io.NopCloser(bytes.NewReader([]byte(
					`{"data": {"id": 123, "provider" : "GCP", 
		"credentials" : {"project_id": "test-project", "region": "us-central1"}}}`))),
				}

				mock.HTTPService.EXPECT().Get(ctx, "cloud-accounts/123/credentials", nil).
					Return(resp, nil)
				mGCP.EXPECT().NewGoogleCredentials(ctx, gomock.Any(), "https://www.googleapis.com/auth/cloud-platform").
					Return(mockCreds, nil)
				mGCP.EXPECT().NewSQLClient(ctx, option.WithCredentials(mockCreds)).
					Return(mockStopper, nil)
			},
		},
		{
			name:   "Invalid resource type",
			input:  ResourceDetails{CloudAccID: 123, Name: "test-instance", Type: "invalid-type"},
			expErr: gofrHttp.ErrorInvalidParam{Params: []string{"type"}},
			mockCalls: func() {
				resp := &http.Response{StatusCode: http.StatusOK, Body: io.NopCloser(bytes.NewReader([]byte(
					`{"data": {"id": 123, "provider" : "GCP", 
		"credentials" : {"project_id": "test-project", "region": "us-central1"}}}`))),
				}

				mock.HTTPService.EXPECT().Get(ctx, "cloud-accounts/123/credentials", nil).
					Return(resp, nil)
			},
		},
		{
			name:   "Error getting credentials",
			input:  ResourceDetails{CloudAccID: 123, Name: "test-instance", Type: SQL},
			expErr: errMock,
			mockCalls: func() {
				mock.HTTPService.EXPECT().Get(ctx, "cloud-accounts/123/credentials", nil).
					Return(nil, errMock)
			},
		},
		{
			name:   "Error getting SQL Client",
			input:  ResourceDetails{CloudAccID: 123, Name: "test-instance", Type: SQL},
			expErr: errMock,
			mockCalls: func() {
				resp := &http.Response{StatusCode: http.StatusOK, Body: io.NopCloser(bytes.NewReader([]byte(
					`{"data": {"id": 123, "provider" : "GCP", 
		"credentials" : {"project_id": "test-project", "region": "us-central1"}}}`))),
				}

				mock.HTTPService.EXPECT().Get(ctx, "cloud-accounts/123/credentials", nil).
					Return(resp, nil)
				mGCP.EXPECT().NewGoogleCredentials(ctx, gomock.Any(), "https://www.googleapis.com/auth/cloud-platform").
					Return(mockCreds, nil)
				mGCP.EXPECT().NewSQLClient(ctx, option.WithCredentials(mockCreds)).
					Return(nil, errMock)
			},
		},
		{
			name:   "Error getting google credentials",
			input:  ResourceDetails{CloudAccID: 123, Name: "test-instance", Type: SQL},
			expErr: errMock,
			mockCalls: func() {
				resp := &http.Response{StatusCode: http.StatusOK, Body: io.NopCloser(bytes.NewReader([]byte(
					`{"data": {"id": 123, "provider" : "GCP", 
		"credentials" : {"project_id": "test-project", "region": "us-central1"}}}`))),
				}

				mock.HTTPService.EXPECT().Get(ctx, "cloud-accounts/123/credentials", nil).
					Return(resp, nil)
				mGCP.EXPECT().NewGoogleCredentials(ctx, gomock.Any(), "https://www.googleapis.com/auth/cloud-platform").
					Return(nil, errMock)
			},
		},
		{
			name:   "Error starting SQL instance",
			input:  ResourceDetails{CloudAccID: 123, Name: "test-instance", Type: SQL},
			expErr: errMock,
			mockCalls: func() {
				mockStopperErr := &mockSQLClient{isError: true}
				resp := &http.Response{StatusCode: http.StatusOK, Body: io.NopCloser(bytes.NewReader([]byte(
					`{"data": {"id": 123, "provider" : "GCP", 
		"credentials" : {"project_id": "test-project", "region": "us-central1"}}}`))),
				}

				mock.HTTPService.EXPECT().Get(ctx, "cloud-accounts/123/credentials", nil).
					Return(resp, nil)
				mGCP.EXPECT().NewGoogleCredentials(ctx, gomock.Any(), "https://www.googleapis.com/auth/cloud-platform").
					Return(mockCreds, nil)
				mGCP.EXPECT().NewSQLClient(ctx, option.WithCredentials(mockCreds)).
					Return(mockStopperErr, nil)
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			tc.mockCalls()

			err := s.stop(ctx, tc.input)

			assert.Equal(t, tc.expErr, err)
		})
	}
}
