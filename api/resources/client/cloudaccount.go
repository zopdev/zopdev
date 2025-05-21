package client

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"

	"gofr.dev/pkg/gofr"
)

var (
	errFailedToGetCloudCredentials = errors.New("failed to get cloud credentials")
	errInvalidResponse             = errors.New("invalid response from cloud account service")
)

type Client struct{}

func New() *Client {
	return &Client{}
}

func (c *Client) GetCloudCredentials(ctx *gofr.Context, cloudAccID int64) (*CloudAccount, error) {
	// Fetch the cloud credentials from the cloud-account entity
	// This is a placeholder function and should be implemented based on the cloud provider
	endpoint := fmt.Sprintf("cloud-accounts/%d/credentials", cloudAccID)

	resp, err := ctx.GetHTTPService("cloud-account").
		Get(ctx, endpoint, nil)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, errFailedToGetCloudCredentials
	}

	var credentials struct {
		Data *CloudAccount `json:"data"`
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(body, &credentials)
	if err != nil {
		return nil, errInvalidResponse
	}

	return credentials.Data, nil
}
