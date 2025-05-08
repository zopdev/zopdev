package overprovision

import (
	"encoding/json"
	"errors"

	"github.com/zopdev/zopdev/api/audit/rules/overprovision/gcp"
)

var (
	errInvalidGCPCreds = errors.New("invalid GCP credentials")
)

func getGCPCredentials(creds any) (*gcp.Credentials, error) {
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
