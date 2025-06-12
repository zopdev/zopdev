package resource

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
	"gofr.dev/pkg/gofr"

	"github.com/zopdev/zopdev/api/resources/client"
)

func TestService_changeSQLState(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mGCP := NewMockCloudResourceProvider(ctrl)
	ctx := &gofr.Context{Context: context.Background()}
	ca := &client.CloudAccount{ID: 123, Name: "MyCloud", Provider: string(GCP),
		Credentials: map[string]any{"project_id": "test-project", "region": "us-central1"}}
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
				mGCP.EXPECT().StartResource(ctx, ca.Credentials, gomock.Any()).Return(nil)
			},
		},
		{
			name:  "Successfully stop SQL instance",
			input: ResourceDetails{CloudAccID: 123, Name: "test-instance", Type: SQL, State: SUSPEND},
			mockCalls: func() {
				mGCP.EXPECT().StopResource(ctx, ca.Credentials, gomock.Any()).Return(nil)
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

	mGCP := NewMockCloudResourceProvider(ctrl)
	ctx := &gofr.Context{Context: context.Background()}
	ca := &client.CloudAccount{ID: 123, Name: "MyCloud", Provider: string(GCP),
		Credentials: map[string]any{"project_id": "test-project", "region": "us-central1"}}
	s := New(mGCP, nil, nil, nil)

	testCases := []struct {
		name      string
		input     ResourceDetails
		cloudAcc  *client.CloudAccount
		expErr    error
		mockCalls func()
	}{
		{
			name:   "Error starting SQL instance",
			input:  ResourceDetails{CloudAccID: 123, Name: "test-instance", Type: SQL, State: START},
			expErr: errMock,
			mockCalls: func() {
				mGCP.EXPECT().StartResource(ctx, ca.Credentials, gomock.Any()).Return(errMock)
			},
		},
		{
			name:   "Error stopping SQL instance",
			input:  ResourceDetails{CloudAccID: 123, Name: "test-instance", Type: SQL, State: SUSPEND},
			expErr: errMock,
			mockCalls: func() {
				mGCP.EXPECT().StopResource(ctx, ca.Credentials, gomock.Any()).Return(errMock)
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
