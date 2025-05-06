package overprovision

import (
	"errors"

	"gofr.dev/pkg/gofr"

	"github.com/zopdev/zopdev/api/audit/client"
	"github.com/zopdev/zopdev/api/audit/rules/overprovision/gcp"
	"github.com/zopdev/zopdev/api/audit/store"
)

var errInvalidGCPCreds = errors.New("invalid GCP credentials")

type SQLInstancePeak struct {
}

func (*SQLInstancePeak) Execute(ctx *gofr.Context, ca *client.CloudAccount) ([]store.Items, error) {
	switch ca.Provider {
	case "gcp":
		creds, err := getGCPCredentials(ctx, ca.Credentials)
		if err != nil {
			return nil, err
		}

		return gcp.CheckCloudSQLProvisionedUsage(ctx, creds)

	default:
		return nil, errors.New("unsupported cloud provider")
	}
}

func (*SQLInstancePeak) GetCategory() string {
	return "overprovision"
}

func (*SQLInstancePeak) GetName() string {
	return "sql_instance_peak"
}
