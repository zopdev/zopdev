package aws

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/stretchr/testify/assert"
	"gofr.dev/pkg/gofr"

	"github.com/zopdev/zopdev/api/resources/models"
)

func TestNew(t *testing.T) {
	c := New()
	assert.NotNil(t, c)
}

func Test_getAWSCredentials(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		creds := map[string]string{"aws_access_key_id": "key", "aws_secret_access_key": "secret"}
		c, err := getAWSCredentials(creds)
		require.NoError(t, err)
		require.Equal(t, "key", c.AccessKey)
		require.Equal(t, "secret", c.AccessSecret)
	})
	t.Run("missing both", func(t *testing.T) {
		creds := map[string]string{}
		_, err := getAWSCredentials(creds)
		require.Error(t, err)
	})
	t.Run("missing key", func(t *testing.T) {
		creds := map[string]string{"aws_secret_access_key": "secret"}
		_, err := getAWSCredentials(creds)
		require.Error(t, err)
	})
	t.Run("missing secret", func(t *testing.T) {
		creds := map[string]string{"aws_access_key_id": "key"}
		_, err := getAWSCredentials(creds)
		require.Error(t, err)
	})
	t.Run("invalid json", func(t *testing.T) {
		_, err := getAWSCredentials(make(chan int))
		require.Error(t, err)
	})
}

func Test_getAWSCredentials_edgeCases(t *testing.T) {
	t.Run("non-map input", func(t *testing.T) {
		_, err := getAWSCredentials(12345)
		require.Error(t, err)
	})
	t.Run("empty struct", func(t *testing.T) {
		type empty struct{}

		_, err := getAWSCredentials(empty{})
		require.Error(t, err)
	})
}

func Test_newSession(t *testing.T) {
	c := &Client{}

	t.Run("success", func(t *testing.T) {
		sess, err := c.newSession("key", "secret")
		require.NoError(t, err)
		require.NotNil(t, sess)
	})
}

func TestNewRDSClient_InvalidCreds(t *testing.T) {
	c := &Client{}
	_, err := c.NewRDSClient(context.Background(), map[string]string{})
	require.Error(t, err)
}

func TestNewEC2Client_InvalidCreds(t *testing.T) {
	c := &Client{}
	_, err := c.NewEC2Client(context.Background(), map[string]string{})
	require.Error(t, err)
}

func TestNewRDSClient_Success(t *testing.T) {
	c := &Client{}
	creds := map[string]string{"aws_access_key_id": "key", "aws_secret_access_key": "secret"}
	client, err := c.NewRDSClient(context.Background(), creds)
	require.NoError(t, err)
	require.NotNil(t, client)
}

func TestNewEC2Client_Success(t *testing.T) {
	c := &Client{}
	creds := map[string]string{"aws_access_key_id": "key", "aws_secret_access_key": "secret"}
	client, err := c.NewEC2Client(context.Background(), creds)
	require.NoError(t, err)
	require.NotNil(t, client)
}

func Test_shouldIncludeResourceType(t *testing.T) {
	tests := []struct {
		name          string
		resourceTypes []string
		targetType    string
		want          bool
	}{
		{
			name:          "empty resource types should include all",
			resourceTypes: []string{},
			targetType:    "EC2",
			want:          true,
		},
		{
			name:          "matching resource type should be included",
			resourceTypes: []string{"EC2", "RDS"},
			targetType:    "EC2",
			want:          true,
		},
		{
			name:          "ALL resource type should include everything",
			resourceTypes: []string{"ALL"},
			targetType:    "EC2",
			want:          true,
		},
		{
			name:          "non-matching resource type should be excluded",
			resourceTypes: []string{"RDS"},
			targetType:    "EC2",
			want:          false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := shouldIncludeResourceType(tt.resourceTypes, tt.targetType)
			assert.Equal(t, tt.want, got)
		})
	}
}

func Test_ListResources(t *testing.T) {
	c := &Client{}
	ctx := &gofr.Context{}

	tests := []struct {
		name         string
		creds        any
		filter       models.ResourceFilter
		wantErr      bool
		errorMessage string
	}{
		{
			name:  "invalid credentials",
			creds: map[string]string{},
			filter: models.ResourceFilter{
				ResourceTypes: []string{"EC2"},
			},
			wantErr:      true,
			errorMessage: "invalid cloud credentials",
		},
		{
			name:  "valid credentials but no resources",
			creds: map[string]string{"aws_access_key_id": "key", "aws_secret_access_key": "secret"},
			filter: models.ResourceFilter{
				ResourceTypes: []string{"UNKNOWN"},
			},
			wantErr:      true,
			errorMessage: "not implemented",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := c.ListResources(ctx, tt.creds, tt.filter)
			if tt.wantErr {
				require.Error(t, err)
				assert.Contains(t, err.Error(), tt.errorMessage)
				assert.Nil(t, got)
			}
		})
	}
}

func Test_StartResource(t *testing.T) {
	c := &Client{}
	ctx := &gofr.Context{}

	tests := []struct {
		name         string
		creds        any
		resource     *models.Resource
		wantErr      bool
		errorMessage string
	}{
		{
			name:  "invalid credentials",
			creds: map[string]string{},
			resource: &models.Resource{
				Type: "EC2",
				UID:  "i-1234567890abcdef0",
			},
			wantErr:      true,
			errorMessage: "invalid cloud credentials",
		},
		{
			name:  "unsupported resource type",
			creds: map[string]string{"aws_access_key_id": "key", "aws_secret_access_key": "secret"},
			resource: &models.Resource{
				Type: "UNKNOWN",
			},
			wantErr:      true,
			errorMessage: "not implemented",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := c.StartResource(ctx, tt.creds, tt.resource)
			if tt.wantErr {
				require.Error(t, err)
				assert.Contains(t, err.Error(), tt.errorMessage)
			}
		})
	}
}

func Test_StopResource(t *testing.T) {
	c := &Client{}
	ctx := &gofr.Context{}

	tests := []struct {
		name         string
		creds        any
		resource     *models.Resource
		wantErr      bool
		errorMessage string
	}{
		{
			name:  "invalid credentials",
			creds: map[string]string{},
			resource: &models.Resource{
				Type: "EC2",
				UID:  "i-1234567890abcdef0",
			},
			wantErr:      true,
			errorMessage: "invalid cloud credentials",
		},
		{
			name:  "unsupported resource type",
			creds: map[string]string{"aws_access_key_id": "key", "aws_secret_access_key": "secret"},
			resource: &models.Resource{
				Type: "UNKNOWN",
			},
			wantErr:      true,
			errorMessage: "not implemented",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := c.StopResource(ctx, tt.creds, tt.resource)
			if tt.wantErr {
				require.Error(t, err)
				assert.Contains(t, err.Error(), tt.errorMessage)
			}
		})
	}
}
