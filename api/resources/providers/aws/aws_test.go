package aws

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNew(t *testing.T) {
	c := New()
	assert.NotNil(t, c)
}

func Test_getAWSCredentials(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		creds := map[string]string{"aws_access_key_id": "key", "aws_secret_access_key": "secret"}
		c, err := getAWSCredentials(creds)
		assert.NoError(t, err)
		assert.Equal(t, "key", c.AccessKey)
		assert.Equal(t, "secret", c.AccessSecret)
	})
	t.Run("missing both", func(t *testing.T) {
		creds := map[string]string{}
		_, err := getAWSCredentials(creds)
		assert.Error(t, err)
	})
	t.Run("missing key", func(t *testing.T) {
		creds := map[string]string{"aws_secret_access_key": "secret"}
		_, err := getAWSCredentials(creds)
		assert.Error(t, err)
	})
	t.Run("missing secret", func(t *testing.T) {
		creds := map[string]string{"aws_access_key_id": "key"}
		_, err := getAWSCredentials(creds)
		assert.Error(t, err)
	})
	t.Run("invalid json", func(t *testing.T) {
		_, err := getAWSCredentials(make(chan int))
		assert.Error(t, err)
	})
}

func Test_getAWSCredentials_edgeCases(t *testing.T) {
	t.Run("non-map input", func(t *testing.T) {
		_, err := getAWSCredentials(12345)
		assert.Error(t, err)
	})
	t.Run("empty struct", func(t *testing.T) {
		type empty struct{}
		_, err := getAWSCredentials(empty{})
		assert.Error(t, err)
	})
}

func Test_newSession(t *testing.T) {
	c := &Client{}
	t.Run("success", func(t *testing.T) {
		sess, err := c.newSession("key", "secret")
		assert.NoError(t, err)
		assert.NotNil(t, sess)
	})
}

func TestNewRDSClient_InvalidCreds(t *testing.T) {
	c := &Client{}
	_, err := c.NewRDSClient(context.Background(), map[string]string{})
	assert.Error(t, err)
}

func TestNewEC2Client_InvalidCreds(t *testing.T) {
	c := &Client{}
	_, err := c.NewEC2Client(context.Background(), map[string]string{})
	assert.Error(t, err)
}

func TestNewRDSClient_Success(t *testing.T) {
	c := &Client{}
	creds := map[string]string{"aws_access_key_id": "key", "aws_secret_access_key": "secret"}
	client, err := c.NewRDSClient(context.Background(), creds)
	assert.NoError(t, err)
	assert.NotNil(t, client)
}

func TestNewEC2Client_Success(t *testing.T) {
	c := &Client{}
	creds := map[string]string{"aws_access_key_id": "key", "aws_secret_access_key": "secret"}
	client, err := c.NewEC2Client(context.Background(), creds)
	assert.NoError(t, err)
	assert.NotNil(t, client)
}
