package resource

import (
	"gofr.dev/pkg/gofr"
	gofrHttp "gofr.dev/pkg/gofr/http"

	"github.com/zopdev/zopdev/api/resources/client"
	"github.com/zopdev/zopdev/api/resources/models"
)

type Service struct {
	gcp   CloudResourceProvider
	aws   CloudResourceProvider
	http  HTTPClient
	store Store
}

func New(gcp, aws CloudResourceProvider, http HTTPClient, store Store) *Service {
	return &Service{gcp: gcp, aws: aws, http: http, store: store}
}

func (s *Service) GetAll(ctx *gofr.Context, id int64, resourceType []string) ([]models.Resource, error) {
	res, err := s.store.GetResources(ctx, id, resourceType)
	if err != nil {
		return nil, err
	}

	return res, nil
}

func (s *Service) GetByID(ctx *gofr.Context, id int64) (*models.Resource, error) {
	res, err := s.store.GetResourceByID(ctx, id)
	if err != nil {
		return nil, err
	}

	if res == nil {
		return nil, gofrHttp.ErrorEntityNotFound{Name: "resource"}
	}

	return res, nil
}

// TODO : This function should be generic such that we dont need to call anything specific to a particular cloud provider type.
// and the error returned in case when the resource is starting or stopping or other operation, error should be returned carefully.
// or state should be managed.
func (s *Service) ChangeState(ctx *gofr.Context, resDetails ResourceDetails) error {
	res, err := s.store.GetResourceByID(ctx, resDetails.ID)
	if err != nil {
		return err
	}

	if res.Status == getStatus(resDetails.State) {
		return nil
	}

	ca, err := s.http.GetCloudCredentials(ctx, resDetails.CloudAccID)
	if err != nil {
		return err
	}

	switch resDetails.Type {
	case SQL, "RDS":
		return s.handleSQLChangeState(ctx, ca, resDetails)
	case AWSCOMPUTE:
		return s.handleAWSComputeChangeState(ctx, ca, resDetails, res)
	default:
		return gofrHttp.ErrorInvalidParam{Params: []string{"req.Type"}}
	}
}

func (s *Service) handleSQLChangeState(ctx *gofr.Context, ca *client.CloudAccount, resDetails ResourceDetails) error {
	err := s.changeSQLState(ctx, ca, resDetails)
	if err != nil {
		ctx.Errorf("failed to change SQL state: %v", err)
		return err
	}

	err = s.store.UpdateStatus(ctx, getStatus(resDetails.State), resDetails.ID)
	if err != nil {
		ctx.Errorf("failed to update resource status: %v", err)
	}

	return nil
}

func (s *Service) handleAWSComputeChangeState(ctx *gofr.Context, ca *client.CloudAccount, resDetails ResourceDetails,
	res *models.Resource) error {
	if resDetails.State == START {
		err := s.aws.StartResource(ctx, ca.Credentials, res)
		if err != nil {
			ctx.Errorf("failed to start EC2 instance: %v", err)
			return err
		}
	} else {
		err := s.aws.StopResource(ctx, ca.Credentials, res)
		if err != nil {
			ctx.Errorf("failed to stop EC2 instance: %v", err)
			return err
		}
	}

	err := s.store.UpdateStatus(ctx, getStatus(resDetails.State), resDetails.ID)
	if err != nil {
		ctx.Errorf("failed to update resource status: %v", err)
	}

	return nil
}

func getStatus(action ResourceState) string {
	switch action {
	case START:
		return RUNNING
	case SUSPEND:
		return STOPPED
	default:
		return ""
	}
}

func (s *Service) SyncResources(ctx *gofr.Context, id int64) ([]models.Resource, error) {
	ca, err := s.http.GetCloudCredentials(ctx, id)
	if err != nil {
		return nil, err
	}

	ins, err := s.getAllInstances(ctx, ca)
	if err != nil {
		ctx.Errorf("failed to get all instances: %v", err)
		return nil, err
	}

	res, err := s.store.GetResources(ctx, id, nil)
	if err != nil {
		ctx.Errorf("failed to get existing resources: %v", err)
		return nil, err
	}

	visited := make([]bool, len(res))

	for i := range ins {
		idx, found := bSearch(res, ins[i].UID)
		ins[i].CloudAccount = models.CloudAccount{ID: id, Type: ca.Provider}

		if !found {
			// This is true when the resource is present in the cloud but not in the store.
			err = s.store.InsertResource(ctx, &ins[i])
			if err != nil {
				ctx.Errorf("failed to insert resource: %v", err)
			}
		} else {
			// else update the existing resource and mark the resource as visited.
			visited[idx] = true
			ins[i].ID = res[idx].ID
			err = s.store.UpdateStatus(ctx, ins[i].Status, ins[i].ID)

			if err != nil {
				ctx.Errorf("failed to update resource: %v", err)
			}
		}
	}

	s.removeStale(ctx, visited, res)

	return s.GetAll(ctx, id, nil)
}

func (s *Service) removeStale(ctx *gofr.Context, visited []bool, res []models.Resource) {
	for i, v := range visited {
		if v {
			continue
		}

		// Remove a resource if it is not visited, i.e., The resource is no longer present on the cloud.
		err := s.store.RemoveResource(ctx, res[i].ID)
		if err != nil {
			ctx.Errorf("failed to remove resource: %v", err)
		}
	}
}

func (s *Service) getALLComputeInstances(ctx *gofr.Context, details CloudDetails) ([]models.Resource, error) {
	switch details.CloudType {
	case AWS:
		filter := models.ResourceFilter{
			ResourceTypes: []string{"EC2"},
		}

		return s.aws.ListResources(ctx, details.Creds, filter)
	case GCP:
		filter := models.ResourceFilter{
			ResourceTypes: []string{"COMPUTE"},
		}

		return s.gcp.ListResources(ctx, details.Creds, filter)
	default:
		// We are not returning any error because the sync process is completely internal, works on the cloud Account ID,
		// if we are getting an unknown cloud type, then this feature is not implemented and we simply return nil.
		return nil, nil
	}
}

// bSearch performs a binary search on the sorted slice of models.Resource.
func bSearch(res []models.Resource, uid string) (int, bool) {
	l, r := 0, len(res)-1

	for l <= r {
		mid := l + (r-l)/2
		if res[mid].UID == uid {
			return mid, true
		}

		if res[mid].UID < uid {
			l = mid + 1
		} else {
			r = mid - 1
		}
	}

	return -1, false
}
