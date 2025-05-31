package service

import (
	"context"
	"testing"

	gofrHttp "gofr.dev/pkg/gofr/http"

	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
	"gofr.dev/pkg/gofr"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/option"

	"github.com/zopdev/zopdev/api/resources/client"
)

func TestService_changeSQLState(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mGCP := NewMockGCPClient(ctrl)
	ctx := &gofr.Context{Context: context.Background()}
	ca := &client.CloudAccount{ID: 123, Name: "MyCloud", Provider: string(GCP),
		Credentials: map[string]any{"project_id": "test-project", "region": "us-central1"}}
	mockCreds := &google.Credentials{ProjectID: "test-project"}
	mockStopper := &mockSQLClient{}
	s := New(mGCP, nil, nil, nil)

	testCases := []struct {
		name      string
		input     ResourceDetails
		cloudAcc  *client.CloudAccount
		expErr    error
		mockCalls func()
	}{
		{
			name:  "Successfully start SQL instance",
			input: ResourceDetails{CloudAccID: 123, Name: "test-instance", Type: SQL, State: START},
			mockCalls: func() {
				mGCP.EXPECT().NewGoogleCredentials(ctx, gomock.Any(), "https://www.googleapis.com/auth/cloud-platform").
					Return(mockCreds, nil)
				mGCP.EXPECT().NewSQLClient(ctx, option.WithCredentials(mockCreds)).
					Return(mockStopper, nil)
			},
		},
		{
			name:  "Successfully stop SQL instance",
			input: ResourceDetails{CloudAccID: 123, Name: "test-instance", Type: SQL, State: SUSPEND},
			mockCalls: func() {
				mGCP.EXPECT().NewGoogleCredentials(ctx, gomock.Any(), "https://www.googleapis.com/auth/cloud-platform").
					Return(mockCreds, nil)
				mGCP.EXPECT().NewSQLClient(ctx, option.WithCredentials(mockCreds)).
					Return(mockStopper, nil)
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			tc.mockCalls()

			err := s.changeSQLState(ctx, ca, tc.input)

			assert.Equal(t, tc.expErr, err)
		})
	}
}

func TestService_changeSQLState_Errors(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mGCP := NewMockGCPClient(ctrl)
	ctx := &gofr.Context{Context: context.Background()}
	ca := &client.CloudAccount{ID: 123, Name: "MyCloud", Provider: string(GCP),
		Credentials: map[string]any{"project_id": "test-project", "region": "us-central1"}}
	mockCreds := &google.Credentials{ProjectID: "test-project"}
	mockStopper := &mockSQLClient{}
	s := New(mGCP, nil, nil, nil)

	testCases := []struct {
		name      string
		input     ResourceDetails
		cloudAcc  *client.CloudAccount
		expErr    error
		mockCalls func()
	}{
		{
			name:   "Error getting google credentials",
			input:  ResourceDetails{CloudAccID: 123, Name: "test-instance", Type: SQL},
			expErr: errMock,
			mockCalls: func() {
				mGCP.EXPECT().NewGoogleCredentials(ctx, gomock.Any(), "https://www.googleapis.com/auth/cloud-platform").
					Return(nil, errMock)
			},
		},
		{
			name:   "Error getting SQL Client",
			input:  ResourceDetails{CloudAccID: 123, Name: "test-instance", Type: SQL, State: START},
			expErr: errMock,
			mockCalls: func() {
				mGCP.EXPECT().NewGoogleCredentials(ctx, gomock.Any(), "https://www.googleapis.com/auth/cloud-platform").
					Return(mockCreds, nil)
				mGCP.EXPECT().NewSQLClient(ctx, option.WithCredentials(mockCreds)).
					Return(nil, errMock)
			},
		},
		{
			name:   "Error starting SQL instance",
			input:  ResourceDetails{CloudAccID: 123, Name: "test-instance", Type: SQL, State: START},
			expErr: errMock,
			mockCalls: func() {
				mockStopperErr := &mockSQLClient{isError: true}
				mGCP.EXPECT().NewGoogleCredentials(ctx, gomock.Any(), "https://www.googleapis.com/auth/cloud-platform").
					Return(mockCreds, nil)
				mGCP.EXPECT().NewSQLClient(ctx, option.WithCredentials(mockCreds)).
					Return(mockStopperErr, nil)
			},
		},
		{
			name:   "Error stopping SQL instance",
			input:  ResourceDetails{CloudAccID: 123, Name: "test-instance", Type: SQL, State: SUSPEND},
			expErr: errMock,
			mockCalls: func() {
				mockStopperErr := &mockSQLClient{isError: true}
				mGCP.EXPECT().NewGoogleCredentials(ctx, gomock.Any(), "https://www.googleapis.com/auth/cloud-platform").
					Return(mockCreds, nil)
				mGCP.EXPECT().NewSQLClient(ctx, option.WithCredentials(mockCreds)).
					Return(mockStopperErr, nil)
			},
		},
		{
			name:   "Error - invalid state",
			input:  ResourceDetails{CloudAccID: 123, Name: "test-instance", Type: SQL, State: "invalid"},
			expErr: gofrHttp.ErrorInvalidParam{Params: []string{"req.State"}},
			mockCalls: func() {
				mGCP.EXPECT().NewGoogleCredentials(ctx, gomock.Any(), "https://www.googleapis.com/auth/cloud-platform").
					Return(mockCreds, nil)
				mGCP.EXPECT().NewSQLClient(ctx, option.WithCredentials(mockCreds)).
					Return(mockStopper, nil)
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			tc.mockCalls()

			err := s.changeSQLState(ctx, ca, tc.input)

			assert.Equal(t, tc.expErr, err)
		})
	}
}
