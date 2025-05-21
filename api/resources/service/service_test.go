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

	"github.com/zopdev/zopdev/api/resources/providers/models"
)

func TestService_GetResources(t *testing.T) {
	ctx, ctrl, mock, mGCP := InitializeTests(t)
	defer ctrl.Finish()

	s := New(mGCP)
	req := CloudDetails{
		CloudType: GCP,
		Creds: map[string]any{
			"project_id": "test-project",
			"region":     "us-central1",
		},
	}
	mockCreds := &google.Credentials{ProjectID: "test-project"}
	mockResp := []models.Instance{
		{Name: "sql-instance-1"}, {Name: "sql-instance-2"},
	}
	mockLister := &mockSQLClient{
		isError:   false,
		instances: mockResp,
	}

	testCases := []struct {
		name      string
		id        int64
		resources []string
		expErr    error
		expResp   []models.Instance
		mockCalls func()
	}{
		{
			name:      "Get all resources",
			id:        123,
			resources: []string{},
			expResp:   mockResp,
			mockCalls: func() {
				resp := &http.Response{
					StatusCode: http.StatusOK,
					Body: io.NopCloser(
						bytes.NewReader([]byte(
							`{"data": {"id": 123, "provider" : "GCP", "credentials" : {"project_id": "test-project", "region": "us-central1"}}}`),
						)),
				}

				mock.HTTPService.EXPECT().Get(ctx, "cloud-accounts/123/credentials", nil).
					Return(resp, nil)
				mGCP.EXPECT().NewGoogleCredentials(ctx, req.Creds, "https://www.googleapis.com/auth/cloud-platform").
					Return(mockCreds, nil)
				mGCP.EXPECT().NewSQLClient(ctx, option.WithCredentials(mockCreds)).
					Return(mockLister, nil)
			},
		},
		{
			name:      "Get SQL resources",
			id:        123,
			resources: []string{string(SQL)},
			expResp:   mockResp,
			mockCalls: func() {
				resp := &http.Response{
					StatusCode: http.StatusOK,
					Body: io.NopCloser(
						bytes.NewReader([]byte(
							`{"data": {"id": 123, "provider" : "GCP", "credentials" : {"project_id": "test-project", "region": "us-central1"}}}`),
						)),
				}

				mock.HTTPService.EXPECT().Get(ctx, "cloud-accounts/123/credentials", nil).
					Return(resp, nil)
				mGCP.EXPECT().NewGoogleCredentials(ctx, req.Creds, "https://www.googleapis.com/auth/cloud-platform").
					Return(mockCreds, nil)
				mGCP.EXPECT().NewSQLClient(ctx, option.WithCredentials(mockCreds)).
					Return(mockLister, nil)
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			tc.mockCalls()

			res, err := s.GetResources(ctx, tc.id, tc.resources)

			assert.Equal(t, tc.expErr, err)
			assert.Equal(t, tc.expResp, res)
		})
	}
}

func TestService_GetResources_Errors(t *testing.T) {
	ctx, ctrl, mock, mGCP := InitializeTests(t)
	defer ctrl.Finish()

	s := New(mGCP)
	req := CloudDetails{
		CloudType: GCP,
		Creds: map[string]any{
			"project_id": "test-project",
			"region":     "us-central1",
		},
	}
	testCases := []struct {
		name      string
		id        int64
		resources []string
		expErr    error
		expResp   []models.Instance
		mockCalls func()
	}{
		{
			name:      "error getting SQL resources",
			id:        123,
			resources: []string{string(SQL)},
			expErr:    errMock,
			mockCalls: func() {
				resp := &http.Response{
					StatusCode: http.StatusOK,
					Body: io.NopCloser(
						bytes.NewReader([]byte(
							`{"data": {"id": 123, "provider" : "GCP", "credentials" : {"project_id": "test-project", "region": "us-central1"}}}`),
						)),
				}

				mock.HTTPService.EXPECT().Get(ctx, "cloud-accounts/123/credentials", nil).
					Return(resp, nil)
				mGCP.EXPECT().NewGoogleCredentials(ctx, req.Creds, "https://www.googleapis.com/auth/cloud-platform").
					Return(nil, errMock)
			},
		},
		{
			name:   "error from GetCloudCredentials",
			id:     123,
			expErr: errMock,
			mockCalls: func() {
				mock.HTTPService.EXPECT().Get(ctx, "cloud-accounts/123/credentials", nil).
					Return(nil, errMock)
			},
		},
		{
			name:      "error from GetAllInstances",
			id:        123,
			resources: []string{},
			expErr:    errMock,
			mockCalls: func() {
				resp := &http.Response{
					StatusCode: http.StatusOK,
					Body: io.NopCloser(
						bytes.NewReader([]byte(
							`{"data": {"id": 123, "provider" : "GCP", "credentials" : {"project_id": "test-project", "region": "us-central1"}}}`),
						)),
				}

				mock.HTTPService.EXPECT().Get(ctx, "cloud-accounts/123/credentials", nil).
					Return(resp, nil)

				mGCP.EXPECT().NewGoogleCredentials(ctx, req.Creds, "https://www.googleapis.com/auth/cloud-platform").
					Return(nil, errMock)
			},
		},
		{
			name:      "unsupported resource",
			id:        123,
			resources: []string{"unsupported"},
			expErr:    gofrHttp.ErrorInvalidParam{Params: []string{"req.CloudType"}},
			mockCalls: func() {
				resp := &http.Response{
					StatusCode: http.StatusOK,
					Body: io.NopCloser(
						bytes.NewReader([]byte(
							`{"data": {"id": 123, "provider" : "GCP", "credentials" : {"project_id": "test-project", "region": "us-central1"}}}`),
						)),
				}

				mock.HTTPService.EXPECT().Get(ctx, "cloud-accounts/123/credentials", nil).
					Return(resp, nil)
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			tc.mockCalls()

			res, err := s.GetResources(ctx, tc.id, tc.resources)

			assert.Equal(t, tc.expErr, err)
			assert.Equal(t, tc.expResp, res)
		})
	}
}

func TestService_ChangeState(t *testing.T) {
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
			name:  "Success - Start",
			input: ResourceDetails{CloudAccID: 123, Name: "test-instance", Type: SQL, State: START},
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
			name:  "Success - Stop",
			input: ResourceDetails{CloudAccID: 123, Name: "test-instance", Type: SQL, State: SUSPEND},
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
			name:      "Error - Invalid State",
			input:     ResourceDetails{CloudAccID: 123, Name: "test-instance", Type: SQL, State: "invalid"},
			expErr:    gofrHttp.ErrorInvalidParam{Params: []string{"state"}},
			mockCalls: func() {},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			tc.mockCalls()

			err := s.ChangeState(ctx, tc.input)

			assert.Equal(t, tc.expErr, err)
		})
	}
}
