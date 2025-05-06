package client

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"gofr.dev/pkg/gofr"
)

func GetCloudCredentials(ctx *gofr.Context, cloudAccId int64) (*CloudAccount, error) {
	// Fetch the cloud credentials from the cloud-account entity
	// This is a placeholder function and should be implemented based on the cloud provider
	endpoint := fmt.Sprintf("/cloud-accounts/%d/credentials", cloudAccId)

	resp, err := ctx.GetHTTPService("cloud-account").Get(ctx, endpoint, nil)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to get cloud credentials: %s", resp.Status)
	}

	var credentials CloudAccount
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(body, &credentials)
	if err != nil {
		return nil, err
	}

	return nil, nil
}
