package gcp

import (
	"context"
	"encoding/json"
	"errors"
	"strings"

	gofrHttp "gofr.dev/pkg/gofr/http"

	"gofr.dev/pkg/gofr"

	"golang.org/x/oauth2/google"
	"google.golang.org/api/option"
	"google.golang.org/api/sqladmin/v1"

	gmonitoring "cloud.google.com/go/monitoring/apiv3/v2"
	"github.com/zopdev/zopdev/api/resources/models"
	sql "github.com/zopdev/zopdev/api/resources/providers/gcp/database"
	metric "github.com/zopdev/zopdev/api/resources/providers/gcp/monitoring"
)

var (
	ErrInvalidCredentials = errors.New("invalid cloud credentials")
	ErrInitializingClient = errors.New("error initializing client")
	ErrNotImplemented     = errors.New("not implemented")
)

const (
	resourceUIDSplitParts = 2
	database
	databaseSQL = "SQL"
	allResource = "ALL"
)

type Client struct{}

func New() *Client { return &Client{} }

func (*Client) NewGoogleCredentials(ctx context.Context, cred any, scopes ...string) (*google.Credentials, error) {
	var gcpCreds credentials

	b, _ := json.Marshal(cred)
	if err := json.Unmarshal(b, &gcpCreds); err != nil {
		return nil, ErrInvalidCredentials
	}

	creds, err := google.CredentialsFromJSON(ctx, b, scopes...)
	if err != nil {
		return nil, ErrInvalidCredentials
	}

	return creds, nil
}

func (*Client) NewSQLClient(ctx context.Context, opts ...option.ClientOption) (SQLClient, error) {
	admin, err := sqladmin.NewService(ctx, opts...)
	if err != nil {
		return nil, ErrInitializingClient
	}

	return &sql.Client{SQL: admin.Instances}, nil
}

func (*Client) NewMetricsClient(ctx context.Context, opts ...option.ClientOption) (MetricsClient, error) {
	mCl, err := gmonitoring.NewMetricClient(ctx, opts...)
	if err != nil {
		return nil, ErrInitializingClient
	}

	return &metric.Client{MetricClient: mCl}, nil
}

// Implement methods matching the CloudResourceProvider interface.
func (c *Client) ListResources(ctx *gofr.Context, creds any, filter models.ResourceFilter) ([]models.Resource, error) {
	var allResources []models.Resource

	resourceTypes := filter.ResourceTypes

	includeSQL := len(resourceTypes) == 0

	for _, t := range resourceTypes {
		if t == databaseSQL || t == allResource {
			includeSQL = true
		}
	}

	if includeSQL {
		credsObj, err := c.NewGoogleCredentials(ctx, creds, sqladmin.SqlserviceAdminScope)
		if err != nil {
			return nil, ErrInvalidCredentials
		}

		sqlClient, err := c.NewSQLClient(ctx, option.WithCredentials(credsObj))
		if err != nil {
			return nil, ErrInitializingClient
		}
		// Extract projectID from creds
		var credStruct struct{ ProjectID string }

		b, _ := json.Marshal(creds)
		_ = json.Unmarshal(b, &credStruct)

		instances, err := sqlClient.GetAllInstances(ctx, credStruct.ProjectID)
		if err != nil {
			return nil, err
		}

		allResources = append(allResources, instances...)
	}

	if len(allResources) == 0 {
		return nil, ErrNotImplemented
	}

	return allResources, nil
}

func (c *Client) StartResource(ctx *gofr.Context, creds any, resource *models.Resource) error {
	switch resource.Type {
	case databaseSQL:
		credsObj, err := c.NewGoogleCredentials(ctx, creds, sqladmin.SqlserviceAdminScope)
		if err != nil {
			return ErrInvalidCredentials
		}

		sqlClient, err := c.NewSQLClient(ctx, option.WithCredentials(credsObj))
		if err != nil {
			return ErrInitializingClient
		}
		// resource.UID is expected to be "projectID/instanceName"
		parts := strings.SplitN(resource.UID, "/", resourceUIDSplitParts)
		if len(parts) != resourceUIDSplitParts {
			return gofrHttp.ErrorInvalidParam{Params: []string{"resourceUID", resource.UID}}
		}

		return sqlClient.StartInstance(ctx, parts[0], parts[1])
	default:
		return ErrNotImplemented
	}
}

func (c *Client) StopResource(ctx *gofr.Context, creds any, resource *models.Resource) error {
	switch resource.Type {
	case "SQL":
		credsObj, err := c.NewGoogleCredentials(ctx, creds, sqladmin.SqlserviceAdminScope)
		if err != nil {
			return ErrInvalidCredentials
		}

		sqlClient, err := c.NewSQLClient(ctx, option.WithCredentials(credsObj))
		if err != nil {
			return ErrInitializingClient
		}
		// resource.UID is expected to be "projectID/instanceName"
		parts := strings.SplitN(resource.UID, "/", resourceUIDSplitParts)
		if len(parts) != resourceUIDSplitParts {
			return gofrHttp.ErrorInvalidParam{Params: []string{"UID", "SQL resource"}}
		}

		return sqlClient.StopInstance(ctx, parts[0], parts[1])
	default:
		return ErrNotImplemented
	}
}
