package providers

import (
	gofrHttp "gofr.dev/pkg/gofr/http"

	service "github.com/zopdev/zopdev/api/resources/service/resource"

	"github.com/zopdev/zopdev/api/resources/providers/aws"
	"github.com/zopdev/zopdev/api/resources/providers/gcp"
)

// NewCloudResourceProvider returns a concrete provider implementation (e.g., *gcp.Client, *aws.Client).
// The caller should type assert to the desired interface.
func NewCloudResourceProvider(cloudType string) (service.CloudResourceProvider, error) {
	switch cloudType {
	case "GCP":
		return gcp.New(), nil
	case "AWS":
		return aws.New(), nil
	default:
		return nil, gofrHttp.ErrorEntityNotFound{
			Name:  "Cloud Provider",
			Value: cloudType,
		}
	}
}
