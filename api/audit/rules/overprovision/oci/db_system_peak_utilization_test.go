package oci

import (
	"testing"

	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"

	"github.com/oracle/oci-go-sdk/v65/common"
	"github.com/oracle/oci-go-sdk/v65/database"
	"github.com/oracle/oci-go-sdk/v65/monitoring"
	"gofr.dev/pkg/gofr"
)

func TestCheckDBSystemProvisionedUsage(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	ctx := &gofr.Context{}

	testCases := []struct {
		name          string
		creds         *Credentials
		expectedError error
	}{
		{
			name: "invalid credentials",
			creds: &Credentials{
				TenancyOCID: "",
				UserOCID:    "",
				Region:      "",
				Fingerprint: "",
				PrivateKey:  "",
				Compartment: "",
			},
			expectedError: errCreateDBClient,
		},
		{
			name: "valid credentials",
			creds: &Credentials{
				TenancyOCID: "test-tenancy",
				UserOCID:    "test-user",
				Region:      "test-region",
				Fingerprint: "test-fingerprint",
				PrivateKey:  "test-key",
				Compartment: "test-compartment",
			},
			expectedError: nil,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Skip the actual OCI client creation for testing
			if tc.expectedError == nil {
				t.Skip("Skipping test that requires actual OCI client")
			}

			_, err := CheckDBSystemProvisionedUsage(ctx, tc.creds)

			if tc.expectedError != nil {
				require.Error(t, err)
				require.Equal(t, tc.expectedError, err)
			} else {
				require.NoError(t, err)
			}
		})
	}
}

func TestGetResult(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	ctx := &gofr.Context{}
	creds := &Credentials{
		TenancyOCID: "test-tenancy",
		UserOCID:    "test-user",
		Region:      "test-region",
		Fingerprint: "test-fingerprint",
		PrivateKey:  "test-key",
		Compartment: "test-compartment",
	}

	dbSystems := []database.DbSystemSummary{
		{
			Id:          common.String("test-id-1"),
			DisplayName: common.String("test-db-1"),
		},
		{
			Id:          common.String("test-id-2"),
			DisplayName: common.String("test-db-2"),
		},
	}

	testCases := []struct {
		name          string
		dbSystems     []database.DbSystemSummary
		expectedError error
	}{
		{
			name:          "empty db systems",
			dbSystems:     []database.DbSystemSummary{},
			expectedError: nil,
		},
		{
			name:          "multiple db systems",
			dbSystems:     dbSystems,
			expectedError: nil,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Skip the actual metrics retrieval for testing
			if len(tc.dbSystems) > 0 {
				t.Skip("Skipping test that requires actual metrics retrieval")
			}

			// Create a real monitoring client for empty DB systems test
			monitoringClient := &monitoring.MonitoringClient{}
			results, err := getResult(ctx, creds, tc.dbSystems, monitoringClient)

			if tc.expectedError != nil {
				require.Error(t, err)
				require.Equal(t, tc.expectedError, err)
			} else {
				require.NoError(t, err)
				require.NotNil(t, results)
			}
		})
	}
}

func TestGetOCICredentials(t *testing.T) {
	testCases := []struct {
		name          string
		creds         any
		expectedError error
	}{
		{
			name:          "nil credentials",
			creds:         nil,
			expectedError: errInvalidOCICreds,
		},
		{
			name: "invalid credentials type",
			creds: func() any {
				// Create a channel which cannot be marshaled to JSON
				return make(chan int)
			}(),
			expectedError: errInvalidOCICreds,
		},
		{
			name: "valid credentials",
			creds: &Credentials{
				TenancyOCID: "test-tenancy",
				UserOCID:    "test-user",
				Region:      "test-region",
				Fingerprint: "test-fingerprint",
				PrivateKey:  "test-key",
				Compartment: "test-compartment",
			},
			expectedError: nil,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			creds, err := getOCICredentials(tc.creds)

			if tc.expectedError != nil {
				require.Error(t, err)
				require.Equal(t, tc.expectedError, err)
			} else {
				require.NoError(t, err)
				require.NotNil(t, creds)
			}
		})
	}
}
