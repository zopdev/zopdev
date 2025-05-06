package overprovision

import (
	"encoding/json"

	"github.com/zopdev/zopdev/api/audit/rules/overprovision/gcp"
	"gofr.dev/pkg/gofr"
)

func getGCPCredentials(ctx *gofr.Context, creds any) (*gcp.Credentials, error) {
	if creds == nil {
		return nil, errInvalidGCPCreds
	}

	b, err := json.Marshal(creds)
	if err != nil {
		return nil, errInvalidGCPCreds
	}

	var gcpCred gcp.Credentials

	err = json.Unmarshal(b, &gcpCred)
	if err != nil {
		return nil, errInvalidGCPCreds
	}

	return &gcpCred, nil
}
