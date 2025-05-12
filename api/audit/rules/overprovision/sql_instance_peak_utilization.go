package overprovision

import (
	"errors"

	"gofr.dev/pkg/gofr"

	"github.com/zopdev/zopdev/api/audit/client"
	"github.com/zopdev/zopdev/api/audit/rules"
	"github.com/zopdev/zopdev/api/audit/rules/overprovision/gcp"
	"github.com/zopdev/zopdev/api/audit/store"
)

var errUnsupportedCloudProvider = errors.New("unsupported cloud provider")

type SQLInstancePeak struct {
}

func (*SQLInstancePeak) Execute(ctx *gofr.Context, ca *client.CloudAccount) ([]store.Items, error) {
	switch ca.Provider {
	case rules.GCP:
		return gcp.CheckCloudSQLProvisionedUsage(ctx, ca.Credentials)
	default:
		return nil, errUnsupportedCloudProvider
	}
}

func (*SQLInstancePeak) GetCategory() string {
	return "overprovision"
}

func (*SQLInstancePeak) GetName() string {
	return "sql_instance_peak"
}
